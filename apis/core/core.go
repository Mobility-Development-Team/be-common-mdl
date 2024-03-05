package core

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
	apiCoreMdlUrlBase     = "apis.internal.core.module.url.base"
	getAllUserInfo        = "%s/users/all"
	getOneContract        = "%s/contracts/%s"
	getAllContracts       = "%s/contracts/all"
	getSupportInfo        = "%s/support/info"
	getAllLocations       = "%s/locations/all"
	getContractUserByUids = "%s/parties/assoc/users"
)

// by id or useKeyRef to get user info
func GetUserById(tk string, id *intstring.IntString, userKeyRef *string) (*model.UserInfo, error) {
	var ids []intstring.IntString
	var userKeyRefs []string
	if id == nil && userKeyRef == nil {
		return nil, nil // Nothing specified, returns nil user
	}
	if id != nil {
		ids = []intstring.IntString{*id}
	}
	if userKeyRef != nil {
		userKeyRefs = []string{*userKeyRef}
	}
	users, err := GetUsersByIds(tk, ids, userKeyRefs)
	if err != nil {
		return nil, err
	}
	var result *model.UserInfo
	if len(users) > 0 {
		result = &users[0]
	}
	return result, nil
}

func GetAllUserInfo(tk string, body map[string]interface{}) ([]model.GetUserResponse, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(getAllUserInfo, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		return []model.GetUserResponse{}, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("user module returned status code: %d", result.StatusCode())
	}
	type respType struct {
		response.Response
		Payload []model.GetUserResponse `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return []model.GetUserResponse{}, err
	}
	return resp.Payload, nil
}

func GetUsersByIds(tk string, ids []intstring.IntString, userKeyRefs []string) ([]model.UserInfo, error) {
	if len(ids) == 0 && len(userKeyRefs) == 0 {
		return []model.UserInfo{}, nil
	}
	client := resty.New()
	body := map[string]interface{}{
		"ids":         ids,
		"userKeyRefs": userKeyRefs,
	}
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(
		fmt.Sprintf(getAllUserInfo, apis.V().GetString(apiCoreMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}

	var resp struct {
		response.Response
		Payload struct {
			Users      []model.UserInfo `json:"users"`
			TotalCount int              `json:"totalCount"`
		} `json:"payload"`
	}
	if !result.IsSuccess() {
		return nil, errors.New("api returns status: " + result.Status())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	for i := range resp.Payload.Users {
		resp.Payload.Users[i].ShouldAddSystemFieldsFromDisplay()

	}

	return resp.Payload.Users, nil
}

func GetOneContract(tk string, contractId intstring.IntString) (*model.GetCoreContractResponse, error) {
	type (
		respType struct {
			response.Response
			Payload *model.GetCoreContractResponse `json:"payload"`
		}
	)
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneContract, apis.V().GetString(apiCoreMdlUrlBase), contractId))
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

func GetAllContracts(tk string, projectId *string, contractIds ...intstring.IntString) (map[intstring.IntString][]model.GetCoreContractResponse, error) {
	var resp struct {
		response.Response
		Payload struct {
			Contracts  []model.GetCoreContractResponse `json:"contracts"`
			TotalCount int                             `json:"totalCount"`
			// Id         intstring.IntString             `json:"id"`
		} `json:"payload"`
	}
	client := resty.New()
	req := map[string]interface{}{
		"contractIds": contractIds,
	}
	if projectId != nil {
		req["projectIdRef"] = *projectId
	}
	result, err := client.R().SetAuthToken(tk).SetBody(
		req,
	).Post(fmt.Sprintf(getAllContracts, apis.V().GetString(apiCoreMdlUrlBase)))
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
	output := map[intstring.IntString][]model.GetCoreContractResponse{}

	for _, c := range resp.Payload.Contracts {
		c.ShouldAddSystemFieldsFromDisplay()
		output[c.Id] = append(output[c.Id], model.GetCoreContractResponse{
			CoreContract: c.CoreContract,
			Parties:      c.Parties,
			UserCount:    c.UserCount,
			PartyCount:   c.PartyCount,
			RoleNames:    c.RoleNames,
			PartyType:    c.PartyType,
		})
	}
	return output, nil
}

func GetSupportInfo() (map[string]string, error) {
	client := resty.New()
	result, err := client.R().Get(fmt.Sprintf(getSupportInfo, apis.V().GetString(apiCoreMdlUrlBase)))
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
func ShouldGetOneContract(tk string, contractId *intstring.IntString) model.GetCoreContractResponse {
	if contractId == nil {
		logger.Warn("[ShouldGetOneContract] Contract id is nil")
		return model.GetCoreContractResponse{}
	}
	contract, err := GetOneContract(tk, *contractId)
	if err != nil {
		logger.Error("[ShouldGetOneContract] Unable to get contract information, ignoring")
		return model.GetCoreContractResponse{}
	}
	if contract == nil {
		logger.Error("[ShouldGetOneContract] Call was successful but a nil contract was received")
		return model.GetCoreContractResponse{}
	}
	return *contract
}

func PopulateUserInfo(tk string, userInfo []*model.UserInfo) error {
	var ids []intstring.IntString
	var keyRefs []string
	idMap := map[intstring.IntString][]*model.UserInfo{}
	keyRefMap := map[string][]*model.UserInfo{}
	for _, info := range userInfo {
		if info == nil {
			logger.Warn("[PopulateUserInfo] Got a nil userInfo, ignoring...")
			continue
		}
		if info.Id > 0 {
			if _, ok := idMap[info.Id]; !ok {
				ids = append(ids, info.Id)
			}
			idMap[info.Id] = append(idMap[info.Id], info)
		}
		if info.UserRefKey != "" {
			if _, ok := keyRefMap[info.UserRefKey]; !ok {
				keyRefs = append(keyRefs, info.UserRefKey)
			}
			keyRefMap[info.UserRefKey] = append(keyRefMap[info.UserRefKey], info)
		}
	}
	if len(ids) == 0 && len(keyRefs) == 0 {
		return nil
	}
	updatedInfos, err := GetUsersByIds(tk, ids, keyRefs)
	if err != nil {
		return err
	}
	for _, updated := range updatedInfos {
		for _, userInfo := range idMap[updated.Id] {
			if userInfo == nil {
				continue
			}
			*userInfo = updated
		}
		for _, userInfo := range keyRefMap[updated.UserRefKey] {
			if userInfo == nil {
				continue
			}
			*userInfo = updated
		}
	}
	return nil
}

func GetLocations(tk string, body map[string]interface{}) (map[intstring.IntString][]*model.Location, error) {
	var resp struct {
		response.Response
		Payload struct {
			Locations  []*model.Location `json:"locations"`
			TotalCount int               `json:"totalCount"`
		} `json:"payload"`
	}
	urlPath := getAllLocations
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(urlPath, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		logger.Error("[GetLocations]", "err:", err)
		return map[intstring.IntString][]*model.Location{}, err
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return map[intstring.IntString][]*model.Location{}, err
	}

	output := map[intstring.IntString][]*model.Location{}

	for _, c := range resp.Payload.Locations {
		c.ShouldAddSystemFieldsFromDisplay()
		if c.Id != 0 {
			output[c.Id] = append(output[c.Id], &model.Location{
				Id:            c.Id,
				Uuid:          c.Uuid,
				Name:          c.Name,
				NameZh:        c.NameZh,
				Status:        c.Status,
				LocationType:  c.LocationType,
				Latitude:      c.Latitude,
				Longitude:     c.Longitude,
				LocationId:    c.LocationId,
				ContractRefId: c.ContractRefId,
			})
		} else {
			output[*c.ContractRefId] = append(output[*c.ContractRefId], &model.Location{
				Id:            c.Id,
				Uuid:          c.Uuid,
				Name:          c.Name,
				NameZh:        c.NameZh,
				Status:        c.Status,
				LocationType:  c.LocationType,
				Latitude:      c.Latitude,
				Longitude:     c.Longitude,
				LocationId:    c.LocationId,
				ContractRefId: c.ContractRefId,
			})
		}

	}

	return output, nil
}

func GetContractUserByUids(tk string, contractId intstring.IntString, uids ...intstring.IntString) (map[intstring.IntString]*intstring.IntString, error) {
	var resp struct {
		response.Response
		Payload map[intstring.IntString]*intstring.IntString `json:"payload"`
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"contractId": contractId,
			"uids":       uids,
		},
	).Post(fmt.Sprintf(getContractUserByUids, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		logger.Error("[GetContractUserByUids] err: ", err)
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("Core module returned status code: %d", result.StatusCode())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
