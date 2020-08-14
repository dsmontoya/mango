package stages

import (
	"reflect"
	"testing"

	"github.com/dsmontoya/mango"
	"github.com/dsmontoya/mango/operators"
)

func TestStages_Match(t *testing.T) {
	type args struct {
		query operators.Query
	}
	tests := []struct {
		name string
		s    Stages
		args args
		want Stages
	}{
		{"$in", Stages{}, args{operators.Query{}.In("name", "John")}, Stages{{"$match": operators.Query{"name": mango.M{"$in": []interface{}{"John"}}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Match(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stages.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
