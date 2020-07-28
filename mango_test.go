package mango

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestStruct struct {
	Document `bson:",inline"`
	A        string
	a        string
	Child    *TestStruct
	Structs  []*TestStruct `bson:",omitempty"`
	Strings  []string
}

func Test_toBsonDoc(t *testing.T) {
	Convey("Given a model", t, func() {
		s := &TestStruct{A: "1", a: "2", Structs: []*TestStruct{
			&TestStruct{A: "3"},
		}, Strings: []string{
			"A", "B",
		}}
		Convey("When it is converted to a document", func() {
			doc := toBsonDoc(s)

			Convey("It should be valid", func() {
				So(doc, ShouldResemble, primitive.D{primitive.E{Key: "a", Value: "1"}, primitive.E{Key: "child", Value: primitive.D{}}, primitive.E{Key: "structs", Value: primitive.A{primitive.D{primitive.E{Key: "a", Value: "3"}, primitive.E{Key: "child", Value: primitive.D{}}, primitive.E{Key: "structs", Value: primitive.A{}}, primitive.E{Key: "strings", Value: primitive.A{}}}}}, primitive.E{Key: "strings", Value: primitive.A{"A", "B"}}})
			})
		})
	})
}

func Test_getCollection(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"given a struct", args{&TestStruct{}}, "testStruct"},
		{"given a string", args{"users"}, "users"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCollection(tt.args.i); got != tt.want {
				t.Errorf("getCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}
