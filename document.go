package mango

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type D bson.D
type M bson.M

type Document struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt *time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt *time.Time         `bson:"updatedAt json:"updatedAt"`
}

func (d *Document) setInsertValues() {
	d.setObjectID()
	t := time.Now()
	d.CreatedAt = &t
	d.UpdatedAt = &t
}

func (d *Document) setObjectID() {
	if d.ID == primitive.NilObjectID {
		d.ID = primitive.NewObjectID()
	}
}

func getDocument(iface interface{}) *Document {
	v := reflect.ValueOf(iface)
	if k := v.Kind(); k != reflect.Ptr {
		panic("should be a pointer")
	}
	el := v.Elem()
	docField := el.FieldByName("Document")
	if docField.IsZero() {
		return nil
	}
	return docField.Addr().Interface().(*Document)
}
