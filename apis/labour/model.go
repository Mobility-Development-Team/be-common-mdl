package labour

import (
	"time"

	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/model/pagination"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
)

type (
	UnsafeCase struct {
		model.Model
		CaseNo        string               `json:"caseNo"`
		CaseStatus    string               `json:"caseStatus"`
		CaseType      string               `json:"caseType"`
		ContractRefId intstring.IntString  `json:"contractId"`
		Worker        interface{}          `json:"worker"`
		Item          interface{}          `json:"item"`
		ApprovalFlow  []UnsafeCaseApproval `json:"approvalFlow"`
	}
	UnsafeCaseCriteria struct {
		CaseStatuses []string `json:"CaseStatuses"`
	}
	UnsafeCaseApproval struct {
		model.Model
		SubmittedBy           *model.UserInfo `json:"submittedBy"`
		SubmittedByActionType string          `json:"submittedByActionType"`
		IsCompleted           bool            `json:"isCompleted"`
		IsRejected            bool            `json:"isRejected"`
		IsCurrent             bool            `json:"isCurrent"`
		SubmittedAt           *time.Time      `json:"submittedAt"`
		Seq                   int             `json:"seq"`
	}
	WorkerSimpleProfileCriteria struct {
		ProjRefId    *intstring.IntString `json:"projRefId"`
		CompanyRefId *intstring.IntString `json:"companyRefId"`
		DesRefId     *intstring.IntString `json:"desRefId"`

		ViewType      *intstring.IntString `json:"viewType"`
		IsBlacklisted *bool                `json:"isBlacklisted"`
		PointFrom     *intstring.IntString `json:"pointFrom"`
		PointTo       *intstring.IntString `json:"pointTo"`
		IdentityNo    *string              `json:"identityNo"`
		SyncAfter     *time.Time           `json:"-"`

		IdentityType *string               `json:"identityType"`
		CertTypeIds  []intstring.IntString `json:"certTypeIds"`
		pagination.Pagination
	}
	WorkerSimpleProfile struct {
		Id               intstring.IntString  `gorm:"primaryKey" json:"id,omitempty"`
		Uuid             string               `json:"uuid" gorm:"<-:create"`
		Status           string               `json:"status"`
		NameEn           *string              `json:"nameEn"`
		NameZh           *string              `json:"nameZh"`
		IsBlacklisted    *bool                `json:"isBlacklisted"`
		DesRefId         *intstring.IntString `json:"desRefId"`
		ProjRefCode      string               `json:"projRefCode"`
		ProjRefId        intstring.IntString  `json:"projRefId"`
		ProjRefName      *string              `json:"projRefName"`
		CompanyRefId     *intstring.IntString `json:"companyRefId"`
		PhoneNo          *string              `json:"phoneNo"`
		UpdatedBy        *string              `json:"updatedBy" gorm:"column:updated_by"`
		UpdatedAt        *string              `json:"updatedAt" gorm:"column:updated_at"`
		UpdatedByDisplay interface{}          `json:"updatedByDisplay" gorm:"-" `
		TotalPoint       int                  `json:"totalPoint" gorm:"-"`
	}
)
