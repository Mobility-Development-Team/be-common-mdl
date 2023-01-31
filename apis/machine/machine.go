package machine

import (
	"encoding/json"
	"fmt"
	"github.com/Mobility-Development-Team/be-common-mdl/apis/system"
	"github.com/Mobility-Development-Team/be-common-mdl/apis/user"
	"github.com/Mobility-Development-Team/be-common-mdl/util/arrutil"
	logger "github.com/sirupsen/logrus"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/go-resty/resty/v2"
)

const (
	apiMachineMdlUrlBase = "apis.internal.machine.module.url.base"
	getOnePlantPermit    = "%s/permits/plantpermits/%s"
	getOneNCAPermit      = "%s/permits/nca/%s"
	getOneHotworkPermit  = "%s/permits/hw/%s"
	getOneEXPermit       = "%s/permits/ex/%s"
	getOneELPermit       = "%s/permits/el/%s"
	getOneLA             = "%s/plant/equip/LA/detail"
	getAllPermits        = "%s/permits/internal/all"
)

func GetOneLA(tk string, criteria LA, isSimple bool) (*LA, error) {
	resp := struct {
		Payload *LA `json:"payload"`
	}{}
	uri := getOneLA
	if isSimple {
		uri += "?isSimple=true"
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(criteria).Post(fmt.Sprintf(uri, apis.V().GetString(apiMachineMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOnePlantPermit(tk string, permitMasterId intstring.IntString) (*PlantPermit, error) {
	resp := struct {
		Payload *PlantPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOnePlantPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOneNCAPermit(tk string, permitMasterId intstring.IntString) (*NCAPermit, error) {
	resp := struct {
		Payload *NCAPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneNCAPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetAllPermits(tk string, userRefKey string, criteria PermitCriteria, opt GetAllPermitOps, preloadNames ...string) ([]*MasterPermit, error) {
	resp := struct {
		Payload []*MasterPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"criteria": criteria,
			"opts":     opt,
			"preloads": preloadNames,
		},
	).Post(
		fmt.Sprintf(getAllPermits, apis.V().GetString(apiMachineMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOneHotworkPermit(tk string, permitMasterId intstring.IntString) (*HotworkPermit, error) {
	resp := struct {
		Payload *HotworkPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneHotworkPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOnePermitToDig(tk string, permitMasterId intstring.IntString) (*EXPermit, error) {
	resp := struct {
		Payload *EXPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneEXPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOneELPermit(tk string, permitMasterId intstring.IntString) (*ELPermit, error) {
	resp := struct {
		Payload *ELPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneELPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func AddPartyIdToParticipants(tk string, contractId intstring.IntString, participants ...*Participant) error {
	if len(participants) == 0 {
		return nil
	}
	uids := make([]intstring.IntString, 0, len(participants))
	var uRefs []string
	uRefToIdMap := map[string]intstring.IntString{}
	for _, p := range participants {
		switch {
		case p.ParticipantUserId > 0:
			uids = append(uids, p.ParticipantUserId)
		case p.ParticipantUserId == 0 && p.ParticipantUserRefKey != "":
			uRefs = append(uRefs, p.ParticipantUserRefKey)
		default:
			logger.Debug("[AddPartyIdToParticipants] Participant is not identifiable, cannot look for party id: ", p)
		}
	}
	if len(uRefs) > 0 {
		logger.Debug("[AddPartyIdToParticipants] Some participant(s) have only refKey available, proceeding to fetch userId: ", uRefs)
		users, err := user.GetUsersByIds(tk, nil, uRefs)
		if err != nil {
			return err
		}
		for _, u := range users {
			uids = append(uids, u.Id)
			uRefToIdMap[u.UserRefKey] = u.Id
		}
	}
	partyInfo, err := system.GetContractUserByUids(tk, contractId, arrutil.Unique(uids)...)
	if err != nil {
		return err
	}
	for _, p := range participants {
		var uid intstring.IntString
		switch {
		case p.ParticipantUserId > 0:
			uid = p.ParticipantUserId
		case p.ParticipantUserId == 0 && p.ParticipantUserRefKey != "":
			uid = uRefToIdMap[p.ParticipantUserRefKey]
		default:
			continue
		}
		if partyInfo, ok := partyInfo[uid]; ok {
			p.PartyRefId = &partyInfo.Info.Id
		} else {
			logger.Warn("[AddPartyIdToParticipants] Unable to find user party, continuing anyways: ", p.ParticipantUserId)
		}
	}
	return nil
}
