package model

// Структура набора пользовательского токена
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string // Идентификатор токена доступа
	RefreshUuid  string
	AtExpires    int64 // Время жизни токена доступа
	RtExpires    int64
}

// Структура метаданных молезной нагрузки Access Token
type AccessDetails struct {
	Uuid       string
	AccessUuid string
	Login      string
	Role       int
}
