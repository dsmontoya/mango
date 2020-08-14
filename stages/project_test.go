package stages

import (
	"reflect"
	"testing"
)

func TestProject_Array(t *testing.T) {
	type args struct {
		field  string
		fields []string
	}
	tests := []struct {
		name string
		p    Project
		args args
		want Project
	}{
		{
			"array",
			NewProject(),
			args{
				"list", []string{
					"x", "y", "z",
				},
			},
			Project{
				"list": []string{
					"$x", "$y", "$z",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Array(tt.args.field, tt.args.fields...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Project.Array() = %v, want %v", got, tt.want)
			}
		})
	}
}
