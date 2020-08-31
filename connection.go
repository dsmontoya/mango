package mango

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/dsmontoya/mango/aggregation"
	"github.com/dsmontoya/mango/options"
	"github.com/dsmontoya/utils/reflectutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	mongooptions.ClientOptions
	Database string
	Context  context.Context
}

type Connection struct {
	client         *mongo.Client
	collectionName string
	context        context.Context
	database       string
	model          interface{}
}

func Connect(config *Config) (*Connection, error) {
	if config.Context == nil {
		config.Context = context.Background()
	}
	client, err := mongo.Connect(config.Context, &config.ClientOptions)
	if err != nil {
		return nil, err
	}
	connection := &Connection{client: client, database: config.Database}
	return connection, nil
}

func (c *Connection) Aggregate(pipeline aggregation.Stages, value interface{}, opts ...*options.Aggregate) error {
	cursor, err := c.aggregateWithCursor(pipeline, opts...)
	if err != nil {
		return err
	}
	return cursor.All(c.context, value)
}

func (c *Connection) aggregateWithCursor(pipeline aggregation.Stages, opts ...*options.Aggregate) (*mongo.Cursor, error) {
	collection := c.collection(nil)
	aggregateOptions := make([]*mongooptions.AggregateOptions, len(opts))
	for i := 0; i < len(opts); i++ {
		opt := opts[i]
		aggregateOptions[i] = &opt.AggregateOptions
	}
	return collection.Aggregate(c.context, pipeline, aggregateOptions...)
}

func (c *Connection) Collection(i interface{}) *Connection {
	c2 := c.clone()
	c2.collectionName = getCollection(i)
	return c2
}

func (c *Connection) Context(ctx context.Context) *Connection {
	c2 := c.clone()
	c2.context = ctx
	return c2
}

func (c *Connection) DeleteOne(filter interface{}, opts ...*options.Delete) error {
	collection := c.collection(c.collectionName)
	deleteOptions := make([]*mongooptions.DeleteOptions, len(opts))
	for i := 0; i < len(opts); i++ {
		opt := opts[i]
		deleteOptions[i] = opt.DeleteOptions
	}
	_, err := collection.DeleteOne(c.context, filter, deleteOptions...)
	return err
}

func (c *Connection) DeleteMany(filter interface{}, opts ...*options.Delete) error {
	collection := c.collection(c.collectionName)
	deleteOptions := make([]*mongooptions.DeleteOptions, len(opts))
	for i := 0; i < len(opts); i++ {
		opt := opts[i]
		deleteOptions[i] = opt.DeleteOptions
	}
	_, err := collection.DeleteMany(c.context, filter, deleteOptions...)
	return err
}

func (c *Connection) Disconnect() error {
	return c.client.Disconnect(c.context)
}

func (c *Connection) Find(filter interface{}, value interface{}, opts ...*options.Find) error {
	cursor, err := c.FindWithCursor(c.context, filter, opts...)
	if err != nil {
		return err
	}
	return cursor.All(c.context, value)
}

func (c *Connection) FindWithCursor(filter interface{}, value interface{}, opts ...*options.Find) (*mongo.Cursor, error) {
	var findOptions []*mongooptions.FindOptions
	collection := c.collection(value)
	for _, op := range opts {
		findOptions = append(findOptions, &op.FindOptions)
	}
	return collection.Find(c.context, filter, findOptions...)
}

func (c *Connection) FindOne(filter interface{}, value interface{}, opts ...*options.Find) error {
	var findOneOptions []*mongooptions.FindOneOptions
	collection := c.collection(value)
	for i := 0; i < len(opts); i++ {
		opt := opts[i]
		findOneOptions[i] = &mongooptions.FindOneOptions{
			AllowPartialResults: opt.AllowPartialResults,
			BatchSize:           opt.BatchSize,
			Collation:           opt.Collation,
			Comment:             opt.Comment,
			CursorType:          opt.CursorType,
			Hint:                opt.Hint,
			Max:                 opt.Max,
			MaxAwaitTime:        opt.MaxAwaitTime,
			MaxTime:             opt.MaxTime,
			Min:                 opt.Min,
			NoCursorTimeout:     opt.NoCursorTimeout,
			OplogReplay:         opt.OplogReplay,
			Projection:          opt.Projection,
			ReturnKey:           opt.ReturnKey,
			ShowRecordID:        opt.ShowRecordID,
			Skip:                opt.Skip,
			Snapshot:            opt.Snapshot,
			Sort:                opt.Sort,
		}
	}
	singleResult := collection.FindOne(c.context, filter, findOneOptions...)
	if err := singleResult.Err(); err != nil {
		return err
	}
	return singleResult.Decode(value)
}

func (c *Connection) InsertMany(values interface{}, opts ...*options.Insert) error {
	var insertValues []interface{}
	collection := c.collection(values)
	n := reflectutils.Each(values, func(i int, v reflect.Value) bool {
		vi := v.Interface()
		setInsertValues(vi)
		insertValues = append(insertValues, vi)
		return true
	})
	if n == 0 {
		return errors.New("you must provide at least one item")
	}
	_, err := collection.InsertMany(c.context, insertValues)
	// id := result.InsertedID.(primitive.ObjectID)
	// doc.ID = id
	return err
}

func (c *Connection) InsertOne(value interface{}, opts ...*options.Insert) error {
	insertOneOptions := make([]*mongooptions.InsertOneOptions, len(opts))
	setInsertValues(value)
	collection := c.collection(value)
	for i := 0; i < len(opts); i++ {
		opt := opts[i]
		insertOneOptions[i] = &mongooptions.InsertOneOptions{
			BypassDocumentValidation: opt.BypassDocumentValidation,
		}
	}
	_, err := collection.InsertOne(c.context, value, insertOneOptions...)
	return err
}

func (c *Connection) Model(value interface{}) *Connection {
	c2 := c.clone()
	c2.model = value
	return c2
}

func (c *Connection) UpdateOne(filter interface{}, update interface{}, opts ...*options.Update) error {
	updateOptions := make([]*mongooptions.UpdateOptions, len(opts))
	collection := c.collection(c.collectionName)
	for i := 0; i < len(opts); i++ {
		opt := opts[i]
		updateOptions[i] = &opt.UpdateOptions
	}
	_, err := collection.UpdateOne(c.context, filter, update, updateOptions...)
	return err
}

func (c *Connection) GetClient() *mongo.Client {
	return c.client
}

func (c *Connection) clone() *Connection {
	return &Connection{
		client:   c.client,
		database: c.database,
		model:    c.model,
		context:  c.context,
	}
}

func (c *Connection) collection(model interface{}) *mongo.Collection {
	client := c.client
	database := client.Database(c.database)
	if c.collectionName != "" {
		return database.Collection(c.collectionName)
	}
	return database.Collection(getCollection(model))
}

func convertCollation(collation *options.Collation) *mongooptions.Collation {
	if collation != nil {
		return &mongooptions.Collation{
			Locale: collation.Locale,
		}
	}
	return nil
}

func convertCursorType(cursorType *options.CursorType) *mongooptions.CursorType {
	if cursorType != nil {
		result := mongooptions.CursorType(*cursorType)
		return &result
	}
	return nil
}

func setInsertValues(value interface{}) {
	now := time.Now()
	v := reflectutils.DeepValue(reflect.ValueOf(value))
	k := v.Kind()
	if k != reflect.Struct {
		return
	}
	if document := v.FieldByName("Document"); document.IsValid() {
		doc := &Document{
			ID: document.FieldByName("ID").Interface().(primitive.ObjectID),
		}
		setInsertValues(doc)
		document.Set(reflect.Indirect(reflect.ValueOf(doc)))
	}
	t := v.Type()
	n := t.NumField()
	for i := 0; i < n; i++ {
		sf := t.Field(i)
		name := sf.Name
		bsonKey := sf.Tag.Get("bson")
		switch bsonKey {
		case "-":
			continue
		case "createdAt", "updatedAt":
			reflectutils.SetField(value, name, now)
		case "_id":
			setObjectID(value, name)
		}
	}
}

func setObjectID(value interface{}, name string) {
	var x interface{}
	v := reflectutils.DeepValue(reflect.ValueOf(value))
	field := v.FieldByName(name)
	if !field.IsValid() {
		return
	}
	id := primitive.NewObjectID()
	fieldType := field.Type().String()
	if fieldType == "primitive.ObjectID" {
		fieldPrimitive := field.Interface().(primitive.ObjectID)
		if fieldPrimitive != primitive.NilObjectID {
			return
		}
		x = id
	} else if fieldType == "string" {
		fieldString := field.Interface().(string)
		if fieldString != "" {
			return
		}
		x = id.Hex()
	}
	reflectutils.SetField(value, name, x)
}

func hostString(address string, port uint) string {
	return fmt.Sprintf("%s%s", address, portString(port))
}

func portString(port uint) string {
	s := ":%d"
	if port > 0 {
		return fmt.Sprintf(s, port)
	}
	return ""
}
