package server

import (
	"context"
	"log"
	"net"

	"github.com/Mldlr/url-shortener/internal/app/config"
	handler "github.com/Mldlr/url-shortener/internal/app/grpc/handlers"
	"github.com/Mldlr/url-shortener/internal/app/grpc/interceptors"
	pb "github.com/Mldlr/url-shortener/internal/app/grpc/proto"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"google.golang.org/grpc"
)

// GRPCServer is a gRPC server shortener api
type GRPCServer struct {
	handler *handler.ShortenerHandler
	cfg     *config.Config
}

// NewGRPCServer creates new gRPC server
func NewGRPCServer(shortener service.ShortenerService, cfg *config.Config) *GRPCServer {
	return &GRPCServer{
		handler: handler.NewShortenerHandler(shortener),
		cfg:     cfg,
	}
}

// Run starts the gRPC server
func (s *GRPCServer) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", s.cfg.GRPCAddress)
	if err != nil {
		log.Fatal(err)
	}
	// auth and trust net interceptors
	trustInterceptor := interceptors.Trusted{Config: s.cfg}
	authInterceptor := interceptors.Auth{Config: s.cfg}
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(trustInterceptor.TrustInterceptor, authInterceptor.AuthInterceptor))
	pb.RegisterShortenerServer(srv, s.handler)
	if err = srv.Serve(listener); err != nil {
		log.Fatal(err)
	}
	<-ctx.Done()
	srv.GracefulStop()
}
