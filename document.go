package mango

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type D bson.D
type M bson.M

type Document struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt,omitempty"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"deletedAt,omitempty"`
}
