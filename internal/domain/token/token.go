package token

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenSchema struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	UserId       string `json:"userId" bson:"userId"`
	RefreshToken string `json:"refreshTokens" bson:"refreshToken"`
}

type TokenPair struct {
	Access  string
	Refresh string
}

type UserClaims struct {
	Email  string `json:"email"`
	UserId string `json:"id"`
	jwt.RegisteredClaims
}

const AccessTokenExpires time.Duration = (time.Minute * 15)
const RefreshokenExpires time.Duration = ((time.Hour * 24) * 30)

func GenerateTokens(userId, email string) (*TokenPair, error) {
	if userId == "" || email == "" {
		return nil, fiber.ErrBadRequest
	}
	currentTime := time.Now()

	claims := UserClaims{
		Email:  email,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(AccessTokenExpires)),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	// claims for refresh token are the same, but it expires later
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(currentTime.Add(RefreshokenExpires))

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	refreshToken, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return &TokenPair{Access: accessToken, Refresh: refreshToken}, nil
}

func ClaimModel(tokenString string, secret []byte) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	} else {
		return nil, err
	}
}
