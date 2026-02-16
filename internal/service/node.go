package service

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type NodeService struct {
	nodeCollection *mongo.Collection
}

func NewNodeService(nodeCollection *mongo.Collection) *NodeService {
	return &NodeService{nodeCollection}
}
