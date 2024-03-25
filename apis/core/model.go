package core

import "github.com/Mobility-Development-Team/be-common-mdl/types/intstring"

type (
	UserAssocRelatedInfo struct {
		UserId     intstring.IntString `json:"userId"`
		ContractId intstring.IntString `json:"contractId"`
		PartyId    intstring.IntString `json:"partyId"`
		RoleNames  []string            `json:"roleNames"`
	}
)
