package service

import (
	"context"
	"errors"
	"time"

	"github.com/Oralkhan-coder/mind-map/internal/core"
	"github.com/Oralkhan-coder/mind-map/internal/http/dto"
	"github.com/Oralkhan-coder/mind-map/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MapService struct {
	mapCollection  *mongo.Collection
	userCollection *mongo.Collection
}

func NewMapService(mapCollection, userCollection *mongo.Collection) *MapService {
	return &MapService{mapCollection, userCollection}
}

func (srv *MapService) GetMaps(ctx context.Context, userId string) ([]*model.Map, error) {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, core.BadRequest("invalid user id: " + err.Error())
	}

	cursor, err := srv.mapCollection.Find(ctx, bson.M{"user_id": oid, "deleted_at": nil})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var maps []*model.Map
	for cursor.Next(ctx) {
		m := new(model.Map)
		if err := cursor.Decode(m); err != nil {
			return nil, err
		}
		maps = append(maps, m)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return maps, nil
}

func (srv *MapService) GetByID(ctx context.Context, mapId, userId string) (*model.Map, error) {
	oid, err := primitive.ObjectIDFromHex(mapId)
	if err != nil {
		return nil, core.BadRequest("invalid map id" + err.Error())
	}

	ouserid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, core.BadRequest("invalid user id" + err.Error())
	}

	var result model.Map
	err = srv.mapCollection.FindOne(ctx, bson.M{"_id": oid, "deleted_at": nil}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, core.BadRequest("map not found" + err.Error())
		}
		return nil, core.InternalServerError(err.Error())
	} else if result.UserID != ouserid {
		return nil, core.BadRequest("user not owned by this map")
	}

	return &result, nil
}

func (srv *MapService) CreateMap(ctx context.Context, req *dto.MapCURequest, userId string) (string, error) {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", core.InternalServerError(err.Error())
	}
	var user model.User
	err = srv.userCollection.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&user)
	if err != nil {
		return "", core.BadRequest("user does not exist")
	} else if req.Title == "" {
		return "", core.BadRequest("title is required")
	}

	newMap, err := model.NewMap(req.Title, req.Description, userId, time.Now(), nil, nil)
	if err != nil {
		return "", core.BadRequest(err.Error())
	}

	res, err := srv.mapCollection.InsertOne(ctx, newMap)
	if err != nil {
		return "", core.InternalServerError(err.Error())
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", core.InternalServerError("failed to get inserted ID")
	}

	return oid.Hex(), nil
}

func (srv *MapService) UpdateMap(ctx context.Context, mapId, userId string, update *dto.MapCURequest) error {
	oid, err := primitive.ObjectIDFromHex(mapId)
	if err != nil {
		return core.BadRequest("invalid map id" + err.Error())
	}

	ouserid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return core.BadRequest("invalid user id" + err.Error())
	}

	m := new(model.Map)
	err = srv.mapCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(m)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.BadRequest("map not found" + err.Error())
		}
		return core.InternalServerError(err.Error())
	} else if m.UserID != ouserid {
		return core.BadRequest("user not owned by this map")
	} else if update.Title == "" {
		return core.BadRequest("title is required")
	}

	res, err := srv.mapCollection.UpdateOne(ctx, bson.D{{"_id", m.ID}}, bson.D{{"$set", update}})
	if err != nil {
		return core.InternalServerError(err.Error())
	} else if res.MatchedCount == 0 {
		return core.BadRequest("map not found")
	}

	return nil
}

func (srv *MapService) DeleteMap(ctx context.Context, mapId, userId string) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return core.BadRequest("invalid user id: " + err.Error())
	}
	var user model.User
	err = srv.userCollection.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&user)
	if err != nil {
		return core.BadRequest("user does not exist: " + err.Error())
	}

	oid, err = primitive.ObjectIDFromHex(mapId)
	if err != nil {
		return core.InternalServerError(err.Error())
	}

	var emap model.Map
	err = srv.mapCollection.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&emap)
	if err != nil {
		return core.BadRequest("map does not exist: " + err.Error())
	} else if emap.UserID != user.ID {
		return core.BadRequest("map does not belong to current user")
	}

	_, err = srv.mapCollection.UpdateOne(ctx, bson.D{{"_id", emap.ID}}, bson.D{{"$set", bson.D{{"deleted_at", time.Now()}}}})
	if err != nil {
		return core.InternalServerError(err.Error())
	}

	return nil
}
