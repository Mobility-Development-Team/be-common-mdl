package system

import (
	"fmt"
	"testing"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/spf13/viper"
)

func TestGetContractUserByUids(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.system.module.url.base", "https://dev.api.fours.app/system/api")
	apis.Init(v)
	tk := "" // Dev token here
	type args struct {
		contractId intstring.IntString
		uids       []intstring.IntString
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TEST",
			args: args{
				contractId: 38,
				uids:       []intstring.IntString{12},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractUserByUids(tk, tt.args.contractId, tt.args.uids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractUserByUids() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GetContractUserByUids() = %v\n", got)
		})
	}
}
