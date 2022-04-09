package server

import (
	"context"
	"fmt"

	"github.com/duckpie/bfb-security-microservice/internal/config"
	pb "github.com/wrs-news/golang-proto/pkg/proto/security"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedSecurityServiceServer

	server   *grpc.Server
	userConn *grpc.ClientConn

	cfg *config.ServerConfig
}

type ServerI interface {
	GetServer() *grpc.Server
	ConnectToUserService(host string, port int) error
}

func (s *Server) GetServer() *grpc.Server {
	return s.server
}

func (s *Server) ConnectToUserService(host string, port int) error {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return err
	}

	s.userConn = conn
	return nil
}

func (s *Server) HeartbeatCheck(ctx context.Context, e *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func InitServer(cfg *config.ServerConfig) *Server {
	return &Server{
		server: grpc.NewServer(),
		cfg:    cfg,
	}
}
