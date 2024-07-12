package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeysize = 32

// JWTMaker is a JSON Web Token Maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker with the provided secretKey
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeysize {
		return nil, fmt.Errorf("invalid key size : must be %d characters", minSecretKeysize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new JSON web token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	// 創建新的Payload架構
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, nil
	}
	// 設定Token -> 由Header, Payload(Claims), Method(signature) 組成
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	/*
		return &Token{
			Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
			},
			Claims: claims,
			Method: method,
			}
	*/

	// 使用Key對Token進行簽名，return的是已經進行過簽名的token
	token, err := jwtToken.SignedString([]byte(maker.secretKey))

	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	return nil, nil
}

// 為了符合 Maker的interface，需要 CreateToken + VerifyToken的function
