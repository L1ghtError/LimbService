package adapters

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"light-backend/internal/domain/user"
)

const MongoUserCollection = "userBase"

type UserMongoRepository struct {
	collection mongo.Collection
}

type UserSchema = user.UserSchema

func NewUserMongoRepository(col mongo.Collection) *UserMongoRepository {
	return &UserMongoRepository{collection: col}
}

func (r *UserMongoRepository) Register(ctx context.Context, schema UserSchema, skipPassword bool) (UserSchema, error) {
	filter := bson.D{{Key: "email", Value: schema.Email}}
	// TODO: CountDocuments is a overkill for isExists use case
	alreadyExists, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return UserSchema{}, err
	}
	if alreadyExists != 0 {
		return UserSchema{}, user.ErrAlreadyExists
	}
	if !skipPassword {
		schema.Password, _ = bcrypt.GenerateFromPassword([]byte(schema.Password), bcrypt.DefaultCost)
	}
	schema.ActivationLink, _ = uuid.NewV7()

	inserted, err := r.collection.InsertOne(ctx, schema)
	if err != nil {
		return UserSchema{}, err
	}
	schema.ID = inserted.InsertedID.(primitive.ObjectID).Hex()

	return schema, nil
}

func (r *UserMongoRepository) GetUserByEmail(ctx context.Context, email string) (UserSchema, error) {
	filter := bson.D{{Key: "email", Value: email}}
	var dbUser UserSchema

	err := r.collection.FindOne(ctx, filter).Decode(&dbUser)
	if err == mongo.ErrNoDocuments {
		return UserSchema{}, user.ErrNotFound
	} else if err != nil {
		return UserSchema{}, err
	}
	return dbUser, nil
}

func (r *UserMongoRepository) GetUserById(ctx context.Context, userId string) (UserSchema, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return UserSchema{}, err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	var dbUser UserSchema

	err = r.collection.FindOne(ctx, filter).Decode(&dbUser)
	if err == mongo.ErrNoDocuments {
		return UserSchema{}, user.ErrNotFound
	} else if err != nil {
		return UserSchema{}, err
	}
	return dbUser, nil
}
func (r *UserMongoRepository) AddImageId(ctx context.Context, schema UserSchema, imgId string) error {
	filter := bson.D{{Key: "_id", Value: schema.ID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "images", Value: imgId}}}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
