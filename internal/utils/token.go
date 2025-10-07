package utils

import (
	"context"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"auth/internal/model"
)

const (
	AuthPrefix = "Bearer "

	RefreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	AccessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="

	RefreshTokenExpiration = 60 * time.Minute
	AccessTokenExpiration  = 35 * time.Minute
)

// revokedTokens хранит отозванные токены до момента их естественного истечения срока действия.
// key: исходная строка токена, value: unix-время истечения токена
var revokedTokens sync.Map

// RevokeToken добавляет токен в список отозванных до указанного времени истечения.
func RevokeToken(tokenStr string, expiresAt int64) {
	if tokenStr == "" || expiresAt == 0 {
		return
	}
	revokedTokens.Store(tokenStr, expiresAt)
}

// IsTokenRevoked проверяет, отозван ли токен. Автоматически очищает просроченные записи.
func IsTokenRevoked(tokenStr string) bool {
	v, ok := revokedTokens.Load(tokenStr)
	if !ok {
		return false
	}

	exp, _ := v.(int64)
	now := time.Now().Unix()
	if now >= exp {
		// истёк — убираем из списка, не считаем отозванным далее
		revokedTokens.Delete(tokenStr)
		return false
	}
	return true
}

func GenerateToken(info model.User, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Name: info.Name,
		Role: info.Role,
		ID:   info.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	if IsTokenRevoked(tokenStr) {
		return nil, errors.Errorf("token revoked")
	}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	if IsTokenRevoked(tokenStr) {
		return nil, errors.Errorf("token revoked")
	}

	return claims, nil
}

// StartRevokedTokensJanitor запускает фоновую очистку просроченных записей revokedTokens.
func StartRevokedTokensJanitor(ctx context.Context, sweepInterval time.Duration) {
	if sweepInterval <= 0 {
		sweepInterval = 5 * time.Minute
	}

	ticker := time.NewTicker(sweepInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now().Unix()
				revokedTokens.Range(func(key, value interface{}) bool {
					exp, ok := value.(int64)
					if !ok || now >= exp {
						revokedTokens.Delete(key)
					}
					return true
				})
			}
		}
	}()
}
