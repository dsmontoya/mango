package aggregation

import (
	"reflect"
	"testing"

	"github.com/dsmontoya/mango"
)

func TestSetUnion(t *testing.T) {
	type args struct {
		expressions []Expression
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"setUnion", args{[]Expression{Field("name"), Var("test")}}, mango.M{"$setUnion": []interface{}{"$name", "$$test"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetUnion(tt.args.expressions).Apply(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetUnion() = %v, want %v", got, tt.want)
			}
		})
	}
}
