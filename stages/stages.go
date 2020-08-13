package stages

import "github.com/dsmontoya/mango"

type Stages []mango.M

//New initializes a new array of aggregation pipeline
//stages
func New() Stages {
	return make(Stages, 0)
}

func (s Stages) Match() Stages {

}

func (s Stages) Sample(size int) Stages {
	stage := mango.M{"$sample": mango.M{"size": size}}
	return append(s, stage)
}
