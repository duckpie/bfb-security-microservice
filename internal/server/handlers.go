package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/duckpie/cherry"
	cherrynet "github.com/duckpie/cherry/net"
	"github.com/golang-jwt/jwt"
	pb "github.com/wrs-news/golang-proto/pkg/proto/security"
	pbu "github.com/wrs-news/golang-proto/pkg/proto/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Login(ctx context.Context, in *pb.LoginReq) (*pb.TokensPair, error) {
	client, err := s.GetConn(cherry.UMS)
	if err != nil {
		es := status.New(codes.Internal, err.Error())
		return nil, es.Err()
	}

	conn := pbu.NewUserServiceClient(client)
	rptr := cherrynet.GrpcRepeater(func(ctx context.Context) (interface{}, error) {
		resp, err := conn.GetUserByLogin(ctx, &pbu.UserReqLogin{
			Login: in.Login,
		})
		if err != nil {
			return nil, err
		}

		return resp, nil
	}, 3, time.Second)

	resp, err := rptr(ctx)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(resp.(*pbu.User).Hash), []byte(in.Password)); err != nil {
		es := status.New(codes.Internal, err.Error())
		return nil, es.Err()
	}

	td, err := s.createToken(resp.(*pbu.User))
	if err != nil {
		es := status.New(codes.Internal, err.Error())
		return nil, es.Err()
	}

	if err := s.createAuth(ctx, td); err != nil {
		es := status.New(codes.Internal, err.Error())
		return nil, es.Err()
	}

	return &pb.TokensPair{
		RefreshToken: td.RefreshToken,
		AccessToken:  td.AccessToken,
	}, nil
}

func (s *Server) AuthCheck(ctx context.Context, in *pb.AuthCheckReq) (*emptypb.Empty, error) {
	tokenDt, err := s.extractTokenMetadata(in.AccessToken)
	if err != nil {
		es := status.New(codes.InvalidArgument, err.Error())
		return nil, es.Err()
	}

	if _, err := s.redis.Get(ctx, tokenDt.AccessUuid); err != nil {
		es := status.New(codes.NotFound, err.Error())
		return nil, es.Err()
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) RefreshToken(ctx context.Context, in *pb.RefreshTokenReq) (*pb.TokensPair, error) {
	// Верификация refresh токена
	token, err := jwt.Parse(in.Token, func(token *jwt.Token) (interface{}, error) {
		// Проверяю соответствие подписи токена с методом SigningMethodHMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			es := status.New(codes.InvalidArgument, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]).Error())
			return nil, es.Err()
		}

		return []byte(s.cfg.RefreshSecret), nil
	})
	if err != nil {
		es := status.New(codes.InvalidArgument, err.Error())
		return nil, es.Err()
	}

	// Проверка валидности токена
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		es := status.New(codes.InvalidArgument, err.Error())
		return nil, es.Err()
	}

	// Проверка на соответствие с MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// Удаляю предыдущий refresh токен
		if err := s.redis.Delete(ctx, claims["refresh_uuid"].(string)); err != nil {
			es := status.New(codes.InvalidArgument, errors.New("unauthorized").Error())
			return nil, es.Err()
		}

		u := pbu.User{
			Uuid:  claims["uuid"].(string),
			Login: claims["login"].(string),
			Role:  int32(claims["role"].(float64)),
		}

		// Создание новой пары токенов
		td, err := s.createToken(&u)
		if err != nil {
			es := status.New(codes.Internal, err.Error())
			return nil, es.Err()
		}

		// Сохранение токенов в redis
		if err := s.createAuth(ctx, td); err != nil {
			es := status.New(codes.Internal, err.Error())
			return nil, es.Err()
		}

		return &pb.TokensPair{
			AccessToken:  td.AccessToken,
			RefreshToken: td.RefreshToken,
		}, nil
	}

	return nil, status.New(codes.InvalidArgument, errors.New("refresh expired").Error()).Err()
}

func (s *Server) Logout(ctx context.Context, in *pb.LogoutReq) (*emptypb.Empty, error) {
	tokenDt, err := s.extractTokenMetadata(in.Token)
	if err != nil {
		es := status.New(codes.Internal, err.Error())
		return nil, es.Err()
	}

	if err := s.redis.Delete(ctx, tokenDt.AccessUuid); err != nil {
		es := status.New(codes.Internal, err.Error())
		return nil, es.Err()
	}

	return &emptypb.Empty{}, nil
}
