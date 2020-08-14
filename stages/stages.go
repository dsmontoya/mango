package stages

import (
	"github.com/dsmontoya/mango"
	"github.com/dsmontoya/mango/operators"
)

type Stages []mango.M

//New initializes a new array of aggregation pipeline
//stages
func New() Stages {
	return make(Stages, 0)
}

func (s Stages) Match(query operators.Query) Stages {
	stage := mango.M{"$match": query}
	return append(s, stage)
}

func (s Stages) Sample(size int) Stages {
	stage := mango.M{"$sample": mango.M{"size": size}}
	return append(s, stage)
}
