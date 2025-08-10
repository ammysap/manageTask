package auth

import (
	"crypto/ecdsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/aman/internal/logging"
)

type ECDSAJWTConfig struct {
	PrivateKey     *ecdsa.PrivateKey
	PublicKey      *ecdsa.PublicKey
	ExpirationTime time.Duration
}

var config ECDSAJWTConfig

func Verify(token string) (*jwt.RegisteredClaims, error) {
	return VerifyWithPublicKey(token, config.PublicKey)
}

func VerifyWithPublicKey(
	token string, publicKey *ecdsa.PublicKey,
) (*jwt.RegisteredClaims, error) {
	log := logging.Default()
	claims := &jwt.RegisteredClaims{}

	tkn, err := jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		},
	)
	if err != nil {
		log.Errorf("token: %s Parsing failed with %s\n", token, err)
		return nil, err
	}

	if !tkn.Valid {
		log.Errorf("token: %s not valid\n", token)
		return nil, errors.New("unauthorized")
	}

	return claims, nil
}