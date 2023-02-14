package labour

import (
	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"time"
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
)
