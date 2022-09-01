// Package grpchandlers provide methods for working with grpc
package grpchandlers

import (
	"context"
	"net"
	"net/http"
	"strconv"

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

func (us *URLServer) CreateShortURL(ctx context.Context, in *pb.CreateShortURLRequest) (*pb.CreateShortURLResponse, error) {
	responseURL, err := us.service.CreateURL(ctx, in.OriginalId, in.UserId)
	if err != nil {
		statusCode := errors.ParseError(err)
		switch statusCode {
		case http.StatusConflict:
			return &pb.CreateShortURLResponse{
				Status: "conflict",
			}, nil
		default:
			return &pb.CreateShortURLResponse{
				Status: "internal server error",
			}, nil
		}
	}
	return &pb.CreateShortURLResponse{
		Status:      "ok",
		ResponseUrl: responseURL,
	}, nil
}

func (us *URLServer) GetUserURLs(ctx context.Context, in *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	urls, err := us.service.GetUserURLs(ctx, in.UserId)
	if err != nil {
		statusCode := errors.ParseError(err)
		switch statusCode {
		case http.StatusNoContent:
			return &pb.GetUserURLsResponse{
				Status: "no content",
			}, nil
		default:
			return &pb.GetUserURLsResponse{
				Status: "internal server error",
			}, nil
		}
	}
	var result []*pb.GetUserURLsResponse_URL
	for i := 0; i < len(urls); i++ {
		result = append(result, &pb.GetUserURLsResponse_URL{
			OriginalUrl: urls[0].OriginalURL,
			ShortUrl:    urls[0].ShortURL,
		})
	}
	return &pb.GetUserURLsResponse{
		Status: "ok",
		Urls:   result,
	}, nil
}

func (us *URLServer) CreateBatch(ctx context.Context, in *pb.CreateBatchRequest) (*pb.CreateBatchResponse, error) {
	var data []handlers.RequestGetURLs
	for i := 0; i < len(in.Urls); i++ {
		data = append(data, handlers.RequestGetURLs{
			CorrelationID: strconv.Itoa(int(in.Urls[i].CorrelationId)),
			OriginalURL:   in.Urls[i].OriginalUrl,
		})
	}
	urls, err := us.service.CreateBatch(ctx, data, in.UserId)
	if err != nil {
		return &pb.CreateBatchResponse{
			Status: "internal server error",
		}, nil
	}
	var response []*pb.CreateBatchResponse_URL
	for i := 0; i < len(urls); i++ {
		id, _ := strconv.ParseInt(urls[i].CorrelationID, 10, 32)
		response = append(response, &pb.CreateBatchResponse_URL{
			CorrelationId: int32(id),
			ShortUrl:      urls[i].ShortURL,
		})
	}
	return &pb.CreateBatchResponse{
		Status: "ok",
		Urls:   response,
	}, nil
}

func (us *URLServer) DeleteBatch(ctx context.Context, in *pb.DeleteBatchRequest) (*pb.DeleteBatchResponse, error) {
	us.service.DeleteBatch(in.Urls, in.UserId)
	return &pb.DeleteBatchResponse{
		Status: "accepted",
	}, nil
}

func (us *URLServer) GetStates(ctx context.Context, in *pb.GetStatesRequest) (*pb.GetStatesResponse, error) {
	hasPermission, response, err := us.service.GetStates(ctx, net.IP(in.IpAddress))
	if !hasPermission {
		return &pb.GetStatesResponse{
			Status: "forbidden",
		}, nil
	}
	if err != nil {
		return &pb.GetStatesResponse{
			Status: "bad request",
		}, nil
	}
	return &pb.GetStatesResponse{
		Status: "ok",
		Users:  int32(response.Users),
		Urls:   int32(response.Urls),
	}, nil
}
