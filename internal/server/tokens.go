package server

import (
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
