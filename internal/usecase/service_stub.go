package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/voita-pro/http-grpc-api/internal/domain"
)

type serviceStub struct {
	ctx  context.Context
	auth domain.Authenticator
}

func NewStub(ctx context.Context, auth domain.Authenticator) domain.Usecase {
	return &serviceStub{
		ctx:  ctx,
		auth: auth,
	}
}

func (s *serviceStub) VerifyToken(idToken string) (*domain.TokenData, error) {
	if idToken == "" {
		return nil, errors.New("id token required")
	}
	claims, err := s.auth.VerifyToken(s.ctx, idToken)
	if err != nil {
		return nil, err
	}
	return &domain.TokenData{
		UID:           claims["uid"].(string),
		Email:         claims["email"].(string),
		EmailVerified: claims["email_verified"].(bool),
		DisplayName:   claims["display_name"].(string),
		Role:          claims["role"].(string),
	}, nil
}

func (s *serviceStub) SignIn(login domain.Login) (string, error) {
	hash := md5.Sum([]byte(login.Email))
	uid := hex.EncodeToString(hash[:])
	log.Info().Msgf("uid: %s", uid)
	log.Info().Msgf("email: %s", login.Email)

	token, err := s.auth.CustomTokenWithClaims(s.ctx, uid, map[string]interface{}{
		"email":          login.Email,
		"uid":            uid,
		"email_verified": false,
		"display_name":   "Just User",
		"role":           "user",
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *serviceStub) Countries() ([]*domain.Country, error) {
	countries := make([]*domain.Country, 0, 3)
	countries = append(countries, &domain.Country{Id: 1, Title: "Portugal", Code: "PT"})
	countries = append(countries, &domain.Country{Id: 2, Title: "Spain", Code: "ES"})
	countries = append(countries, &domain.Country{Id: 3, Title: "Italy", Code: "IT"})
	return countries, nil
}
func (s *serviceStub) CountryAdd(user *domain.TokenData, country domain.Country) (*domain.Country, error) {
	if country.Title == "" || country.Code == "" {
		return nil, fmt.Errorf("invalid country data")
	}
	if country.Id != 0 {
		return nil, fmt.Errorf("country already exists")
	}
	res := country
	res.Id = 4
	return &res, nil
}

func (s *serviceStub) CountrySave(user *domain.TokenData, country domain.Country) (*domain.Country, error) {
	if country.Title == "" || country.Code == "" {
		return nil, fmt.Errorf("invalid country data")
	}
	if country.Id != 0 {
		return s.CountryAdd(user, country)
	}

	res := country
	return &res, nil
}
