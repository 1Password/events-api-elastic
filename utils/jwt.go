package utils

import (
	"errors"
	"fmt"

	"gopkg.in/square/go-jose.v2/jwt"
)

type Features []string

type JWTClaims struct {
	Audience []string `json:"aud"`
	Features Features `json:"1password.com/fts"`
}

const AudienceDEPRECATED = "com.1password.streamingservice"

const ItemUsageFeatureScope = "itemusages"
const SignInAttemptsFeatureScope = "signinattempts"

func ParseJWTClaims(token string) (*JWTClaims, error) {
	t, err := jwt.ParseSigned(token)

	if err != nil {
		return nil, err
	}

	claims := &JWTClaims{}

	// We don't have the ECDSA Pub Key for verification at this point
	// as we need the URL to retrieve the metadata
	// This is fine as the server will properly verify the token
	err = t.UnsafeClaimsWithoutVerification(claims)

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (t *JWTClaims) GetEventsURL() (string, error) {
	if t.Audience[0] == AudienceDEPRECATED {
		return "", errors.New("token does not have a url")
	}

	return fmt.Sprintf("https://%s", t.Audience[0]), nil
}

func (s Features) Contains(v string) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}
