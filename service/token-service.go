package service

import (
	"light-backend/config"
	"light-backend/model"
	"light-backend/mongoclient"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func GenerateTokens(user *model.UserSchema) (*TokenPair, error) {
	if user.Email == "" || user.ID.IsZero() {
		return nil, fiber.ErrBadRequest
	}
	currentTime := time.Now()

	claims := UserClaims{
		Email:  user.Email,
		UserId: user.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(AccessTokenExpires)),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(config.Config("JWT_ACCESS_SECRET")))
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	// claims for refresh token are the same, but it expires later
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(currentTime.Add(RefreshokenExpires))

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(config.Config("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return &TokenPair{Access: accessToken, Refresh: refreshToken}, nil
}

func ClaimModel(tokenString *string, secret []byte) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(*tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
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

func SaveToken(c *fiber.Ctx, token *model.TokenSchema) error {
	collection := mongoclient.DB.Collection("tokenBase")
	filter := bson.D{{Key: "userId", Value: token.UserId}}
	var existingToken model.TokenSchema
	err := collection.FindOne(c.Context(), filter).Decode(&existingToken)

	// Check if the token exists or if there was an error
	if err == mongo.ErrNoDocuments {
		_, err := collection.InsertOne(c.Context(), token)
		return err
	} else if err != nil {
		return err
	}

	_, err = collection.ReplaceOne(c.Context(), filter, token)
	return err
}

func GetToken(c *fiber.Ctx, token *model.TokenSchema) (*model.TokenSchema, error) {
	collection := mongoclient.DB.Collection("tokenBase")
	filter := bson.D{{Key: "userId", Value: token.UserId}}
	var existingToken model.TokenSchema
	err := collection.FindOne(c.Context(), filter).Decode(&existingToken)
	if err != nil {
		return nil, err
	}
	return &existingToken, nil
}

func RemoveToken(c *fiber.Ctx, token *model.TokenSchema) error {
	collection := mongoclient.DB.Collection("tokenBase")
	filter := bson.D{{Key: "userId", Value: token.UserId}}

	_, err := collection.DeleteOne(c.Context(), filter)
	return err
}
