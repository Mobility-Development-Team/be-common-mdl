package user

import (
	"fmt"
	"testing"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/spf13/viper"
)

func TestGetAllGroupInfo(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.user.module.url.base", "https://dev.api.fours.app/user/api")
	apis.Init(v)
	tk := ""
	type args struct {
		body map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TEST",
			args: args{
				body: map[string]interface{}{
					"contractId":    38,
					"isSystemGroup": true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllGroupInfo(tk, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllGroupInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GetAllGroupInfo() = %v\n", got)
		})
	}
}

func TestGetUserSignatures(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.user.module.url.base", "https://dev.api.fours.app/user/api")
	apis.Init(v)
	tk := ""
	type args struct {
		ids []intstring.IntString
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TEST",
			args: args{
				ids: []intstring.IntString{12, 179, 119},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserSignatures(tk, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserSignatures() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GetUserSignatures() = %v\n", got)
		})
	}
}
