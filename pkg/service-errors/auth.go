package service_errors

import "errors"

var (
	UserBlocked        = errors.New("пользователь заблокирован")
	PasswordNotMatched = errors.New("не правильный пароль")
	RefreshTokenFailed = errors.New("не удалось сформировать токен")
)
