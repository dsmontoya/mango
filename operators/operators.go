package operators

import (
	"github.com/dsmontoya/mango"
	"github.com/dsmontoya/utils/maputils"
)

type Query mango.M

func (q Query) Equal(field string, value interface{}) Query {
	c := q.copy()
	c[field] = value
	return c
}

func (q Query) copy() Query {
	c := Query{}
	maputils.Copy(q, c)
	return c
}
