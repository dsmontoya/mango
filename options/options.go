package options

import (
	"github.com/dsmontoya/mango/aggregation"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CursorType int8

// Aggregate represents options that can be used to configure an Aggregate operation.
type Aggregate struct {
	options.AggregateOptions
}

// Delete represents options that can be used to configure DeleteOne and DeleteMany operations.
type Delete struct {
	*options.DeleteOptions
}

// Find represent all possible options to the Find() function.
type Find struct {
	options.FindOptions
	Projection aggregation.Project
}

//Insert represents all possible options to the InsertMany() and InsertOne() functions.
type Insert struct {
	BypassDocumentValidation *bool // If true, allows the write to opt-out of document level validation
	Ordered                  *bool // If true, when InsertMany() fails, return without performing the remaining inserts. Defaults to true.
}

type Update struct {
	options.UpdateOptions
}

type Client struct {
	options.ClientOptions
}

type Collation struct {
	options.Collation
}
