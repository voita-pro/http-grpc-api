package usecase

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/rs/zerolog/log"
	"github.com/voita-pro/http-grpc-api/internal/domain"
)

func (s *service) SignIn(login domain.Login) (string, error) {
	hash := md5.Sum([]byte(login.Email))
	md5Sum := hex.EncodeToString(hash[:])
	log.Info().Msgf("uid: %s", md5Sum)
	log.Info().Msgf("email: %s", login.Email)

	token, err := s.auth.CustomTokenWithClaims(s.ctx, md5Sum, map[string]interface{}{
		"email": login.Email,
		"uid":   md5Sum,
		"roles": "user",
	})
	log.Info().Msgf("token: %s", token)
	if err != nil {
		return "", err
	}
	return token, nil
}
