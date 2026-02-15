package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Point struct {
	X float64 `json:"x" bson:"x"`
	Y float64 `json:"y" bson:"y"`
}

type NodeData struct {
	Label string `json:"label" bson:"label"`
}

type Node struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	MapID     primitive.ObjectID `bson:"map_id"`
	ParentID  primitive.ObjectID `bson:"parent_id,omitempty"`
	Type      string             `bson:"type"`
	Position  Point              `bson:"position"`
	Data      NodeData           `bson:"data"`
	Width     float64            `bson:"width,omitempty"`
	Height    float64            `bson:"height,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at"`
}
