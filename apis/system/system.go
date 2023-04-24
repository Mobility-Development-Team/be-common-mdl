package system

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"

	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
)

const (
	apiSystemMdlUrlBase   = "apis.internal.system.module.url.base"
	getOneContract        = "%s/contracts/%s"
	getAllLocations       = "%s/locations/all"
	getSupportInfo        = "%s/config/supportinfo"
	getAllContracts       = "%s/contracts/all"
	getContractParties    = "%s/parties/assoc/%s?groupBy=party"
	getManyParitesById    = "%s/parties/many"
	getContractUserByUids = "%s/parties/assoc/users"
)

func GetAllContracts(tk string, contractId ...intstring.IntString) (map[intstring.IntString]model.Contract, error) {
	var resp struct {
		response.Response
		Payload []*struct {
			model.Contract
			Id intstring.IntString `json:"id"`
		} `json:"payload"`
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"contractIds": append([]intstring.IntString{}, contractId...),
		},
	).Post(fmt.Sprintf(getAllContracts, apis.V().GetString(apiSystemMdlUrlBase)))
	if err != nil {
		logger.Error("[GetAllContracts] err: ", err)
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("system module returned status code: %d", result.StatusCode())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	output := map[intstring.IntString]model.Contract{}
	for _, r := range resp.Payload {
		if r == nil || r.Id == 0 {
			continue
		}
		output[r.Id] = r.Contract
	}
	return output, nil
}

func GetOneContract(tk string, contractId intstring.IntString) (*model.Contract, error) {
	type (
		respType struct {
			response.Response
			Payload *model.Contract `json:"payload"`
		}
	)
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneContract, apis.V().GetString(apiSystemMdlUrlBase), contractId))
	if err != nil {
		logger.Error("[GetOneContract]", "err:", err)
		return nil, err
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func GetManyPartiesById(tk string, ids ...intstring.IntString) ([]*model.PartyInfo, error) {
	var resp struct {
		response.Response
		Payload []*model.PartyInfo `json:"payload"`
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"ids": ids,
		},
	).Post(fmt.Sprintf(getManyParitesById, apis.V().GetString(apiSystemMdlUrlBase)))
	if err != nil {
		logger.Error("[GetManyPartiesById] err: ", err)
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("system module returned status code: %d", result.StatusCode())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

type partyInfoWithType struct {
	Info         model.PartyInfo `json:"info"`
	PartyTypeRef string          `json:"partyType"`
}

func GetContractUserByUids(tk string, contractId intstring.IntString, uids ...intstring.IntString) (map[intstring.IntString]*partyInfoWithType, error) {
	var resp struct {
		response.Response
		Payload map[intstring.IntString]*partyInfoWithType `json:"payload"`
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"contractId": contractId,
			"uids":       uids,
		},
	).Post(fmt.Sprintf(getContractUserByUids, apis.V().GetString(apiSystemMdlUrlBase)))
	if err != nil {
		logger.Error("[GetContractUserByUids] err: ", err)
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("system module returned status code: %d", result.StatusCode())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func GetLocations(tk string, body map[string]interface{}) (map[intstring.IntString]*model.Location, error) {
	urlPath := getAllLocations + "?showAsMap=true"
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(urlPath, apis.V().GetString(apiSystemMdlUrlBase)))
	if err != nil {
		logger.Error("[GetLocations]", "err:", err)
		return map[intstring.IntString]*model.Location{}, err
	}
	type respType struct {
		response.Response
		Payload map[intstring.IntString]*model.Location `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return map[intstring.IntString]*model.Location{}, err
	}
	return resp.Payload, nil
}

// A version of GetSupportInfo that logs the error and retruns the initialized map value on error
func ShouldGetSupportInfo() map[string]string {
	supportInfo, err := GetSupportInfo()
	if err != nil {
		logger.Error("[ShouldGetSupportInfo] Unable to get support info, support info would be missing for functions depending on it")
		return map[string]string{}
	}
	if supportInfo == nil {
		logger.Error("[ShouldGetSupportInfo] Call was successful but a nil support info was received")
		return map[string]string{}
	}
	return supportInfo
}

// A version of GetOneContract that logs the error and retruns an empty value on error
func ShouldGetOneContract(tk string, contractId *intstring.IntString) model.Contract {
	if contractId == nil {
		logger.Warn("[ShouldGetOneContract] Contract id is nil")
		return model.Contract{}
	}
	contract, err := GetOneContract(tk, *contractId)
	if err != nil {
		logger.Error("[ShouldGetOneContract] Unable to get contract information, ignoring")
		return model.Contract{}
	}
	if contract == nil {
		logger.Error("[ShouldGetOneContract] Call was successful but a nil contract was received")
		return model.Contract{}
	}
	return *contract
}

func GetSupportInfo() (map[string]string, error) {
	client := resty.New()
	result, err := client.R().Get(fmt.Sprintf(getSupportInfo, apis.V().GetString(apiSystemMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	if result.IsError() {
		return nil, errors.New("API returns status: " + result.Status())
	}
	type respType struct {
		response.Response
		Payload map[string]string `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func GetContractParties(tk string, contractId intstring.IntString) (map[string]ContractParty, error) {
	resp := struct {
		Payload map[string]ContractParty `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getContractParties, apis.V().GetString(apiSystemMdlUrlBase), contractId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("system module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func ShouldPopulatePartyInfo(tk string, partyInfo []*model.PartyInfo) {
	if err := PopulatePartyInfo(tk, partyInfo); err != nil {
		logger.Error("[ShouldPopulatePartyInfo] Failed getting parties, ignoring ", err)
	}
}

// PopulatePartyInfo Gets all parties in partyInfo, replace them with the updated version
// It tries to look for the records by their id
func PopulatePartyInfo(tk string, partyInfo []*model.PartyInfo) error {
	var ids []intstring.IntString
	idMap := map[intstring.IntString][]*model.PartyInfo{}
	for _, info := range partyInfo {
		if info == nil {
			logger.Warn("[PopulatePartyInfo] Got a nil partyInfo, ignoring...")
			continue
		}
		if info.Id > 0 {
			if _, ok := idMap[info.Id]; !ok {
				ids = append(ids, info.Id)
			}
			idMap[info.Id] = append(idMap[info.Id], info)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	updatedInfos, err := GetManyPartiesById(tk, ids...)
	if err != nil {
		return err
	}
	for _, updated := range updatedInfos {
		if updated == nil {
			continue
		}
		for _, partyInfo := range idMap[updated.Id] {
			if partyInfo == nil {
				continue
			}
			*partyInfo = *updated
		}
	}
	return nil
}
