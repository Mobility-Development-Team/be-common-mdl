package apiutil

import (
	"strings"

	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/gin-gonic/gin"
)

const (
	HeaderRegular         = "Authorization"
	HeaderCustom          = "Authorization-ext"
	AuthHeaderPrefixBasic = "Basic "
)

func ParseBearerAuth(c *gin.Context) (string, bool) {
	auth := c.Request.Header.Get(HeaderRegular)
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token, token != ""
}

func ParseCustAuthExt(c *gin.Context, prefix string) (string, bool) {
	auth := c.Request.Header.Get(HeaderCustom)
	pf, token := "Basic ", ""
	if prefix != "" {
		pf = prefix
	}
	if auth != "" && strings.HasPrefix(auth, pf) {
		token = auth[len(pf):]
	}
	return token, token != ""
}

func GenerateResponse(c *gin.Context, payload interface{}, message response.Message, v ...interface{}) {
	var resp interface{}
	if v != nil {
		resp = response.NewResponseByMessage(payload, message, v...)
	} else {
		resp = response.NewResponseByMessage(payload, message)
	}
	c.JSON(message.StatusCode, resp)
	c.Abort()
}

func CommonResultIndicator(isSuccess bool) map[string]interface{} {
	return map[string]interface{}{
		"isSuccess": isSuccess,
	}
}
