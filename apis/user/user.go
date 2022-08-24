package user

import (
	"encoding/json"
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
	apiUserMdlUrlBase  = "apis.internal.user.module.url.base"
	getCurrentUserInfo = "%s/users/profile?isSimple=true"
	getAllUserInfo     = "%s/users/all"
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

func GetUsersByIds(tk string, ids []intstring.IntString, userKeyRefs []string) ([]model.UserInfo, error) {
	client := resty.New()
	body := map[string]interface{}{
		"ids":         ids,
		"userKeyRefs": userKeyRefs,
	}
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(
		fmt.Sprintf(getAllUserInfo, apis.V().GetString(apiUserMdlUrlBase)),
	)
	if err != nil {
		return nil, err
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

func PopulateModelUserDisplay(tk string, models ...*model.Model) error {
	userInfos := make([]*model.UserInfo, 0, len(models)*2)
	for _, m := range models {
		if m.CreatedBy != "" {
			u := &model.UserInfo{
				UserRefKey: m.CreatedBy,
			}
			m.CreatedByDisplay = u
			userInfos = append(userInfos, u)
		}
		if m.UpdatedBy != nil {
			u := &model.UserInfo{
				UserRefKey: *m.UpdatedBy,
			}
			m.UpdatedByDisplay = u
			userInfos = append(userInfos, u)
		}
	}
	return PopulateUserInfo(tk, userInfos)
}

// PopulateUserInfo Gets all users in userInfo, replace them with the updated version
// It tries to look for the records by either userKeyRef or id
func PopulateUserInfo(tk string, userInfo []*model.UserInfo) error {
	var ids []intstring.IntString
	var keyRefs []string
	idMap := map[intstring.IntString]*model.UserInfo{}
	keyRefMap := map[string]*model.UserInfo{}
	for _, info := range userInfo {
		if info == nil {
			logger.Warn("[PopulateUserInfo] Got a nil userInfo, ignoring...")
			continue
		}
		if info.Id > 0 {
			idMap[info.Id] = info
			ids = append(ids, info.Id)
		}
		if info.UserRefKey != "" {
			keyRefMap[info.UserRefKey] = info
			keyRefs = append(keyRefs, info.UserRefKey)
		}
	}
	updatedInfos, err := GetUsersByIds(tk, ids, keyRefs)
	if err != nil {
		return err
	}
	for _, updated := range updatedInfos {
		// Always try refKey first
		if hit, ok := keyRefMap[updated.UserRefKey]; ok {
			*hit = updated
			continue
		}
		if hit, ok := idMap[updated.Id]; ok {
			*hit = updated
			continue
		}
		logger.Warnf("[PopulateUserInfo] Skipped mapping, userInfo does not have the related record, id=%s refKey=%s", updated.Id, updated.UserRefKey)
	}
	return nil
}
