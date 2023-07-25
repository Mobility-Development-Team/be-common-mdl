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
	apiAuthMdlUrlBase       = "apis.internal.auth.module.url.base"
	getTokenInfo            = "%s/tokeninfo"
	validateEmatToken       = "%s/validate/smm/user"
	validateEmatTokenWithTk = "%s/validate/smm/user?tk=true"
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
