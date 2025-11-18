package domain

import "context"

type Usecase interface {
	SignIn(Login) (string, error)
	Countries() ([]*Country, error)
	CountrySave(*TokenData, Country) (*Country, error)
	VerifyToken(string) (*TokenData, error)
}

type Authenticator interface {
	VerifyToken(context.Context, string) (map[string]interface{}, error)
	CustomTokenWithClaims(context.Context, string, map[string]interface{}) (string, error)
}
