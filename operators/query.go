package operators

import (
	"github.com/dsmontoya/mango/bson"
	"github.com/dsmontoya/utils/maputils"
)

type Query bson.M

func (q Query) Equal(field string, value interface{}) Query {
	c := q.copy()
	c[field] = value
	return c
}

//In selects the documents where the value of a field
//equals any value in the specified array. Calling In
//multiple times for the same field will append new
//values.
func (q Query) In(field string, values ...interface{}) Query {
	c := q.copy()
	if _, ok := c[field]; ok {
		m := c[field].(bson.M)
		v := m["$in"].([]interface{})
		v = append(v, values...)
		m["$in"] = v
	} else {
		c[field] = bson.M{"$in": values}
	}
	return c
}

func (q Query) copy() Query {
	c := Query{}
	maputils.Copy(q, c)
	return c
}
