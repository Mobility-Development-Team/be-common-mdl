package media

import (
	"fmt"
	"testing"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/spf13/viper"
)

func TestGetMediaBatches(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.media.module.url.base", "https://dev.api.fours.app/media/api")
	apis.Init(v)
	type args struct {
		tk      string
		batchId []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]model.MediaParam
		wantErr bool
	}{
		{
			name: "Test call",
			args: args{
				tk:      "", // Input token here to test (dev env)
				batchId: []string{"895d9fbb-e2cd-45bc-b9de-883a2ec0904b", "93b7804a-5c25-4101-8ea3-3571bce3e5fa"},
			},
			want:    map[string]model.MediaParam{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMediaBatches(tt.args.tk, tt.args.batchId...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMediaBatches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("GetMediaBatches() = ", got)
		})
	}
}

func TestGetMediaByRefId(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.media.module.url.base", "https://dev.api.fours.app/media/api")
	apis.Init(v)
	type args struct {
		tk    string
		refId []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test call",
			args: args{
				tk:    "", // Input token here to test
				refId: []string{"9edb9686-e8bf-4e88-819c-ee939798bf83"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMediaByRefId(tt.args.tk, tt.args.refId...)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestGetMediaByRefId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("TestGetMediaByRefId() = ", got)
		})
	}
}
