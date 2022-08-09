package inspection

import (
	"fmt"
	"testing"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/spf13/viper"
)

func TestFindUserPendingAppointments(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.inspection.module.url.base", "https://dev.api.fours.app/inspection/api")
	apis.Init(v)
	tk := "" // Dev token here
	type args struct {
		userRefKey string
		isSimple   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test not simple",
			args: args{
				userRefKey: "2157402f-ef03-4758-8e91-6b78ef98f2a7",
				isSimple:   false,
			},
			wantErr: false,
		},
		{
			name: "Test simple",
			args: args{
				userRefKey: "2157402f-ef03-4758-8e91-6b78ef98f2a7",
				isSimple:   true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindUserPendingAppointments(tk, tt.args.userRefKey, tt.args.isSimple)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserPendingAppointments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("FindUserPendingAppointments() isSimple = %t, got = %v\n", tt.args.isSimple, got)
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.inspection.module.url.base", "https://dev.api.fours.app/inspection/api")
	apis.Init(v)
	tk := "" // Dev token here
	type args struct {
		cri GetAllTasksCriteria
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test searchType, taskStatuses",
			args: args{
				cri: GetAllTasksCriteria{
					SearchType: "ASSIGNED",
					TaskStatuses: []string{
						"WORK_IN_PROGRESS",
						"FURTHER_FOLLOW_UP",
						"AWAITING_APPROVAL",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllTasks(tk, tt.args.cri)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GetAllTasks() got = %v\n", got)
		})
	}
}

func TestGetLatestTasksByParentRefIds(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.inspection.module.url.base", "https://dev.api.fours.app/inspection/api")
	apis.Init(v)
	tk := "" // Dev token here
	type args struct {
		taskParentRefIds []intstring.IntString
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				taskParentRefIds: []intstring.IntString{514},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestTasksByParentRefIds(tk, tt.args.taskParentRefIds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestTasksByParentRefIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GetLatestTasksByParentRefIds() got = %v\n", got)
		})
	}
}

func TestGetSitePlanBySiteWalkId(t *testing.T) {
	v := viper.New()
	v.Set("apis.internal.inspection.module.url.base", "https://dev.api.fours.app/inspection/api")
	apis.Init(v)
	tk := "" // Dev token here
	type args struct {
		siteWalkId intstring.IntString
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				siteWalkId: 520,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSitePlanBySiteWalkId(tk, tt.args.siteWalkId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSitePlanBySiteWalkId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GetSitePlanBySiteWalkId() got = %v\n", got)
		})
	}
}
