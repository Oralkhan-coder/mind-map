package model

import (
	"errors"
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

func NewNode(mapId, parentId, types, label string, x, y, w, h float64,
	createdAt time.Time, updatedAt, deletedAt *time.Time) (*Node, error) {
	if mapId == "" {
		return nil, errors.New("map id is required")
	} else if types == "" {
		return nil, errors.New("types is required")
	} else if w <= 0.0 {
		return nil, errors.New("x is required")
	} else if h <= 0.0 {
		return nil, errors.New("y is required")
	}
	position := Point{X: x, Y: y}
	data := NodeData{Label: label}

	mapOId, err := primitive.ObjectIDFromHex(mapId)
	if err != nil {
		return nil, errors.New("invalid map id: " + mapId)
	}
	var parentOId primitive.ObjectID
	if parentId != "" {
		parentOId, err = primitive.ObjectIDFromHex(parentId)
		if err != nil {
			return nil, errors.New("invalid parent id: " + parentId)
		}
	}

	return &Node{
		MapID:     mapOId,
		ParentID:  parentOId,
		Type:      types,
		Position:  position,
		Data:      data,
		Width:     w,
		Height:    h,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}, nil
}
