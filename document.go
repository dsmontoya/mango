package mango

import (
	"time"

	"github.com/dsmontoya/mango/bson"
)

type Document struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt,omitempty"`
	DeletedAt *time.Time    `bson:"deletedAt" json:"deletedAt,omitempty"`
}
