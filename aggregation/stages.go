package aggregation

import (
	"github.com/dsmontoya/mango/bson"
	"github.com/dsmontoya/mango/operators"
)

type Stages []bson.M

//New initializes a new array of aggregation pipeline
//stages
func New() Stages {
	return make(Stages, 0)
}

//Facet processes multiple aggregation pipelines within a single stage. Each sub-pipeline has its own field in the output document.
func (s Stages) Facet(pipelines map[string]Stages) Stages {
	stage := bson.M{"$facet": pipelines}
	return append(s, stage)
}

func (s Stages) Match(query operators.Query) Stages {
	stage := bson.M{"$match": query}
	return append(s, stage)
}

func (s Stages) Sample(size int) Stages {
	stage := bson.M{"$sample": bson.M{"size": size}}
	return append(s, stage)
}
