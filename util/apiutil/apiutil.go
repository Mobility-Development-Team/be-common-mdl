package apiutil

import (
	"strings"

	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/gin-gonic/gin"
)

func ParseBearerAuth(c *gin.Context) (string, bool) {
	auth := c.Request.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token, token != ""
}

func GenerateResponse(c *gin.Context, payload interface{}, message response.Message, v ...interface{}) {
	var resp interface{}
	if v != nil {
		resp = response.NewResponseByMessage(payload, message, v)
	} else {
		resp = response.NewResponseByMessage(payload, message)
	}
	c.JSON(message.StatusCode, resp)
	c.Abort()
}
