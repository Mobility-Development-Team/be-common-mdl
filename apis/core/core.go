package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/apis/auth"
	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/Mobility-Development-Team/be-common-mdl/util/apiutil"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
)

// Gin context storage keys
const (
	tokenInfoUser = "userInfo"
)

const (
	apiCoreMdlUrlBase       = "apis.internal.core.module.url.base"
	getAllUserInfo          = "%s/users/all"
	getOneContract          = "%s/contracts/%s"
	getAllContracts         = "%s/contracts/all"
	getManyContractMapUsers = "%s/contracts/users/map"
	getContractIdsUserMap   = "%s/contracts/users/map/contractIds"
	getSupportInfo          = "%s/support/info"
	getAllLocations         = "%s/locations/all"
	getContractUserByUids   = "%s/parties/assoc/users"
	getUserByRole           = "%s/roles/users"
	getContractParties      = "%s/parties/assoc/%s?groupBy=party"
	getManyParitesById      = "%s/parties/all"
	getUserByRoleAndParty   = "%s/users/role/party"
	getAdminUser            = "%s/users/admin/user"
	findAllRolesUnderUser   = "%s/users/roles/assoc/all"
	getUserHashtags         = "%s/users/hashtags/all"
	getAllRoles             = "%s/roles/all"
	inactiveUser            = "%s/users/inactivate"
	getUsersByGroupCriteria = "%s/roles/users/grouped"
	getUserList             = "%s/users/list"
	getSimpleUserList       = "%s/users/simple/list"
)

var muGetCurrentUserInfoFromContext sync.Mutex

// by id or useKeyRef to get user info
func GetUserById(tk string, id *intstring.IntString, userKeyRef *string, withSign *bool) (*model.UserInfo, error) {
	var ids []intstring.IntString
	var userKeyRefs []string
	var users []model.UserInfo
	var err error
	if id == nil && userKeyRef == nil {
		return nil, nil // Nothing specified, returns nil user
	}
	if id != nil {
		ids = []intstring.IntString{*id}
	}
	if userKeyRef != nil {
		userKeyRefs = []string{*userKeyRef}
	}

	users, err = GetUsersByIds(tk, ids, userKeyRefs, withSign)
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

func GetUsersByIds(tk string, ids []intstring.IntString, userKeyRefs []string, withSign *bool) ([]model.UserInfo, error) {
	if len(ids) == 0 && len(userKeyRefs) == 0 {
		return []model.UserInfo{}, nil
	}
	client := resty.New()
	body := map[string]interface{}{
		"ids":           ids,
		"userKeyRefs":   userKeyRefs,
		"withSignature": withSign,
	}
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(
		fmt.Sprintf(getUserList, apis.V().GetString(apiCoreMdlUrlBase)),
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

func GetSimpleUsersByIds(tk string, ids []intstring.IntString, userKeyRefs []string) ([]model.SimpleUserInfo, error) {
	if len(ids) == 0 && len(userKeyRefs) == 0 {
		return []model.SimpleUserInfo{}, nil
	}
	client := resty.New()
	body := map[string]interface{}{
		"ids":         ids,
		"userKeyRefs": userKeyRefs,
	}
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(
		fmt.Sprintf(getSimpleUserList, apis.V().GetString(apiCoreMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}

	var resp struct {
		response.Response
		Payload struct {
			Users      []model.SimpleUserInfo `json:"users"`
			TotalCount int                    `json:"totalCount"`
		} `json:"payload"`
	}
	if !result.IsSuccess() {
		return nil, errors.New("api returns status: " + result.Status())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
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

func GetManyContractMapUsers(tk string, req map[string]interface{}) ([]model.ContractToUserDetailMap, error) {
	var resp struct {
		response.Response
		Payload []model.ContractToUserDetailMap `json:"payload"`
	}
	client := resty.New()

	result, err := client.R().SetAuthToken(tk).SetBody(
		req,
	).Post(fmt.Sprintf(getManyContractMapUsers, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		logger.Error("[GetManyContractMapUsers] err: ", err)
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

func GetContractIdsUserMap(tk string, req map[string]interface{}) (*model.ContractIdsUserMap, error) {
	var resp struct {
		response.Response
		Payload *model.ContractIdsUserMap `json:"payload"`
	}
	client := resty.New()

	result, err := client.R().SetAuthToken(tk).SetBody(
		req,
	).Post(fmt.Sprintf(getContractIdsUserMap, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		logger.Error("[GetContractIdsUserMap] err: ", err)
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
	updatedInfos, err := GetUsersByIds(tk, ids, keyRefs, nil)
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

func GetUsersIdByRole(tk string, body map[string]interface{}) ([]intstring.IntString, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(getUserByRole, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		return []intstring.IntString{}, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("user module returned status code: %d", result.StatusCode())
	}
	type respType struct {
		response.Response
		Payload []intstring.IntString `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return []intstring.IntString{}, err
	}
	return resp.Payload, nil
}

func GetCurrentUserInfoFromContext(c *gin.Context) (*model.UserInfo, error) {
	muGetCurrentUserInfoFromContext.Lock()
	defer muGetCurrentUserInfoFromContext.Unlock()
	refKey := auth.GetUserRefKeyFromContext(c)
	tk, _ := apiutil.ParseBearerAuth(c)
	v, ok := c.Get(tokenInfoUser)
	switch {
	case ok && v != nil:
		logger.Debugf("[GetCurrentUserInfoFromContext] Reusing user info of creater %s...", refKey)
	case ok && v == nil:
		logger.Error("[GetCurrentUserInfoFromContext] Got a nil user info from cache, trying to get another one...")
		fallthrough
	default:
		logger.Debugf("[GetCurrentUserInfoFromContext] Getting user info of creater %s...", refKey)
		var err error
		var withSign *bool
		v, err = GetUserById(tk, nil, &refKey, withSign)
		if err != nil {
			return nil, err
		}
		c.Set(tokenInfoUser, v.(*model.UserInfo))
	}
	return v.(*model.UserInfo), nil
}

// move from user
// GenerateModelUserDisplay generates empty userInfo for the models and returns them in a single list.
// The UserInfo are filled with the users' refKey only.
// To load the UserInfo with the correct values, call PopulateUserInfo() with the returned list
func GenerateModelUserDisplay(models ...*model.Model) []*model.UserInfo {
	userInfos := make([]*model.UserInfo, 0, len(models)*2)
	for _, m := range models {
		u := &model.UserInfo{
			UserRefKey: m.CreatedBy,
		}
		m.CreatedByDisplay = u
		userInfos = append(userInfos, u)
		if m.UpdatedBy != nil {
			u := &model.UserInfo{
				UserRefKey: *m.UpdatedBy,
			}
			m.UpdatedByDisplay = u
			userInfos = append(userInfos, u)
		}
	}
	return userInfos
}

// move from system
func ShouldPopulatePartyInfo(tk string, partyInfo []*model.CorePartyInfoDisplay) {
	if err := PopulatePartyInfo(tk, partyInfo); err != nil {
		logger.Error("[ShouldPopulatePartyInfo] Failed getting parties, ignoring ", err)
	}
}

func PopulatePartyInfo(tk string, partyInfo []*model.CorePartyInfoDisplay) error {
	var ids []intstring.IntString
	var keyRefs []intstring.IntString
	idMap := map[intstring.IntString][]*model.CorePartyInfoDisplay{}
	keyRefMap := map[intstring.IntString][]*model.CorePartyInfoDisplay{}
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
		if info.Id != 0 {
			if _, ok := keyRefMap[info.Id]; !ok {
				keyRefs = append(keyRefs, info.Id)
			}
			keyRefMap[info.Id] = append(keyRefMap[info.Id], info)
		}
	}
	if len(ids) == 0 && len(keyRefs) == 0 {
		return nil
	}
	updatedInfos, err := GetManyPartiesById(tk, ids...)
	if err != nil {
		return err
	}
	for _, updated := range updatedInfos {
		for _, pInfo := range idMap[updated.Id] {
			if pInfo == nil {
				continue
			}
			*pInfo = *updated
		}
		for _, partyInfo := range keyRefMap[updated.Id] {
			if partyInfo == nil {
				continue
			}
			*partyInfo = *updated
		}
	}
	return nil
}

// move from system
func GetManyPartiesById(tk string, ids ...intstring.IntString) ([]*model.CorePartyInfoDisplay, error) {
	if len(ids) == 0 {
		return []*model.CorePartyInfoDisplay{}, nil
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"ids": ids,
		}).Post(
		fmt.Sprintf(getManyParitesById, apis.V().GetString(apiCoreMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}

	var resp struct {
		response.Response
		Payload struct {
			ContractId *intstring.IntString          `json:"contractId"`
			Parties    []*model.CorePartyInfoDisplay `json:"parties"`
			TotalCount int                           `json:"totalCount"`
		} `json:"payload"`
	}
	if !result.IsSuccess() {
		return nil, errors.New("api returns status: " + result.Status())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	for i := range resp.Payload.Parties {
		resp.Payload.Parties[i].ShouldAddSystemFieldsFromDisplay()

	}
	return resp.Payload.Parties, nil
}

func GetContractParties(tk string, contractId intstring.IntString, showModuleInfo bool) (model.CoreContractPartyInfoDisplay, error) {
	resp := struct {
		Payload model.CoreContractPartyInfoDisplay `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(map[string]interface{}{
		"contractId":     contractId,
		"showModuleInfo": showModuleInfo,
	}).Post(fmt.Sprintf(getManyParitesById, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		return model.CoreContractPartyInfoDisplay{}, err
	}
	if !result.IsSuccess() {
		return model.CoreContractPartyInfoDisplay{}, fmt.Errorf("core module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return model.CoreContractPartyInfoDisplay{}, err
	}
	resp.Payload.ShouldAddSystemFieldsFromDisplay()
	for i := range resp.Payload.Parties {
		resp.Payload.Parties[i].ShouldAddSystemFieldsFromDisplay()
	}

	return resp.Payload, err
}

func GetUsersByRoleAndParty(tk string, roleName string, contractId, partyId intstring.IntString) ([]model.UserInfo, error) {
	var resp struct {
		response.Response
		Payload []model.UserInfo `json:"payload"`
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"roleName":   roleName,
			"contractId": contractId,
			"partyId":    partyId,
		},
	).Post(fmt.Sprintf(getUserByRoleAndParty, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		logger.Error("[GetUserByRoleAndParty] err: ", err)
		return []model.UserInfo{}, err
	}
	if !result.IsSuccess() {
		return []model.UserInfo{}, fmt.Errorf("Core module returned status code: %d", result.StatusCode())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return []model.UserInfo{}, err
	}
	return resp.Payload, nil
}

func GetAdminUsers(tk string, contractId, partyId intstring.IntString) ([]model.UserInfo, error) {
	if contractId == 0 && partyId == 0 {
		return []model.UserInfo{}, nil
	}
	client := resty.New()
	body := map[string]interface{}{
		"contractId": contractId,
		"partyId":    partyId,
	}
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(
		fmt.Sprintf(getAdminUser, apis.V().GetString(apiCoreMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}

	var resp struct {
		response.Response
		Payload []model.UserInfo `json:"payload"`
	}
	if !result.IsSuccess() {
		return nil, errors.New("api returns status: " + result.Status())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	// for i := range resp.Payload {
	// 	resp.Payload[i].ShouldAddSystemFieldsFromDisplay()

	// }

	return resp.Payload, nil
}

func FindAllRolesUnderUser(tk string, userId, partyId, contractId intstring.IntString, userKey string) (result []UserAssocRelatedInfo, err error) {
	var (
		resp struct {
			response.Response
			Payload []UserAssocRelatedInfo `json:"payload"`
		}
	)
	client := resty.New()
	r, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"userKey":    userKey,
			"userId":     userId,
			"partyId":    partyId,
			"contractId": contractId,
		},
	).Post(fmt.Sprintf(findAllRolesUnderUser, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		logger.Error("[FindAllRolesUnderUser] err: ", err)
		return
	}
	if !r.IsSuccess() {
		err = fmt.Errorf("core module returned status code: %d", r.StatusCode())
		return
	}
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return
	}
	result = resp.Payload
	return
}

func GetAllUserHashTag(tk string, contractId intstring.IntString) ([]model.HashtagInfo, error) {
	var resp struct {
		response.Response
		Payload []model.HashtagInfo `json:"payload"`
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"contractId": contractId,
		},
	).Post(fmt.Sprintf(getUserHashtags, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		logger.Error("[GetAllUserHashTag] err: ", err)
		return []model.HashtagInfo{}, err
	}
	if !result.IsSuccess() {
		return []model.HashtagInfo{}, fmt.Errorf("Core module returned status code: %d", result.StatusCode())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return []model.HashtagInfo{}, err
	}
	return resp.Payload, nil
}

func GetAllRole(tk string) ([]model.CoreRole, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(nil).Post(
		fmt.Sprintf(getAllRoles, apis.V().GetString(apiCoreMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}

	var resp struct {
		response.Response
		Payload struct {
			Roles      []model.CoreRole `json:"roles"`
			TotalCount int              `json:"totalCount"`
		} `json:"payload"`
	}
	if !result.IsSuccess() {
		return nil, errors.New("api returns status: " + result.Status())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	for i := range resp.Payload.Roles {
		resp.Payload.Roles[i].ShouldAddSystemFieldsFromDisplay()

	}
	return resp.Payload.Roles, nil
}

func GetRoleHastag(tk string) ([]model.HashtagInfo, error) {

	rhi := []model.HashtagInfo{}

	roleInfos, err := GetAllRole(tk)
	if err != nil {
		return nil, err
	}

	if len(roleInfos) > 0 {
		for _, r := range roleInfos {
			rhi = append(rhi, model.HashtagInfo{
				Id:         r.Id,
				Name:       r.RoleName,
				Type:       "ROLE",
				UserRefKey: r.Uuid,
			})
		}
	}

	return rhi, nil
}

func GetLocationHashtagByContractId(tk string, contractId intstring.IntString) ([]model.HashtagInfo, error) {
	lhi := []model.HashtagInfo{}
	locInfos, err := GetLocations(tk, map[string]interface{}{
		"contractId": contractId,
	})
	if err != nil {
		return nil, err
	}

	if len(locInfos) > 0 {
		for _, r := range locInfos {
			for _, c := range r {
				lhi = append(lhi, model.HashtagInfo{
					Id:         c.Id,
					Name:       fmt.Sprintf("%s %s", c.Name, *c.NameZh),
					Type:       "LOCATION",
					UserRefKey: c.Uuid,
				})
			}
		}
	}

	return lhi, nil
}

func InactivateUserAcc(tk string, userRefKey string) error {
	body := map[string]interface{}{
		"userRefKey": userRefKey,
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(inactiveUser, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil || result.StatusCode() != 200 {
		return errors.New(result.String())
	}
	return nil
}

func GetUsersByGroupCriteria(tk string, body map[string]interface{}) (map[string][]model.UserInfo, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(getUsersByGroupCriteria, apis.V().GetString(apiCoreMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("user module returned status code: %d", result.StatusCode())
	}
	type respType struct {
		response.Response
		Payload map[string][]model.UserInfo `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
