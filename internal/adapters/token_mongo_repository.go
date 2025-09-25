package adapters

import (
	"context"
	"light-backend/internal/domain/token"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MongoTokenCollection = "tokenBase"

type TokenMongoRepository struct {
	collection mongo.Collection
}

type TokenSchema = token.TokenSchema

type tokenDoc struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"userId"`
	RefreshToken string             `bson:"refreshToken"`
}

func NewTokenMongoRepository(col mongo.Collection) *TokenMongoRepository {
	return &TokenMongoRepository{collection: col}
}

func (r TokenMongoRepository) GetToken(ctx context.Context, userId string) (TokenSchema, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return TokenSchema{}, err
	}
	filter := bson.D{{Key: "userId", Value: id}}

	var existingToken TokenSchema
	err = r.collection.FindOne(ctx, filter).Decode(&existingToken)
	if err == mongo.ErrNoDocuments {
		return TokenSchema{}, token.ErrNotFound
	} else if err != nil {
		return TokenSchema{}, err
	}
	return existingToken, nil
}

func (r TokenMongoRepository) RemoveToken(ctx context.Context, schema TokenSchema) error {
	id, err := primitive.ObjectIDFromHex(schema.UserId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "userId", Value: id}}

	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}

func (r TokenMongoRepository) SaveToken(ctx context.Context, schema TokenSchema) error {
	userID, err := primitive.ObjectIDFromHex(schema.UserId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "userId", Value: userID}}

	opts := options.FindOne().SetProjection(bson.M{"_id": 1})
	err = r.collection.FindOne(ctx, filter, opts).Err()

	// Check if the token exists or if there was an error
	if err == mongo.ErrNoDocuments {
		doc := bson.M{
			"userId":       userID,
			"refreshToken": schema.RefreshToken,
		}

		_, err := r.collection.InsertOne(ctx, doc)
		return err
	} else if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"refreshToken": schema.RefreshToken,
		},
	}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}
