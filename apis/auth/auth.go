package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/Mobility-Development-Team/be-common-mdl/util/apiutil"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
)

// Gin context storage keys
const (
	tokenInfoUserRefKey = "userRefKey"
	keyTokenInfo        = "tokenInfo"
	AuthHeader          = "Authorization"
	AuthHeaderCust      = "Authorization-ext"
	AuthorizationBearer = "Bearer"
	AuthorizationBasic  = "Basic"
	// API constant
	apiAuthMdlUrlBase             = "apis.internal.auth.module.url.base"
	getTokenInfo                  = "%s/tokeninfo"
	validateEmatToken             = "%s/validate/smm/user"
	validateEmatTokenWithTk       = "%s/validate/smm/user?tk=true"
	findIdentitiesByUserKey       = "%s/users/identities"
	validateExternalByIdentity    = "%s/users/validate/external"
	createUserWithIdentities      = "%s/users"
	updateDeviceIdRegisterAttempt = "%s/users/device/regd"
	linkUserWithIdentity          = "%s/users/identity/link"
	unlinkUserWithIdentity        = "%s/users/identity/unlink"
	resetUserIdentityCredential   = "%s/users/identity/credential/reset"
	findAllLoginHistory           = "%s/users/login/histories"
	updateAuthUserlockStatus      = "%s/users/lock/status"
	getManyUserLockInfo           = "%s/users/lock/info"
	getUserInactive               = "%s/users/inactive/batch"
)

type (
	TokenInfoResp struct {
		UserId     string `json:"userId"`
		CExpiresIn int    `json:"cExpiresIn"`
		AExpiresIn int    `json:"aExpiresIn"`
		RExpiresIn int    `json:"rExpiresIn"`
	}
	ValidateEmatTokenResp struct {
		IsValid    bool   `json:"isValid"`
		Message    string `json:"message"`
		Email      string `json:"email"`
		UserRefKey string `json:"userRefKey"`
		Token      string `json:"tk"`
	}
)

func GetUserRefKeyFromContext(c *gin.Context) string {
	k, ok := c.Get(tokenInfoUserRefKey) // Assume
	if !ok {
		return "undefined"
	}
	return fmt.Sprintf("%v", k)
}

// NewTokenVerifierInterceptor Gets a gin middleware for handing token verifications. The returned intercepter should be registered
// as a middleware in gin for protected API calls such that ParseBearerAuth and GetUserRefKeyFromContext
// can be used. invalidHeaderMsg or invalidTokenMsg is returned to the user in case of error.
func NewTokenVerifierInterceptor(invalidHeaderMsg, invalidTokenMsg response.Message) gin.HandlerFunc {
	return func(c *gin.Context) {
		b := strings.ToLower(c.Query("smm")) == "true"
		if b {
			// Handle EMat token
			tk, isValid := parseCustomAuthHeader(c, fmt.Sprintf("%s ", AuthorizationBearer))
			if !isValid {
				logger.Warn("[NewTokenVerifierInterceptor] unable to get token from Authorization-ext")
				apiutil.GenerateResponse(c, nil, invalidTokenMsg)
				c.Abort()
				return
			}
			r, err := ValidateEMatToken(c, tk)
			if err != nil {
				logger.Warn("[NewTokenVerifierInterceptor] invalid EMat token")
				apiutil.GenerateResponse(c, nil, invalidTokenMsg)
				c.Abort()
				return
			}
			c.Set("somTk", r.Token)
			c.Set("userRefKey", r.UserRefKey)
			c.Next()
			return
		}
		tk, ok := apiutil.ParseBearerAuth(c)
		if !ok {
			logger.Warn("[ValidateInternalToken] unable to parse the given token from the header")
			apiutil.GenerateResponse(c, nil, invalidHeaderMsg)
			c.Abort()
			return
		}
		info, err := GetTokenInfo(c, tk)
		if err != nil {
			logger.Warn("[ValidateInternalToken] ", err)
			apiutil.GenerateResponse(c, nil, invalidTokenMsg)
			c.Abort()
			return
		}
		// Reserve Token User Key
		c.Set(keyTokenInfo, info)
		c.Set(tokenInfoUserRefKey, info.UserId)
		c.Next()
	}
}

func GetTokenInfo(c *gin.Context, tk string) (TokenInfoResp, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getTokenInfo, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil {
		return TokenInfoResp{}, err
	}
	if result.StatusCode() != 200 {
		return TokenInfoResp{}, errors.New("the given token is invalid")
	}
	var info TokenInfoResp
	if err := json.Unmarshal(result.Body(), &info); err != nil {
		return TokenInfoResp{}, err
	}
	return info, err
}

func CreateAuthUser(c *gin.Context, body map[string]interface{}) (*resty.Response, error) {
	client := resty.New()
	tk, _ := apiutil.ParseBearerAuth(c)
	v, _ := apiutil.ParseCustAuthExt(c, "")
	result, err := client.R().SetAuthToken(tk).SetHeader(apiutil.HeaderCustom, fmt.Sprintf("%s%s", apiutil.AuthHeaderPrefixBasic, v)).
		SetBody(body).Post(fmt.Sprintf(createUserWithIdentities, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil || result.StatusCode() != 200 {
		// c.Abort()
		// api.GenerateResponse(c, nil, message.MsgCodeCommon19002)
		return nil, errors.New(result.String())
	}
	return result, nil
}

func GetTokenInfoFromContext(c *gin.Context) (TokenInfoResp, error) {
	value, ok := c.Get(keyTokenInfo)
	if !ok {
		return TokenInfoResp{}, errors.New("the internal token is not yet validated or no internal token found")
	}
	info, ok := value.(TokenInfoResp)
	if !ok {
		return TokenInfoResp{}, errors.New("the validated token is not TokenInfoResp")
	}
	return info, nil
}

func ValidateEMatToken(c *gin.Context, tk string) (*ValidateEmatTokenResp, error) {
	client := resty.New()
	url := strings.TrimSpace(fmt.Sprintf(validateEmatTokenWithTk, apis.V().GetString(apiAuthMdlUrlBase)))
	result, err := client.R().SetHeaders(map[string]string{
		AuthHeaderCust: fmt.Sprintf("%s %s", AuthorizationBearer, tk),
		AuthHeader:     fmt.Sprintf("%s %s", AuthorizationBasic, "YzBhZjVlMDZiNTdlYmJlYTlhYTQ6ZGI4MDBjNzQ3ZjQ2MzgzOGM2NTQwMDQwYmM4ODM3MmNlZjVkNGVkMTlhNDU="),
	}).Get(url)
	if err != nil {
		return nil, err
	}
	if result.StatusCode() != 200 {
		return nil, errors.New("the given token is invalid")
	}
	var info ValidateEmatTokenResp
	if err := json.Unmarshal(result.Body(), &info); err != nil {
		return nil, err
	}
	return &info, err
}

func parseCustomAuthHeader(c *gin.Context, prefix string) (string, bool) {
	auth := c.Request.Header.Get(AuthHeaderCust) // Customized token for exchange
	pf, token := AuthorizationBearer, ""
	if prefix != "" {
		pf = prefix
	}
	if auth != "" && strings.HasPrefix(auth, pf) {
		token = auth[len(pf):]
	}
	return token, token != ""
}

func FindAuthUserIdentities(tk string, body map[string]interface{}) (AuthUserMaster, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(findIdentitiesByUserKey, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil || result.StatusCode() != 200 {
		return AuthUserMaster{}, errors.New(result.String())
	}
	// var ids []user.UserIdentity
	var u AuthUserMaster
	_ = json.Unmarshal(result.Body(), &u)
	return u, nil
}

func ValidateExternalByIdentity(tk, phoneNo, email string) (*ValidateExternalResp, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(map[string]interface{}{
		"phoneNo": phoneNo,
		"email":   email,
	}).Post(fmt.Sprintf(validateExternalByIdentity, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, errors.New("api returns status: " + result.Status())
	}
	resp := ValidateExternalResp{}
	if err := json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func LinkUserWithOneIdentity(tk string, body map[string]interface{}) error {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Patch(fmt.Sprintf(linkUserWithIdentity, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil || result.StatusCode() != 200 {
		return errors.New(result.String())
	}
	return nil
}

func UnlinkUserWithOneIdentity(tk string, body map[string]interface{}) error {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Patch(fmt.Sprintf(unlinkUserWithIdentity, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil || result.StatusCode() != 200 {
		return errors.New(result.String())
	}
	return nil
}

func ResetUserIdentityCredential(tk string, body map[string]interface{}) error {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Patch(fmt.Sprintf(resetUserIdentityCredential, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil || result.StatusCode() != 200 {
		return errors.New(result.String())
	}
	return nil
}

func UpdateAuthUserLockStatus(tk string, userRefKey string, lock bool, isActive *bool) error {
	var body map[string]interface{}
	if isActive == nil {
		body = map[string]interface{}{
			"userKey": userRefKey,
			"lock":    lock,
		}
	} else {
		body = map[string]interface{}{
			"userKey":  userRefKey,
			"lock":     lock,
			"isActive": isActive,
		}
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(updateAuthUserlockStatus, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil || result.StatusCode() != 200 {
		return errors.New(result.String())
	}
	return nil
}

func UpdateAuthUserDeviceRegisterAttempt(tk string, userRefKey string) error {
	body := map[string]interface{}{
		"userKey": userRefKey,
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(updateDeviceIdRegisterAttempt, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil || result.StatusCode() != 200 {
		return errors.New(result.String())
	}
	return nil
}

func FindAllUserLoginHistory(tk string, userRefKey string, isDesc bool) (interface{}, error) {
	body := map[string]interface{}{
		"userKey":    userRefKey,
		"descending": isDesc,
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(findAllLoginHistory, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil || result.StatusCode() != 200 {
		return nil, errors.New(result.String())
	}
	var histories []interface{}
	_ = json.Unmarshal(result.Body(), &histories)
	return histories, nil
}

func GetAuthStatusByUserRefKeys(tk string, userRefKeys []string) (map[string]*AuthUserMaster, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(map[string]interface{}{
		"userRefKeys": userRefKeys,
	}).Post(fmt.Sprintf(getManyUserLockInfo, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, errors.New("api returns status: " + result.Status())
	}
	values := map[string]*AuthUserMaster{}
	if err := json.Unmarshal(result.Body(), &values); err != nil {
		return nil, err
	}
	return values, nil
}

func UpdateAccInactive(tk string) error {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Post(fmt.Sprintf(getUserInactive, apis.V().GetString(apiAuthMdlUrlBase)))
	if err != nil || result.StatusCode() != 200 {
		return errors.New(result.String())
	}
	return nil
}
