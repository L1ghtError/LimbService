package service

import (
	"fmt"
	"light-backend/config"
	"light-backend/model"
	"light-backend/mongoclient"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx, user model.UserSchema) (*TokenPair, error) {
	collection := mongoclient.DB.Collection("userBase")
	filter := bson.D{{Key: "email", Value: user.Email}}
	// TODO: cout is a overkill for isExists use case
	alreadyExists, _ := collection.CountDocuments(c.Context(), filter)
	if alreadyExists != 0 {
		return nil, fmt.Errorf("user already registred")
	}

	user.Password, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.ActivationLink, _ = uuid.NewV7()

	inserted, err := collection.InsertOne(c.Context(), user)
	user.ID = inserted.InsertedID.(primitive.ObjectID)
	if err != nil {
		return nil, err
	}

	tokens, err := GenerateTokens(&user)
	if err != nil {
		return nil, err
	}

	err = SaveToken(c, &model.TokenSchema{UserId: user.ID, RefreshToken: tokens.Refresh})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func Login(c *fiber.Ctx, user model.UserSchema) (*TokenPair, error) {
	collection := mongoclient.DB.Collection("userBase")
	filter := bson.D{{Key: "email", Value: user.Email}}
	var dbUser model.UserSchema

	err := collection.FindOne(c.Context(), filter).Decode(&dbUser)
	if err == mongo.ErrNoDocuments {
		return nil, fiber.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, err
	}
	user.ID = dbUser.ID
	tokens, err := GenerateTokens(&user)
	if err != nil {
		return nil, err
	}

	err = SaveToken(c, &model.TokenSchema{UserId: user.ID, RefreshToken: tokens.Refresh})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func Logout(c *fiber.Ctx, user *UserClaims) error {
	userId, err := primitive.ObjectIDFromHex(user.UserId)
	if err != nil {
		return err
	}
	err = RemoveToken(c, &model.TokenSchema{UserId: userId})
	return err
}

func Refresh(c *fiber.Ctx, token *jwt.Token) (*TokenPair, error) {
	claims, err := ClaimModel(&token.Raw, []byte(config.Config("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(claims.UserId)
	if err != nil {
		return nil, err
	}

	dbToken, err := GetToken(c, &model.TokenSchema{UserId: userId})

	if err != nil {
		return nil, fiber.ErrUnauthorized
	} else if dbToken.RefreshToken != token.Raw {
		return nil, fiber.ErrUnauthorized
	}

	tokens, err := GenerateTokens(&model.UserSchema{ID: userId, Email: claims.Email})
	if err != nil {
		return nil, err
	}

	err = SaveToken(c, &model.TokenSchema{UserId: userId, RefreshToken: tokens.Refresh})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func GetBasics(c *fiber.Ctx, token *jwt.Token) (*model.UserSchema, error) {
	collection := mongoclient.DB.Collection("userBase")

	// in most cases err was validated in jwt middleware
	claims, _ := ClaimModel(&token.Raw, []byte(config.Config("JWT_ACCESS_SECRET")))
	fmt.Printf("id is: %v", claims.UserId)

	objectID, err := primitive.ObjectIDFromHex(claims.UserId)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	var dbUser model.UserSchema

	err = collection.FindOne(c.Context(), filter).Decode(&dbUser)
	if err == mongo.ErrNoDocuments {
		return nil, fiber.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	dbUser.Password = nil
	dbUser.ActivationLink = [16]byte{}
	return &dbUser, nil
}
