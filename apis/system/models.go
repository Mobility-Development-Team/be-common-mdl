package system

import (
	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
)

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
	PartyTypeRef          string                `json:"partyType"`
	UserIds               []intstring.IntString `json:"userIds"`
	GroupIds              []intstring.IntString `json:"groupIds"`
	RoleIds               []intstring.IntString `json:"roleIds"`
	ContractPartyMappings []PartyMap            `json:"-" gorm:"foreignKey:PartyId"`
}

type PartyMap struct {
	model.Model
	Status       string              `json:"status"`
	ContractId   intstring.IntString `json:"contractId" gorm:"index:uk_contract_id_party_id,unique"`
	PartyId      intstring.IntString `json:"partyId" gorm:"index:uk_contract_id_party_id,unique"`
	PartyTypeRef string              `json:"partyTypeRef" gorm:"index:uk_contract_id_party_id,unique"`
	IsPrimary    bool                `json:"isPrimary"`
}
