package server

import (
	"context"
	"fmt"

	"github.com/duckpie/bfb-security-microservice/internal/config"
	"github.com/duckpie/bfb-security-microservice/internal/core"
	"github.com/duckpie/bfb-security-microservice/internal/db/redisstore"
	pb "github.com/wrs-news/golang-proto/pkg/proto/security"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedSecurityServiceServer

	server *grpc.Server
	redis  redisstore.RedisStoreI

	conn map[core.ResolverConnKey]*grpc.ClientConn
	cfg  *config.ServerConfig
}

type ServerI interface {
	GetServer() *grpc.Server
	GetConn(key core.ResolverConnKey) (*grpc.ClientConn, error)
	AddConnection(key core.ResolverConnKey, connect func() (*grpc.ClientConn, error)) error
}

func (s *Server) GetServer() *grpc.Server {
	return s.server
}

func (s *Server) GetConn(key core.ResolverConnKey) (*grpc.ClientConn, error) {
	if val, ok := s.conn[key]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("established connection to the %s not found", key)
}

func (s *Server) AddConnection(key core.ResolverConnKey, connect func() (*grpc.ClientConn, error)) error {
	conn, err := connect()
	if err != nil {
		return err
	}

	s.conn[key] = conn
	return err
}

func (s *Server) HeartbeatCheck(ctx context.Context, e *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func InitServer(cfg *config.ServerConfig, r redisstore.RedisStoreI) *Server {
	return &Server{
		server: grpc.NewServer(),
		cfg:    cfg,
		redis:  r,
		conn:   make(map[core.ResolverConnKey]*grpc.ClientConn),
	}
}
