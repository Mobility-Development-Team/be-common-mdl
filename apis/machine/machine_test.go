package machine

import (
	"fmt"
	"testing"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/spf13/viper"
)

func TestGetOnePlantPermit(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.machine.module.url.base", "https://dev.api.fours.app/machine/api")
	apis.Init(v)
	tk := ""
	type args struct {
		permitMasterId intstring.IntString
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TEST",
			args: args{
				permitMasterId: 80,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOnePlantPermit(tk, tt.args.permitMasterId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOnePlantPermit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("GetOnePlantPermit() = ", got)
		})
	}
}