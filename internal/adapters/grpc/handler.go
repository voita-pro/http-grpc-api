package grpcadapter

import (
	"context"

	"github.com/voita-pro/http-grpc-api/internal/domain"
	"github.com/voita-pro/http-grpc-api/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedServicePBServer
	uc   domain.Usecase
	auth domain.Authenticator
}

func NewHandler(auth domain.Authenticator, uc domain.Usecase) *Handler {
	return &Handler{
		uc:   uc,
		auth: auth,
	}
}

// GetTokenData verify token from headers and return user data from token
func (s *Handler) GetTokenData(ctx context.Context) (*domain.TokenData, error) {
	idToken, err := bearerFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// Check token and get user claims
	validUser, err := s.auth.VerifyToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	user := domain.TokenData{
		UID:           validUser["uid"].(string),
		Email:         validUser["email"].(string),
		EmailVerified: validUser["email_verified"].(bool),
		DisplayName:   validUser["display_name"].(string),
		Role:          validUser["roles"].(string),
	}
	return &user, nil
}

// Login
func (s *Handler) Login(ctx context.Context, in *pb.LoginIN) (*pb.LoginOUT, error) {
	token, err := s.uc.SignIn(domain.Login{Email: in.Email, Password: in.Password})
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}
	return &pb.LoginOUT{
		Token: token,
	}, nil
}

// Countries
func (s *Handler) Countries(ctx context.Context, in *emptypb.Empty) (*pb.CountriesOUT, error) {
	data, err := s.uc.Countries()
	if err != nil {
		return nil, err
	}
	out := pb.CountriesOUT{Data: make([]*pb.Country, 0, len(data))}
	for _, item := range data {
		out.Data = append(out.Data, &pb.Country{
			Id:    item.Id,
			Title: item.Title,
			Code:  item.Code,
		})
	}
	return &out, nil
}

// CountryAdd
func (s *Handler) CountryAdd(ctx context.Context, in *pb.CountryIN) (*pb.Country, error) {
	user, err := s.GetTokenData(ctx)
	if err != nil {
		return nil, err
	}
	res, err := s.uc.CountrySave(user, domain.Country{
		Title:   in.Title,
		Code:    in.Code,
		IsoCode: &in.IsoCode,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Country{
		Id:      res.Id,
		Title:   res.Title,
		Code:    res.Code,
		IsoCode: *res.IsoCode,
	}, nil
}

// CountrySave
func (s *Handler) CountrySave(ctx context.Context, in *pb.Country) (*pb.Country, error) {
	user, err := s.GetTokenData(ctx)
	if err != nil {
		return nil, err
	}
	res, err := s.uc.CountrySave(user, domain.Country{
		Id:      in.Id,
		Title:   in.Title,
		Code:    in.Code,
		IsoCode: &in.IsoCode,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Country{
		Id:      res.Id,
		Title:   res.Title,
		Code:    res.Code,
		IsoCode: *res.IsoCode,
	}, nil
}
