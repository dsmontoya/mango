package mango

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/dsmontoya/mango/aggregation"
	"github.com/dsmontoya/mango/options"
	"github.com/dsmontoya/utils/strutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

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

func TestConnection_Aggregate(t *testing.T) {
	type args struct {
		pipeline aggregation.Stages
		value    interface{}
		opts     []*options.Aggregate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantLen int
	}{
		{"sample", args{aggregation.New().Sample(3), &[]bson.M{}, nil}, false, 3},
	}
	connection := newConnection(context.Background()).
		Collection("aggregationTests")
	defer connection.Disconnect()
	ms := make([]bson.M, 100)
	a := strutils.Rand(10)
	for i := 0; i < 100; i++ {
		ms[i] = bson.M{"a": a}
	}
	if err := connection.InsertMany(ms); err != nil {
		panic(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := connection.Aggregate(tt.args.pipeline, tt.args.value, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Connection.Aggregate() error = %v, wantErr %v", err, tt.wantErr)
			}
			v := *(tt.args.value.(*[]bson.M))
			if len(v) != tt.wantLen {
				t.Errorf("len = %v, wantLen %v", len(v), tt.wantLen)
			}
			if v[0]["a"] != a {
				t.Errorf("v[0]['a'] = %v, a %v", v[0]["a"], a)
			}
		})
	}
	if err := connection.DeleteMany(bson.M{"a": a}); err != nil {
		panic(err)
	}
}

func newConnection(ctx context.Context) *Connection {
	config := &Config{
		ClientOptions: *new(mongooptions.ClientOptions).ApplyURI(os.Getenv("MONGO_URI")),
		Database:      "test",
		Context:       ctx,
	}
	connection, err := Connect(config)
	if err != nil {
		panic(err)
	}
	return connection
}
