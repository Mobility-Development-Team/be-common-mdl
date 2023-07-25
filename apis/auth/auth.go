package auth

import (
	"encoding/json"
	"errors"
	"fmt"

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
	CustAuthHeader      = "Authorization-ext"
	// API constant
	apiAuthMdlUrlBase = "apis.internal.auth.module.url.base"
	getTokenInfo      = "%s/tokeninfo"
	validateEmatToken = "%s/validate/smm/user"
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
	}
)

func GetUserRefKeyFromContext(c *gin.Context) string {
	k, ok := c.Get(tokenInfoUserRefKey) // Assume
	if !ok {
		return "undefined"
	}
	return fmt.Sprintf("%v", k)
}

// Gets a gin middleware for handing token verifications. The returned intercepter should be registered
// as a middleware in gin for protected API calls such that ParseBearerAuth and GetUserRefKeyFromContext
// can be used. invalidHeaderMsg or invalidTokenMsg is returned to the user in case of error.
func NewTokenVerifierInterceptor(invalidHeaderMsg, invalidTokenMsg response.Message) gin.HandlerFunc {
	return func(c *gin.Context) {
		tk, ok := apiutil.ParseBearerAuth(c)
		if !ok {
			logger.Warn("[ValidateInternalToken] unable to parse the given token from the header")
			apiutil.GenerateResponse(c, nil, invalidHeaderMsg)
			c.Abort()
		}
		info, err := GetTokenInfo(c, tk)
		if err != nil {
			logger.Warn("[ValidateInternalToken] ", err)
			apiutil.GenerateResponse(c, nil, invalidTokenMsg)
			c.Abort()
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
	result, err := client.R().SetHeader(keyTokenInfo, fmt.Sprintf("Bearer %s", tk)).Get(fmt.Sprintf(validateEmatToken, apis.V().GetString(apiAuthMdlUrlBase)))
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
