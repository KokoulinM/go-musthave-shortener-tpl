package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/handlers/middlewares"
	_ "github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/database"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers/db"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/router"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	cfg := configs.New()

	conn, err := db.Conn("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Println("Closing connect to db")
		err := conn.Close()
		if err != nil {
			log.Println("Closing don't close")
		}
	}

	db := database.New(conn)

	defer db.Conn.Close()

	handler := router.New(db, cfg)

	//serv := server.New(db, cfg.ServerAddress, handler, cfg)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: middlewares.Conveyor(handler, middlewares.GzipMiddleware, middlewares.CookieMiddleware),
	}

	//shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	//
	//defer shutdownCancel()
	//
	//if srv != nil {
	//	_ = srv.Shutdown(shutdownCtx)
	//}

	go func() {
		log.Fatal("app error exit", zap.Error(http.ListenAndServe(srv.Addr, srv.Handler)))
		//serv.Start()
	}()

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		fmt.Println("Got SIGINT...")
	case syscall.SIGTERM:
		fmt.Println("Got SIGTERM...")
	}
}
