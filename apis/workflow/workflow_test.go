package workflow

import (
	"testing"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/spf13/viper"
)

func TestDeleteWorkflow(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.workflow.module.url.base", "https://dev.api.fours.app/workflow/api")
	apis.Init(v)
	tk := ""
	type args struct {
		id intstring.IntString
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TEST",
			args: args{
				id: 191,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteWorkflow(tk, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteWorkflow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
