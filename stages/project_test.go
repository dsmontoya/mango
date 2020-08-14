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

func TestProject_Exclude(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		p    Project
		args args
		want Project
	}{
		{"exclude", NewProject(), args{"name"}, Project{"name": 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Exclude(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Project.Exclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProject_Include(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		p    Project
		args args
		want Project
	}{
		{"include", NewProject(), args{"name"}, Project{"name": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Include(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Project.Include() = %v, want %v", got, tt.want)
			}
		})
	}
}
