package mango

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_getDocument(t *testing.T) {
	Convey("Given a struct with a reference to a document", t, func() {
		s := &TestStruct{}

		Convey("When the document is got", func() {
			doc := getDocument(s)

			Convey("doc should be valid", func() {
				So(doc, ShouldNotBeNil)
			})

			Convey("And the id is modified", func() {
				doc.ID = primitive.NewObjectID()

				Convey("The new context should be valid", func() {
					newDoc := getDocument(s)
					So(newDoc.ID, ShouldNotEqual, primitive.NilObjectID)
				})
			})
		})
	})
}

func TestDocument_setObjectID(t *testing.T) {
	type fields struct {
		ID        primitive.ObjectID
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"document without _id", fields{}},
		{"document with _id", fields{ID: primitive.NewObjectID()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Document{
				ID:        tt.fields.ID,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			d.setObjectID()
			if got := tt.fields.ID; got == primitive.NilObjectID {
				t.Errorf("_id = %v", got)
			}
		})
	}
}
