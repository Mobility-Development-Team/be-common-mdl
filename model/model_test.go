package model

import (
	"fmt"
	"testing"
)

func TestSanitizeForCreate(t *testing.T) {
	type args struct {
		mdl interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TEST",
			args: args{
				mdl: &struct {
					Model
					AAA string
					bbb string
					CCC *Model
					DDD []*Model
					EEE []Model
				}{
					Model: Model{
						Id:               124,
						CreatedBy:        "SYSTEM",
						CreatedByDisplay: "SYSTEM",
						UpdatedByDisplay: "SYSTEM",
					},
					AAA: "aaaaaaaaaaaa",
					bbb: "bbbbbbbbbbbb",
					CCC: &Model{
						Id:               126,
						CreatedBy:        "SYSTEM",
						CreatedByDisplay: "SYSTEM",
						UpdatedByDisplay: "SYSTEM",
					},
					DDD: []*Model{
						{
							Id:               127,
							CreatedBy:        "SYSTEM",
							CreatedByDisplay: "SYSTEM",
							UpdatedByDisplay: "SYSTEM",
						},
						{
							Id:               128,
							CreatedBy:        "SYSTEM",
							CreatedByDisplay: "SYSTEM",
							UpdatedByDisplay: "SYSTEM",
						},
					},
					EEE: []Model{
						{
							Id:               129,
							CreatedBy:        "SYSTEM",
							CreatedByDisplay: "SYSTEM",
							UpdatedByDisplay: "SYSTEM",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("SanitizeForCreate() = %v", SanitizeForCreate(tt.args.mdl))
		})
	}
}
