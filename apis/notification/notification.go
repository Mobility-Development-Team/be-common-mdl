package notification

import (
	"fmt"

	"time"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/apis/notification/models"
	"github.com/Mobility-Development-Team/be-common-mdl/apis/user"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
)

const (
	apiNotificationMdlUrlBase = "apis.internal.notification.module.url.base"
	createNotification        = "%s/notification"
)

const (
	notificationTypeSystem = "SYSTEM"
)

func FilterSelfFromUserIds(c *gin.Context, ids []intstring.IntString) []intstring.IntString {
	self, err := user.GetCurrentUserInfoFromContext(c)
	if err != nil {
		logger.Errorf("[FilterSelfFromUserIds] Unable to retrieve current user info: %s, current user id is not filtered.", err.Error())
		return ids
	}
	result := make([]intstring.IntString, 0, len(ids))
	for _, id := range ids {
		if id != self.Id {
			result = append(result, id)
		}
	}
	return result
}

func NewNotification(contractId *intstring.IntString, templateType string, templateParams ...interface{}) *models.Notification {
	return &models.Notification{
		NotificationType: notificationTypeSystem,
		TemplateType:     templateType,
		ContractID:       contractId,
		Params:           templateParams,
	}
}

func CreateNotifications(tk string, notifications ...*models.Notification) error {
	if len(notifications) == 0 {
		logger.Debug("[CreateNotification]: No notifications to create.")
		return nil
	}
	logger.Debugf("[CreateNotification]: Creating %d notification(s).", len(notifications))
	return createOneOrManyNotifications(tk, notifications)
}

func createOneOrManyNotifications(tk string, body interface{}) error {
	client := resty.New()
	// Set retries
	client.
		SetRetryCount(5).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if !r.IsSuccess() {
				logger.Warnf("[CreateNotification]: API returns status code %d. Retrying...", r.StatusCode())
				return true
			}
			return false
		})
	// Send request
	result, err := client.R().
		SetAuthToken(tk).
		SetBody(body).
		Post(fmt.Sprintf(createNotification, apis.V().GetString(apiNotificationMdlUrlBase)))
	if err != nil {
		return err
	}
	if !result.IsSuccess() {
		return fmt.Errorf("[CreateNotification] given up: API returned status code %d", result.StatusCode())
	}
	return nil
}

// This function requires env.short to be set in config
func AppendRefKeyWithEnvironment(refKey string) string {
	return fmt.Sprintf("%s-%s", apis.V().GetString("env.short"), refKey)
}
