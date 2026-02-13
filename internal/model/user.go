package model

import (
	"time"

	"github.com/Oralkhan-coder/mind-map/internal/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email"`
	Nickname   string             `bson:"nickname"`
	Password   string             `bson:"password"`
	IsVerified bool               `bson:"is_verified"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  *time.Time         `bson:"updated_at,omitempty"`
	DeletedAt  *time.Time         `bson:"deleted_at,omitempty"`
}

func NewUser(email, nickname, password string, isVerified bool, createdAt time.Time,
	updatedAt, deletedAt *time.Time) (*User, error) {
	if email == "" {
		return nil, core.BadRequest("email is required")
	} else if nickname == "" {
		return nil, core.BadRequest("nickname is required")
	} else if password == "" {
		return nil, core.BadRequest("password is required")
	}

	return &User{
		Nickname:   nickname,
		Email:      email,
		Password:   password,
		IsVerified: isVerified,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		DeletedAt:  deletedAt,
	}, nil
}
