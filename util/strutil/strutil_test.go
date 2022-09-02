package strutil

import (
	"reflect"
	"sort"
	"testing"
)

func TestSorterNumberFirst(t *testing.T) {
	type args struct {
		strs         []string
		numberDesc   bool
		strDesc      bool
		emptyStrLast bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test number descending string ascending empty last",
			args: args{
				strs: []string{
					"1049",
					"Apple",
					"",
					"Azzzz",
					"8088",
					"apple",
					"1047",
				},
				numberDesc:   true,
				strDesc:      false,
				emptyStrLast: true,
			},
			want: []string{
				"8088",
				"1049",
				"1047",
				"Apple",
				"Azzzz",
				"apple",
				"",
			},
		},
		{
			name: "Test number ascending string ascending empty last",
			args: args{
				strs: []string{
					"1049",
					"Apple",
					"",
					"Azzzz",
					"8088",
					"apple",
					"1047",
				},
				numberDesc:   false,
				strDesc:      false,
				emptyStrLast: true,
			},
			want: []string{
				"1047",
				"1049",
				"8088",
				"Apple",
				"Azzzz",
				"apple",
				"",
			},
		},
		{
			name: "Test number ascending string ascending empty last",
			args: args{
				strs: []string{
					"1049",
					"Apple",
					"",
					"Azzzz",
					"8088",
					"apple",
					"1047",
				},
				numberDesc:   false,
				strDesc:      false,
				emptyStrLast: false,
			},
			want: []string{
				"1047",
				"1049",
				"8088",
				"",
				"Apple",
				"Azzzz",
				"apple",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Slice(tt.args.strs, func(i, j int) bool {
				return CmpNumberFirst(tt.args.strs[i], tt.args.strs[j], tt.args.numberDesc, tt.args.strDesc, tt.args.emptyStrLast)
			})
			if !reflect.DeepEqual(tt.args.strs, tt.want) {
				t.Errorf("SorterNumberFirst() = %v, want %v", tt.args.strs, tt.want)
			}
		})
	}
}

func TestScreamCaseToLowerCamel(t *testing.T) {
	type args struct {
		scream string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "SCREAM_TEST",
			args: args{
				scream: "SCREAM_TEST",
			},
			want: "screamTest",
		},
		{
			name: "SCREAM",
			args: args{
				scream: "SCREAM",
			},
			want: "scream",
		},
		{
			name: "",
			args: args{
				scream: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScreamCaseToLowerCamel(tt.args.scream); got != tt.want {
				t.Errorf("ScreamCaseToLowerCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}
