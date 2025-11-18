package firebaseadapter

import (
	"context"
	"net/http"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/voita-pro/http-grpc-api/internal/config"
	"google.golang.org/api/option"
)

// Client wraps Firebase Admin SDK auth client and Web API key powered REST flows.
type Client struct {
	Auth   *auth.Client
	APIKey string
	HTTP   *http.Client
}

// New creates a Firebase client using service account credentials file and optional API key.
func New(ctx context.Context, conf *config.Firebase) (*Client, error) {
	opts := []option.ClientOption{}
	if conf.CredsFile != "" {
		opts = append(opts, option.WithCredentialsFile(conf.CredsFile))
	}
	if conf.ProjectID != "" {
		opts = append(opts, option.WithQuotaProject(conf.ProjectID))
	}
	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, err
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		Auth:   authClient,
		APIKey: conf.WebAPIKey,
		HTTP:   &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// VerifyToken verifies an ID token and returns essential claims.
func (c *Client) VerifyToken(ctx context.Context, idToken string) (claims map[string]interface{}, err error) {
	tok, err := c.Auth.VerifyIDToken(ctx, idToken)

	if err != nil {
		return nil, err
	}
	claims["uid"] = tok.UID
	if em, ok := tok.Claims["email"].(string); ok {
		claims["email"] = em
	}
	if ev, ok := tok.Claims["email_verified"].(bool); ok {
		claims["email_verified"] = ev

	}
	if un, ok := tok.Claims["name"].(string); ok {
		claims["display_name"] = un
	}
	return claims, nil
}

func (c *Client) CustomTokenWithClaims(ctx context.Context, uid string, devClaims map[string]interface{}) (string, error) {
	return c.Auth.CustomTokenWithClaims(ctx, uid, devClaims)
}
