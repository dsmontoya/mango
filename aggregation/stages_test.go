package aggregation

import (
	"reflect"
	"testing"

	"github.com/dsmontoya/mango/bson"
	"github.com/dsmontoya/mango/operators"
)

func TestStages_Facet(t *testing.T) {
	type args struct {
		pipelines map[string]Stages
	}
	tests := []struct {
		name string
		s    Stages
		args args
		want Stages
	}{
		{
			"sample",
			New(),
			args{
				map[string]Stages{
					"field": New().Sample(3),
				},
			},
			Stages{
				{
					"$facet": map[string]Stages{
						"field": New().Sample(3),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Facet(tt.args.pipelines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stages.Facet() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		{"$in", Stages{}, args{operators.Query{}.In("name", "John")}, Stages{{"$match": operators.Query{"name": bson.M{"$in": []interface{}{"John"}}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Match(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stages.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStages_ReplaceRoot(t *testing.T) {
	type args struct {
		newRoot Expression
	}
	tests := []struct {
		name string
		s    Stages
		args args
		want Stages
	}{
		{
			"replaceRoot",
			New(),
			args{
				Field("name"),
			},
			Stages{
				bson.M{
					"$replaceRoot": bson.M{
						"newRoot": "$name",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ReplaceRoot(tt.args.newRoot); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stages.ReplaceRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStages_Unwind(t *testing.T) {
	type args struct {
		path                       string
		includeArrayIndex          string
		preserveNullAndEmptyArrays bool
	}
	tests := []struct {
		name string
		s    Stages
		args args
		want Stages
	}{
		{
			"unwind",
			New(),
			args{"list", "index", true},
			Stages{bson.M{"$unwind": bson.M{
				"path":                       "list",
				"includeArrayIndex":          "index",
				"preserveNullAndEmptyArrays": true,
			}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Unwind(tt.args.path, tt.args.includeArrayIndex, tt.args.preserveNullAndEmptyArrays); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stages.Unwind() = %v, want %v", got, tt.want)
			}
		})
	}
}
