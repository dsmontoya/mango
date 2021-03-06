package aggregation

import (
	"github.com/dsmontoya/mango/bson"
	"github.com/dsmontoya/utils/maputils"
)

type Project bson.M

func NewProject() Project {
	return Project{}
}

//Array projects fields as elements in field.
func (p Project) Array(field string, fields ...string) Project {
	c := p.copy()
	array := make([]string, len(fields))
	for i, f := range fields {
		array[i] = "$" + f
	}
	c[field] = array
	return c
}

//Exclude specifies the exclusion of a field..
func (p Project) Exclude(field string) Project {
	c := p.copy()
	c[field] = 0
	return c
}

//Include specifies the inclusion of a field.
func (p Project) Include(field string) Project {
	c := p.copy()
	c[field] = 1
	return c
}

//Expression adds a new field or resets the value of an
//existing field.
func (p Project) Expression(field string, expression Expression) Project {
	c := p.copy()
	c[field] = expression.Apply()
	return c
}

//Rename modifies the name of a field.
func (p Project) Rename(old, new string) Project {
	c := p.copy()
	c[new] = "$" + old
	return c
}

func (p Project) copy() Project {
	c := Project{}
	maputils.Copy(p, c)
	return c
}
