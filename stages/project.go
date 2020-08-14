package stages

import (
	"github.com/dsmontoya/mango"
	"github.com/dsmontoya/utils/maputils"
)

type Project mango.M

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

func (p Project) copy() Project {
	c := Project{}
	maputils.Copy(p, c)
	return c
}
