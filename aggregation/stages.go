package aggregation

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

//Facet processes multiple aggregation pipelines within a single stage. Each sub-pipeline has its own field in the output document.
func (s Stages) Facet(pipelines map[string]Stages) Stages {
	stage := mango.M{"$facet": pipelines}
	return append(s, stage)
}

func (s Stages) Match(query operators.Query) Stages {
	stage := mango.M{"$match": query}
	return append(s, stage)
}

func (s Stages) Sample(size int) Stages {
	stage := mango.M{"$sample": mango.M{"size": size}}
	return append(s, stage)
}
