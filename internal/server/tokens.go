package server

import (
	"context"
	"fmt"
	"time"

	"github.com/duckpie/bfb-security-microservice/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
	pbu "github.com/wrs-news/golang-proto/pkg/proto/user"
)

func (s *Server) createToken(u *pbu.User) (*model.TokenDetails, error) {
	// Набор информации о пользовательских токенах и иж сроки действия
	td := &model.TokenDetails{}
	var err error

	/* Определение времени жизни для токенов */

	// Определяю время жизни в 15 МИНУТ для токена ДОСТУПА
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	// Создаю идентификатор для токена доступа
	td.AccessUuid = uuid.NewV4().String()

	// Определяю время жизни в 7 ДНЕЙ для токена ОБНОВЛЕНИЯ
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	/* Генерация токена доступа */

	// Создаю полезную нагрузку токена
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["uuid"] = u.Uuid
	atClaims["login"] = u.Login
	atClaims["role"] = u.Role
	atClaims["exp"] = td.AtExpires

	// Кодирую полезную нагрузку создавая ТОКЕН ДОСТУПА
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(s.cfg.AccessSecret))
	if err != nil {
		return nil, err
	}

	/* Генерация токена обновления */
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["uuid"] = u.Uuid
	rtClaims["login"] = u.Login
	rtClaims["role"] = u.Role
	rtClaims["exp"] = td.RtExpires

	// Кодирую полезную нагрузку создавая ТОКЕН ОБНОВЛЕНИЯ
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(s.cfg.RefreshSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *Server) createAuth(ctx context.Context, td *model.TokenDetails) error {
	// Конвертация access_token из Unix формата в UTC
	at := time.Unix(td.AtExpires, 0)
	// Конвертация refresh_token из Unix формата в UTC
	rt := time.Unix(td.RtExpires, 0)

	now := time.Now()

	// Сохранение access_tokenа
	if errAccess := s.redis.Save(ctx, td.AccessUuid, td.AccessToken, at.Sub(now)); errAccess != nil {
		return errAccess
	}

	// Сохранение refresh_tokenа
	if errRefresh := s.redis.Save(ctx, td.RefreshUuid, td.RefreshToken, rt.Sub(now)); errRefresh != nil {
		return errRefresh
	}

	return nil
}

// Верификация JWT токена
func (s *Server) verifyToken(tokenStr string) (*jwt.Token, error) {
	// Извлекаю токен в виде структуры
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Проверяю соответствие подписи токена с методом SigningMethodHMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.cfg.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Извлечение мета-данных из токена
func (s *Server) extractTokenMetadata(tokenStr string) (*model.AccessDetails, error) {
	token, err := s.verifyToken(tokenStr)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		return &model.AccessDetails{
			AccessUuid: accessUuid,
			Uuid:       claims["uuid"].(string),
			Login:      claims["login"].(string),
			Role:       int(claims["role"].(float64)),
		}, nil
	}

	return nil, err
}
