package ginrecovery

import (
	"fmt"
	"strings"

	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/Mobility-Development-Team/be-common-mdl/util/apiutil"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

// Returns a gin recovery middleware that logs panic to logrus as well as stdout
// errMsg is returned to the caller after panic recovery
func NewIntercepter(errMsg response.Message) gin.HandlerFunc {
	var writer strings.Builder
	return gin.CustomRecoveryWithWriter(&writer, func(c *gin.Context, err interface{}) {
		msg := writer.String()
		fmt.Print(msg)
		logger.StandardLogger().Log(logger.FatalLevel, msg)
		writer.Reset()
		apiutil.GenerateResponse(c, nil, errMsg)
		c.Abort()
	})
}
