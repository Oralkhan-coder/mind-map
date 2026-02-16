package service

import (
	"context"
	"time"

	"github.com/Oralkhan-coder/mind-map/internal/core"
	"github.com/Oralkhan-coder/mind-map/internal/http/dto"
	"github.com/Oralkhan-coder/mind-map/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NodeService struct {
	nodeCollection *mongo.Collection
	userService    *UserService
	mapService     *MapService
}

func NewNodeService(nodeCollection *mongo.Collection, userService *UserService, mapService *MapService) *NodeService {
	return &NodeService{nodeCollection, userService, mapService}
}

func (srv *NodeService) GetByMapId(ctx context.Context, mapId, userId string) ([]*model.Node, error) {
	if _, err := srv.userService.GetUserById(ctx, userId); err != nil {
		return nil, err
	} else if _, err := srv.mapService.GetByID(ctx, mapId, userId); err != nil {
		return nil, err
	}

	cursor, err := srv.nodeCollection.Find(ctx, bson.M{"map_id": mapId})
	if err != nil {
		return nil, core.BadRequest(err.Error())
	}
	defer cursor.Close(ctx)

	var nodes []*model.Node
	for cursor.Next(ctx) {
		node := new(model.Node)
		if err := cursor.Decode(node); err != nil {
			return nil, core.BadRequest(err.Error())
		}
		nodes = append(nodes, node)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (srv *NodeService) CreateNode(ctx context.Context, req *dto.NodeCreateRequest, userId, mapId string) (string, error) {
	if _, err := srv.userService.GetUserById(ctx, userId); err != nil {
		return "", err
	} else if _, err := srv.mapService.GetByID(ctx, mapId, userId); err != nil {
		return "", err
	}

	newNode, err := model.NewNode(mapId, req.ParentID, req.Type, req.Data.Label,
		req.Position.X, req.Position.Y, req.Width, req.Height, time.Now(), nil, nil)
	res, err := srv.nodeCollection.InsertOne(ctx, newNode)
	if err != nil {
		return "", core.InternalServerError(err.Error())
	}

	return res.InsertedID.(string), nil
}
