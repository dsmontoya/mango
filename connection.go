package mango

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/entropyx/mango/options"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Config struct {
	options.Client
	Address  string
	Port     uint
	Timeout  *time.Time
	Database string
	Username string
	Password string
	Source   string
	Context  context.Context
}

type Connection struct {
	client   *mongo.Client
	database string
}

func Connect(config *Config) (*Connection, error) {
	u := url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(config.Username, config.Password),
		Host:   fmt.Sprintf("%s:%d", config.Address, config.Port),
	}
	if config.Context == nil {
		config.Context = context.Background()
	}
	client, err := mongo.Connect(config.Context, u.String(), &config.ClientOptions)
	if err != nil {
		return nil, err
	}
	connection := &Connection{client: client, database: config.Database}
	return connection, nil
}

func (c *Connection) Register(ctx context.Context, models ...interface{}) {
	newCtx := context.WithValue(ctx, "connection", c)
	for _, model := range models {
		doc := getDocument(model)
		doc.Context = newCtx
	}
}

func (c *Connection) collection(model interface{}) *mongo.Collection {
	client := c.client
	return client.Database(c.database).Collection(getCollection(model))
}