package grpc_handlers

import (
	"context"
	"net/http"

	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/errors"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/handlers"
	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/pb"
)

func NewGRPCHandler(service handlers.URLServiceInterface) *URLServer {
	return &URLServer{
		service: service,
	}
}

type URLServer struct {
	pb.UnimplementedURLServer
	service handlers.URLServiceInterface
}

func (us *URLServer) RetrieveShortURL(ctx context.Context, in *pb.RetrieveShortURLRequest) (*pb.RetrieveShortURLResponse, error) {
	longURL, err := us.service.GetURL(ctx, in.ShortUrlId)
	if err != nil {
		statusCode := errors.ParseError(err)
		switch statusCode {
		case http.StatusGone:
			return &pb.RetrieveShortURLResponse{
				Status: "gone",
			}, nil
		case http.StatusNotFound:
			return &pb.RetrieveShortURLResponse{
				Status: "not found",
			}, nil
		default:
			return &pb.RetrieveShortURLResponse{
				Status: "internal server error",
			}, nil
		}
	}

	return &pb.RetrieveShortURLResponse{
		RedirectUrl: longURL,
		Status:      "ok",
	}, nil
}
