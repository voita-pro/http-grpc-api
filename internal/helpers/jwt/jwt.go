package helperjwt

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/voita-pro/http-grpc-api/internal/config"
	"time"
)

const INTERNAL_SECRET string = "0b8a8664f7a84607bb89059526e3350989d02275"

type Auth struct {
	ctx context.Context
}

func New(ctx context.Context) (*Auth, error) {
	return &Auth{
		ctx: ctx,
	}, nil
}
func (s *Auth) CustomTokenWithClaims(ctx context.Context, uid string, devClaims map[string]interface{}) (string, error) {
	cfg := config.Get()
	claims := jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(cfg.HTTP.ExpToken).Unix(),
	}
	for key, value := range devClaims {
		if key == "uid" || key == "exp" {
			continue
		}
		claims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.GRPC.Secret + INTERNAL_SECRET + uid))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *Auth) VerifyToken(ctx context.Context, tokenString string) (map[string]interface{}, error) {
	cfg := config.Get()
	tokenData, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("error in parsing token")
			}
			uid, _ := token.Claims.(jwt.MapClaims)["uid"].(string)
			secretKey := []byte(cfg.GRPC.Secret + INTERNAL_SECRET + uid)
			return secretKey, nil
		})

	if err != nil {
		return nil, err
	}
	if claims, ok := tokenData.Claims.(jwt.MapClaims); ok && tokenData.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
