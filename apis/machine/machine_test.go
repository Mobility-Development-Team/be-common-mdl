package machine

import (
	"fmt"
	"testing"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/Mobility-Development-Team/be-common-mdl/util/strutil"
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

func TestGetOneAsset(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.machine.module.url.base", "https://dev.api.fours.app/machine/api")
	apis.Init(v)
	tk := ""
	type args struct {
		criteria Equipment
		isSimple bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TEST",
			args: args{
				criteria: Equipment{
					SerialNo: strutil.NewPtr("S1234"),
				},
				isSimple: true,
			},
			wantErr: false,
		},
		{
			name: "TEST",
			args: args{
				criteria: Equipment{
					SerialNo: strutil.NewPtr("S1234"),
				},
				isSimple: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOneAsset(tk, tt.args.criteria, tt.args.isSimple)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOneAsset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("GetOneAsset() = ", got)
		})
	}
}

func TestGetAllPermits(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.machine.module.url.base", "https://dev.api.fours.app/machine/api")
	apis.Init(v)
	tk := ""
	type args struct {
		userRefKey   string
		criteria     PermitCriteria
		opt          GetAllPermitOps
		preloadNames []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TEST",
			args: args{
				userRefKey: "2157402f-ef03-4758-8e91-6b78ef98f2a7",
				criteria: PermitCriteria{
					SearchType: PermitSearchActionRequired,
				},
				opt: GetAllPermitOps{
					GetApprovalStage: false,
				},
				preloadNames: []string{
					PreloadApprovalFlows,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllPermits(tk, tt.args.userRefKey, tt.args.criteria, tt.args.opt, tt.args.preloadNames...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPermits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GetAllPermits() = %v\n", got)
		})
	}
}

func TestGetOneNCAPermit(t *testing.T) {
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
				permitMasterId: 128,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOneNCAPermit(tk, tt.args.permitMasterId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOneNCAPermit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("GetOneNCAPermit() = ", got)
		})
	}
}

func TestGetOnePermitToDig(t *testing.T) {
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
				permitMasterId: 502,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOnePermitToDig(tk, tt.args.permitMasterId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOnePermitToDig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("GetOnePermitToDig() = ", got)
		})
	}
}

func TestGetOneELPermit(t *testing.T) {
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
				permitMasterId: 649,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOneELPermit(tk, tt.args.permitMasterId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOneELPermit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("GetOneELPermit() = ", got)
		})
	}
}
