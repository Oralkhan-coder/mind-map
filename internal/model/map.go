package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Map struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt   *time.Time         `json:"updatedAt" bson:"updated_at"`
	DeletedAt   *time.Time         `json:"deletedAt" bson:"deleted_at"`
}

func NewMap(title, description, userId string, createdAt time.Time, updateAt, deleteAt *time.Time) (*Map, error) {
	if title == "" {
		return nil, errors.New("title is required")
	} else if userId == "" {
		return nil, errors.New("userId is required")
	}
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.New("wrong userId")
	}

	return &Map{
		Title:       title,
		Description: description,
		UserID:      objID,
		CreatedAt:   createdAt,
		UpdatedAt:   updateAt,
		DeletedAt:   deleteAt,
	}, nil
}
