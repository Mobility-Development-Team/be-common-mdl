package system

import "github.com/Mobility-Development-Team/be-common-mdl/types/intstring"

type ContractParty struct {
	Info struct {
		Id            intstring.IntString `json:"id"`
		PartyName     string              `json:"partyName"`
		PartyNameZh   string              `json:"partyNameZh"`
		Address       string              `json:"address"`
		Email         string              `json:"email"`
		Br            string              `json:"br"`
		TradeCategory string              `json:"tradeCategory"`
		PartyIconUrl  string              `json:"partyIconUrl"`
		PartyPrefix   string              `json:"partyPrefix"`
		SubconRefId   string              `json:"subconRefId"`
	} `json:"info"`
	PartyTypeRef string                `json:"partyType"`
	UserIds      []intstring.IntString `json:"userIds"`
	GroupIds     []intstring.IntString `json:"groupIds"`
	RoleIds      []intstring.IntString `json:"roleIds"`
}
