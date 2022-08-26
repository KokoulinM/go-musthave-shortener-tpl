package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/configs"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/database/filebase"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/database/postgres"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/grpc_handlers"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/handlers"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/helpers/certificate"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/pb"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/router"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/server"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/services"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/workers"
)

var (
	httpServer   *http.Server
	grpcServer   *grpc.Server
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	log.Printf("Build version: %v\n", buildVersion)
	log.Printf("Build date: %v\n", buildDate)
	log.Printf("Build commit: %v\n", buildCommit)

	err := certificate.Generate()
	if err != nil {
		log.Fatal("There was a problem when generating the certificate")
	}

	var httpServer *server.Server

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	defer signal.Stop(interrupt)

	cfg := configs.New()

	var service *services.URLService
	_, subnet, err := net.ParseCIDR(cfg.TrustedSubnet)

	wp := workers.New(ctx, cfg.Workers, cfg.WorkersBuffer)

	go func() {
		wp.Run(ctx)
	}()
	defer wp.Stop()
	if cfg.DatabaseDSN != "" {
		conn, err := postgres.Conn("postgres", cfg.DatabaseDSN)
		if err != nil {
			log.Printf("Unable to connect to the database: %s", err.Error())
		}

		err = postgres.SetUpDataBase(ctx, conn)

		if err != nil {
			log.Printf("Unable to create database struct: %s", err.Error())
		}

		service = services.New(postgres.NewDatabaseRepository(cfg.BaseURL, conn), cfg.BaseURL, wp, subnet)
	} else {
		service = services.New(filebase.NewFileRepository(ctx, cfg.FileStoragePath, cfg.BaseURL), cfg.BaseURL, wp, subnet)
	}

	g, ctx := errgroup.WithContext(ctx)

	h := handlers.New(service, cfg.BaseURL, wp)
	grpcHandler := grpc_handlers.NewGRPCHandler(service)

	mux := router.New(h, cfg)

	g.Go(func() error {
		httpServer = server.New(cfg.ServerAddress, cfg.Key, mux)

		var err error

		if cfg.EnableHttps {
			err = httpServer.StartTLS("cert.pem", "key.pem")
		} else {
			err = httpServer.Start()
		}

		if err != nil {
			return err
		}

		log.Printf("httpServer starting at: %v", cfg.ServerAddress)

		return nil
	})

	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
		if err != nil {
			log.Printf("gRPC server failed to listen: %v", err.Error())
			return err
		}
		grpcServer = grpc.NewServer()
		pb.RegisterURLServer(grpcServer, grpcHandler)
		log.Printf("server listening at %v", lis.Addr())
		return grpcServer.Serve(lis)
	})

	select {
	case <-interrupt:
		log.Println("Stop server")
		break
	case <-ctx.Done():
		break
	}

	log.Println("Receive shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdownCancel()

	if httpServer != nil {
		_ = httpServer.Shutdown(shutdownCtx)
	}

	if grpcServer != nil {
		grpcServer.GracefulStop()
	}

	err = g.Wait()
	if err != nil {
		log.Printf("server returning an error: %v", err)
	}

	log.Println("Server Shutdown gracefully")
}
