package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/duckpie/bfb-security-microservice/internal/core"
	"github.com/golang-jwt/jwt"
	pb "github.com/wrs-news/golang-proto/pkg/proto/security"
	pbu "github.com/wrs-news/golang-proto/pkg/proto/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Login(ctx context.Context, in *pb.LoginReq) (*pb.TokensPair, error) {
	client, err := s.GetConn(core.UMS)
	if err != nil {
		return nil, err
	}

	conn := pbu.NewUserServiceClient(client)

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

	if err := s.createAuth(ctx, td); err != nil {
		return nil, err
	}

	return &pb.TokensPair{
		RefreshToken: td.RefreshToken,
		AccessToken:  td.AccessToken,
	}, nil
}

func (s *Server) AuthCheck(ctx context.Context, in *pb.AuthCheckReq) (*emptypb.Empty, error) {
	tokenDt, err := s.extractTokenMetadata(in.AccessToken)
	if err != nil {
		return nil, err
	}

	if _, err := s.redis.Get(ctx, tokenDt.AccessUuid); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) RefreshToken(ctx context.Context, in *pb.RefreshTokenReq) (*pb.TokensPair, error) {
	// Верификация refresh токена
	token, err := jwt.Parse(in.Token, func(token *jwt.Token) (interface{}, error) {
		// Проверяю соответствие подписи токена с методом SigningMethodHMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.cfg.RefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// Проверка валидности токена
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, err
	}

	// Проверка на соответствие с MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// Удаляю предыдущий refresh токен
		if err := s.redis.Delete(ctx, claims["refresh_uuid"].(string)); err != nil {
			return nil, errors.New("unauthorized")
		}

		u := pbu.User{
			Uuid:  claims["uuid"].(string),
			Login: claims["login"].(string),
			Role:  int32(claims["role"].(float64)),
		}

		// Создание новой пары токенов
		td, err := s.createToken(&u)
		if err != nil {
			return nil, err
		}

		// Сохранение токенов в redis
		if err := s.createAuth(ctx, td); err != nil {
			return nil, err
		}

		return &pb.TokensPair{
			AccessToken:  td.AccessToken,
			RefreshToken: td.RefreshToken,
		}, nil
	}

	return nil, errors.New("refresh expired")
}

func (s *Server) Logout(ctx context.Context, in *pb.LogoutReq) (*emptypb.Empty, error) {
	tokenDt, err := s.extractTokenMetadata(in.Token)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, s.redis.Delete(ctx, tokenDt.AccessUuid)
}
