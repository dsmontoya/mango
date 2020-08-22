package aggregation

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

func TestProject_Rename(t *testing.T) {
	type args struct {
		old string
		new string
	}
	tests := []struct {
		name string
		p    Project
		args args
		want Project
	}{
		{"rename", NewProject(), args{"old", "new"}, Project{"new": "$old"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Rename(tt.args.old, tt.args.new); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Project.Rename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProject_Expression(t *testing.T) {
	type args struct {
		field      string
		expression Expression
	}
	tests := []struct {
		name string
		p    Project
		args args
		want Project
	}{
		{"adf", NewProject(), args{"name", Field("userName")}, Project{"name": "$userName"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Expression(tt.args.field, tt.args.expression); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Project.Expression() = %v, want %v", got, tt.want)
			}
		})
	}
}
