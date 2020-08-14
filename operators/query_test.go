package operators

import (
	"reflect"
	"testing"

	"github.com/dsmontoya/mango"
)

func TestQuery_Equal(t *testing.T) {
	t.Run("copy", func(t *testing.T) {
		q := Query{}
		got := q.Equal("name", "John").Equal("age", 26)
		want := Query{"name": "John", "age": 26}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Query.Equal() = %v, want %v", got, want)
		}
		if q["name"] == "John" {
			t.Errorf("Query was modified")
		}
	})
}

func TestQuery_In(t *testing.T) {
	type args struct {
		field  string
		values []interface{}
	}
	tests := []struct {
		name string
		q    Query
		args args
		want Query
	}{
		{
			"empty",
			Query{},
			args{"name", []interface{}{"John", "Diana"}},
			Query{
				"name": mango.M{
					"$in": []interface{}{"John", "Diana"},
				},
			},
		},
		{
			"non-empty",
			Query{
				"name": mango.M{
					"$in": []interface{}{"John", "Peter"},
				},
			},
			args{"name", []interface{}{"Diana", "Hanna"}},
			Query{
				"name": mango.M{
					"$in": []interface{}{"John", "Peter", "Diana", "Hanna"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.In(tt.args.field, tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.In() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_copy(t *testing.T) {
	t.Run("copy", func(t *testing.T) {
		q := Query{"a": 1}
		got := q.copy()
		got["b"] = 2
		want := Query{"a": 1}
		if !reflect.DeepEqual(q, want) {
			t.Errorf("Query.copy() = %v, want %v", got, want)
		}
	})
}
