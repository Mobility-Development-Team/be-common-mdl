package user

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
	apiUserMdlUrlBase      = "apis.internal.user.module.url.base"
	getCurrentUserInfo     = "%s/users/profile?isSimple=true"
	getAllUserInfo         = "%s/users/all"
	getUserList            = "%s/users/list"
	getUsersByGroupDetails = "%s/users/groups/details"
	getUserSignatures      = "%s/users/signatures"
	getAllGroupInfo        = "%s/groups/all"
)

// ShouldGetCurrentUserInfoFromContext Similar to GetCurrentUserInfoFromContext, but returns an empty value if it fails with an error message
func ShouldGetCurrentUserInfoFromContext(c *gin.Context) model.UserInfo {
	userInfo, err := GetCurrentUserInfoFromContext(c)
	if err != nil {
		logger.Errorf("[ShouldGetCurrentUserInfoFromContext] Unable to get current user info, ignoring... userRef:%s err:%s", auth.GetUserRefKeyFromContext(c), err.Error())
		return model.UserInfo{}
	}
	if userInfo == nil {
		logger.Errorf("[ShouldGetCurrentUserInfoFromContext] Got an empty user info, ignoring... userRef:%s", auth.GetUserRefKeyFromContext(c))
		return model.UserInfo{}
	}
	return *userInfo
}

var muGetCurrentUserInfoFromContext sync.Mutex

// GetCurrentUserInfoFromContext Gets the user object by the refKey that is passed with the token.
// This call is lazy loaded, and would reuse the retrieved object if called more than once
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
		v, err = GetUserById(tk, nil, &refKey)
		if err != nil {
			return nil, err
		}
		c.Set(tokenInfoUser, v.(*model.UserInfo))
	}
	return v.(*model.UserInfo), nil
}

func GetAllUserInfoAsMap(tk string, body map[string]interface{}) (map[string]model.UserInfo, error) {
	urlPath := getAllUserInfo + "?showAsMap=true"
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(urlPath, apis.V().GetString(apiUserMdlUrlBase)))
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
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(getAllUserInfo, apis.V().GetString(apiUserMdlUrlBase)))
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

func GetAllGroupInfo(tk string, body map[string]interface{}) ([]model.GroupInfo, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(getAllGroupInfo, apis.V().GetString(apiUserMdlUrlBase)))
	if err != nil {
		return []model.GroupInfo{}, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("user module returned status code: %d", result.StatusCode())
	}
	type respType struct {
		response.Response
		Payload []model.GroupInfo `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return []model.GroupInfo{}, err
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
		fmt.Sprintf(getUserList, apis.V().GetString(apiUserMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, errors.New("api returns status: " + result.Status())
	}
	type respType struct {
		response.Response
		Payload []model.UserInfo `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// GetUserById Gets a user by id or userKeyRef (either is fine), returns the user information if found, nil if not found / error
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

// GetUsersByGroupDetails Gets a user by id or userKeyRef (either is fine), returns the user information if found, nil if not found / error
func GetUsersByGroupDetails(tk string, groupName *string, contractId, partyId *intstring.IntString) ([]model.UserInfo, error) {
	if nil == groupName || nil == contractId || nil == partyId {
		return []model.UserInfo{}, nil // Nothing specified, returns nil user
	}
	client := resty.New()
	body := map[string]interface{}{
		"groupName":  groupName,
		"contractId": contractId,
		"partyId":    partyId,
	}
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(
		fmt.Sprintf(getUsersByGroupDetails, apis.V().GetString(apiUserMdlUrlBase)),
	)
	if !result.IsSuccess() {
		return nil, errors.New("api returns status: " + result.Status())
	}
	type respType struct {
		response.Response
		Payload []model.UserInfo `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// GetUserSignatures gets user signatures by given user ids
//
// Retuns a map[userId]sig where sig is a base64 encoded string of a png image.
// sig is empty "" if the user has no signature saved.
func GetUserSignatures(tk string, ids []intstring.IntString) (map[intstring.IntString]string, error) {
	if len(ids) == 0 {
		return map[intstring.IntString]string{}, nil
	}
	client := resty.New()
	body := map[string]interface{}{
		"ids": ids,
	}
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(
		fmt.Sprintf(getUserSignatures, apis.V().GetString(apiUserMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, errors.New("api returns status: " + result.Status())
	}
	type respType struct {
		response.Response
		Payload map[intstring.IntString]string `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// GenerateModelUserDisplay generates empty userInfo for the models and returns them in a single list.
// The UserInfo are filled with the users' refKey only.
//
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

func PopulateModelUserDisplay(tk string, models ...*model.Model) error {
	return PopulateUserInfo(tk, GenerateModelUserDisplay(models...))
}

func ShouldPopulateModelUserDisplay(tk string, models ...*model.Model) {
	ShouldPopulateUserInfo(tk, GenerateModelUserDisplay(models...))
}

func ShouldPopulateUserInfo(tk string, userInfo []*model.UserInfo) {
	if err := PopulateUserInfo(tk, userInfo); err != nil {
		logger.Error("[ShouldPopulateUserInfo] Failed getting user, ignoring ", err)
	}
}

// PopulateUserInfo Gets all users in userInfo, replace them with the updated version
// It tries to look for the records by either userKeyRef or id
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
