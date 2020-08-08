package mango

import (
	"testing"
)

type TestStruct struct {
	Document `bson:",inline"`
	A        string
	a        string
	Child    *TestStruct
	Structs  []*TestStruct `bson:",omitempty"`
	Strings  []string
}

func Test_getCollection(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"given a struct", args{&TestStruct{}}, "testStructs"},
		{"given an array of structs", args{[]*TestStruct{}}, "testStructs"},
		{"given a string", args{"users"}, "users"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCollection(tt.args.i); got != tt.want {
				t.Errorf("getCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}
