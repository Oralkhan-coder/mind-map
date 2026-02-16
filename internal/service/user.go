package service

import (
	"context"
	"errors"

	"github.com/Oralkhan-coder/mind-map/internal/core"
	"github.com/Oralkhan-coder/mind-map/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(collection *mongo.Collection) *UserService {
	return &UserService{collection}
}

func (srv *UserService) GetUserById(ctx context.Context, id string) (*model.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, core.BadRequest("invalid user id: " + err.Error())
	}

	user := new(model.User)
	err = srv.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, core.BadRequest("user does not exists")
		}
		return nil, core.InternalServerError("internal server error: " + err.Error())
	}

	return user, nil
}

func (srv *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, core.BadRequest("invalid email")
	}
	user := new(model.User)
	err := srv.collection.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, core.InternalServerError("internal server error: " + err.Error())
	}

	return user, nil
}

func (srv *UserService) CreateUser(ctx context.Context, user *model.User) (string, error) {
	res, err := srv.collection.InsertOne(ctx, user)
	if err != nil {
		return "", core.InternalServerError(err.Error())
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", core.InternalServerError("failed to get inserted ID")
	}

	return oid.Hex(), nil
}

func (srv *UserService) VerifyUser(ctx context.Context, userId string) error {
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return core.BadRequest("invalid user id: " + err.Error())
	}

	res, err := srv.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID, "is_verified": false},
		bson.M{"$set": bson.M{"is_verified": true}},
	)
	if err != nil {
		return core.InternalServerError(err.Error())
	}

	if res.MatchedCount == 0 {
		return core.NotFound("user not found or already verified")
	}
	return nil
}
