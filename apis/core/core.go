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
	apiCoreMdlUrlBase = "apis.internal.core.module.url.base"
	getAllUserInfo    = "%s/users/all"
	getOneContract    = "%s/contracts/%s"
	getAllContracts   = "%s/contracts/all"
	getSupportInfo    = "%s/support/info"
)

func GetAllUserInfoAsMap(tk string, body map[string]interface{}) (map[string]model.UserInfo, error) {
	urlPath := getAllUserInfo + "?showAsMap=true"
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(urlPath, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		return map[string]model.UserInfo{}, err
	}
	type respType struct {
		response.Response
		Payload map[string]model.UserInfo `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return map[string]model.UserInfo{}, err
	}
	return resp.Payload, nil
}

func GetAllUserInfo(tk string, body map[string]interface{}) ([]model.UserInfo, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(getAllUserInfo, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		return []model.UserInfo{}, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("user module returned status code: %d", result.StatusCode())
	}
	type respType struct {
		response.Response
		Payload []model.UserInfo `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return []model.UserInfo{}, err
	}
	return resp.Payload, nil
}

func GetUsersByIds(tk string, ids []intstring.IntString, userKeyRefs []string) ([]model.GetUserResponse, error) {
	if len(ids) == 0 && len(userKeyRefs) == 0 {
		return []model.GetUserResponse{}, nil
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
			Users      []model.GetUserResponse `json:"users"`
			TotalCount int                     `json:"totalCount"`
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

func GetOneContract(tk string, contractId intstring.IntString) (*model.Contract, error) {
	type (
		respType struct {
			response.Response
			Payload *model.Contract `json:"payload"`
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

func GetAllContracts(tk string, projectId *string, contractId ...intstring.IntString) (map[intstring.IntString]model.Contract, error) {
	var resp struct {
		response.Response
		Payload []*struct {
			model.Contract
			Id intstring.IntString `json:"id"`
		} `json:"payload"`
	}
	client := resty.New()
	req := map[string]interface{}{
		"contractIds": append([]intstring.IntString{}, contractId...),
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
	output := map[intstring.IntString]model.Contract{}
	for _, r := range resp.Payload {
		if r == nil || r.Id == 0 {
			continue
		}
		output[r.Id] = r.Contract
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
