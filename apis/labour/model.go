package labour

import (
	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
)

type (
	UnsafeCase struct {
		model.Model
		CaseNo        string              `json:"caseNo"`
		CaseStatus    string              `json:"caseStatus"`
		CaseType      string              `json:"caseType"`
		ContractRefId intstring.IntString `json:"contractId"`
		Worker        interface{}         `json:"worker"`
		Item          interface{}         `json:"item"`
	}
	UnsafeCaseCriteria struct {
		CaseStatuses []string `json:"CaseStatuses"`
	}
)
