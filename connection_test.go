package mango

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

func Test_portString(t *testing.T) {
	Convey("Given a valid port number", t, func() {
		var port uint = 3000
		Convey("When the port is parsed to string", func() {
			s := portString(port)
			Convey("The string should be valid", func() {
				So(s, ShouldEqual, ":3000")
			})
		})
	})

	Convey("Given a zero port number", t, func() {
		var port uint = 0
		Convey("When the port is parsed to string", func() {
			s := portString(port)
			Convey("The string should be empty", func() {
				So(s, ShouldBeEmpty)
			})
		})
	})
}

func Test_hostString(t *testing.T) {
	Convey("Given a address and port", t, func() {
		var port uint = 3000
		host := "example.com"

		Convey("When the host string is generated", func() {
			s := hostString(host, port)
			Convey("The string should be valid", func() {
				So(s, ShouldEqual, "example.com:3000")
			})
		})
	})
}

func TestConnect(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name         string
		args         args
		wantDatabase string
		wantErr      bool
	}{
		{
			"connect",
			args{
				&Config{
					ClientOptions: *new(mongooptions.ClientOptions).ApplyURI(os.Getenv("MONGO_URI")),
					Database:      "test",
					Context:       context.Background(),
				},
			},
			"test",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, err := Connect(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = connection.client.Ping(context.Background(), nil)
			if err != nil {
				t.Errorf("Unable to ping")
			}
		})
	}
}

func Test_setInsertValues_withDocument(t *testing.T) {
	type model struct {
		Document
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"model", args{&model{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := tt.args.v
			setInsertValues(value)
			v := reflect.ValueOf(value).Elem()
			fieldCreatedAt := v.FieldByName("CreatedAt").Interface().(time.Time)
			fieldUpdatedAt := v.FieldByName("UpdatedAt").Interface().(time.Time)
			fieldID := v.FieldByName("ID").Interface().(primitive.ObjectID)

			if fieldCreatedAt.IsZero() || fieldUpdatedAt.IsZero() {
				t.Errorf("zero value date in %v or %v", fieldCreatedAt, fieldUpdatedAt)
				return
			}
			if fieldID.IsZero() {
				t.Errorf("zero value id")
			}
		})
	}
}

func Test_setInsertValues_withoutDocument(t *testing.T) {
	type model struct {
		ID        primitive.ObjectID `bson:"_id"`
		CreatedAt time.Time          `bson:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt"`
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"model", args{&model{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := tt.args.v
			setInsertValues(value)
			v := reflect.ValueOf(value).Elem()
			fieldCreatedAt := v.FieldByName("CreatedAt").Interface().(time.Time)
			fieldUpdatedAt := v.FieldByName("UpdatedAt").Interface().(time.Time)
			fieldID := v.FieldByName("ID").Interface().(primitive.ObjectID)

			if fieldCreatedAt.IsZero() || fieldUpdatedAt.IsZero() {
				t.Errorf("zero value date in %v or %v", fieldCreatedAt, fieldUpdatedAt)
				return
			}
			if fieldID.IsZero() {
				t.Errorf("zero value id")
			}
		})
	}
}

func Test_setInsertValues_stringID(t *testing.T) {
	type model struct {
		ID string `bson:"_id"`
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"model string id", args{&model{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := tt.args.v
			setInsertValues(value)
			v := reflect.ValueOf(value).Elem()
			fieldID := v.FieldByName("ID").Interface().(string)

			if fieldID == "" {
				t.Errorf("id is empty")
				return
			}
		})
	}
}

func Test_setInsertValues_preexistentID(t *testing.T) {
	type model struct {
		ID string `bson:"_id"`
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"model string id", args{&model{"stringID"}}, "stringID"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := tt.args.v
			setInsertValues(value)
			v := reflect.ValueOf(value).Elem()
			fieldID := v.FieldByName("ID").Interface().(string)

			if fieldID != tt.want {
				t.Errorf("got %s, want %s", fieldID, tt.want)
				return
			}
		})
	}
}
