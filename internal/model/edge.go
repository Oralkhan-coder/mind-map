package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Edge struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	MapID primitive.ObjectID `bson:"map_id"`

	From string `bson:"from"`
	To   string `bson:"to"`
}
