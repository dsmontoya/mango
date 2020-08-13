package operators

import (
	"reflect"
	"testing"
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
