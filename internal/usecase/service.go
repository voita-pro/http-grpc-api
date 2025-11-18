package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/voita-pro/http-grpc-api/internal/domain"
	"github.com/voita-pro/http-grpc-api/internal/repository/pgdb"
)

type service struct {
	ctx  context.Context
	auth domain.Authenticator
	repo pgdb.Querier
}

func New(ctx context.Context, auth domain.Authenticator, repo pgdb.Querier) domain.Usecase {
	return &service{
		ctx:  ctx,
		auth: auth,
		repo: repo,
	}
}

// VerifyToken
func (s *service) VerifyToken(idToken string) (*domain.TokenData, error) {
	if idToken == "" {
		return nil, errors.New("id token required")
	}
	claims, err := s.auth.VerifyToken(s.ctx, idToken)
	if err != nil {
		return nil, err
	}

	tData := domain.TokenData{
		UID:   claims["uid"].(string),
		Email: claims["email"].(string),
	}
	if ev, ok := claims["email_verified"].(bool); ok {
		tData.EmailVerified = ev
	}
	if dn, ok := claims["display_name"].(string); ok {
		tData.DisplayName = dn
	}
	if rl, ok := claims["roles"].(string); ok {
		tData.Roles = rl
	}
	return &tData, nil
}

// Countries
func (s *service) Countries() ([]*domain.Country, error) {
	res, err := s.repo.Countries(s.ctx)
	if err != nil {
		return nil, err
	}
	list := make([]*domain.Country, 0, len(res))
	for _, item := range res {
		list = append(list, &domain.Country{
			Id:      item.ID,
			Title:   item.Title,
			Code:    item.Code,
			IsoCode: item.IsoCode,
		})
	}
	return list, nil
}

// CountryAdd
func (s *service) CountryAdd(user *domain.TokenData, country domain.Country) (*domain.Country, error) {
	if country.Title == "" || country.Code == "" {
		return nil, fmt.Errorf("invalid country data")
	}
	if country.Id > 0 {
		return nil, fmt.Errorf("country already exists")
	}
	res, err := s.repo.CountryInsert(s.ctx, pgdb.CountryInsertParams{
		Title:   country.Title,
		Code:    country.Code,
		IsoCode: country.IsoCode,
	})
	if err != nil {
		return nil, err
	}
	return &domain.Country{
		Id:      res.ID,
		Title:   res.Title,
		Code:    res.Code,
		IsoCode: res.IsoCode,
	}, nil
}

// CountrySave
func (s *service) CountrySave(user *domain.TokenData, country domain.Country) (*domain.Country, error) {
	var (
		res *pgdb.Country
		err error
	)
	if country.Title == "" || country.Code == "" {
		return nil, fmt.Errorf("invalid country data")
	}
	if country.Id == 0 {
		// Insert country
		return s.CountryAdd(user, country)
	}

	res, err = s.repo.CountryUpdate(s.ctx, pgdb.CountryUpdateParams{
		ID:      country.Id,
		Title:   country.Title,
		Code:    country.Code,
		IsoCode: country.IsoCode,
	})
	if err != nil {
		return nil, err
	}
	return &domain.Country{
		Id:      res.ID,
		Title:   res.Title,
		Code:    res.Code,
		IsoCode: res.IsoCode,
	}, nil
}
