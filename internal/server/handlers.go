package server

import (
	"context"

	pb "github.com/wrs-news/golang-proto/pkg/proto/security"
	pbu "github.com/wrs-news/golang-proto/pkg/proto/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Login(ctx context.Context, in *pb.LoginReq) (*pb.Token, error) {
	conn := pbu.NewUserServiceClient(s.userConn)

	resp, err := conn.GetUserByLogin(ctx, &pbu.UserReqLogin{
		Login: in.Login,
	})
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(resp.Hash), []byte(in.Password)); err != nil {
		return nil, err
	}

	td, err := s.createToken(resp)
	if err != nil {
		return nil, err
	}

	return &pb.Token{RefreshToken: td.RefreshToken}, nil
}

func (s *Server) AuthCheck(ctx context.Context, in *emptypb.Empty) (*pb.AuthCheckRes, error) {
	return &pb.AuthCheckRes{}, nil
}

func (s *Server) RefreshToken(ctx context.Context, in *pb.Token) (*pb.Token, error) {
	return &pb.Token{}, nil
}

func (s *Server) Logout(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
