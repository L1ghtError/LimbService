package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TokenSchema struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId       primitive.ObjectID `json:"userId" bson:"userId"`
	RefreshToken string             `json:"refreshTokens" bson:"refreshToken"`
}
