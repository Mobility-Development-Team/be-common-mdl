package notification

import (
	"fmt"
	"strings"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/apis/core"
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

func handleExtraNotification(tk string, notifications ...*Notification) (map[intstring.IntString][]intstring.IntString, error) {
	// 筛选出需要额外发送的user组合成map
	contractUserIDMap := make(map[intstring.IntString][]intstring.IntString)
	var extraContractIDs []intstring.IntString
	// 循环取出需要去检查的项目
	for _, notification := range notifications {
		if notification.EnableExtraSendUser {
			extraContractIDs = append(extraContractIDs, *notification.ContractID)
		}
	}
	if len(extraContractIDs) == 0 {
		return contractUserIDMap, nil
	}
	// 查询配置的 contract 是否能开启发送额外用户，有组装map返回
	req := map[string]interface{}{
		"contractIds":         extraContractIDs,
		"contractIsExtraSend": ContractIsExtraSendOn,
		"userIsExtraSend":     UserIsExtraSendOn,
		"userStatus":          UserStatusActive,
		"contractStatus":      ContractStatusActive,
	}
	userDetails, err := core.GetManyContractMapUsers(
		tk, req)
	if err != nil {
		return nil, err
	}
	for _, userDetail := range userDetails {
		contractId := userDetail.ContractToUserMap.ContractId
		contractUserIDMap[contractId] = append(contractUserIDMap[contractId], userDetail.ContractToUserMap.UserId)
	}
	return contractUserIDMap, nil
}

// checkSuffix 判断哪些是需要根据额外发送配置去发送额外用户的
func checkSuffix(typeName string) bool {
	suffixes := []string{"_APPROVED", "_ACCEPTED", "_REJECTED", "_CANCELLED"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(typeName, suffix) {
			return true
		}
	}
	return false
}
func CreateNotifications(tk string, notifications ...*Notification) error {
	validNotifications := make([]*Notification, 0, len(notifications))
	for _, notification := range notifications {
		notification.SetupEnableExtraSendUser(checkSuffix(notification.TemplateType))
	}
	contractUserIDMap, err := handleExtraNotification(tk, notifications...)
	if err != nil {
		logger.Warn(fmt.Sprintf("[CreateNotification] Skipped: extra notifcation has error: %v", err))
		return err
	}
	for _, noti := range notifications {
		if noti == nil {
			logger.Warn("[CreateNotification] Skipped: notifcation is nil.")
			continue
		}
		if noti.TemplateType == "" {
			logger.Warn("[CreateNotification] Skipped: notifcation does not have a templateType: ", *noti)
			continue
		}
		if !noti.PermitEmptyRecipients && len(noti.Recipients.Groups) == 0 && len(noti.Recipients.Users) == 0 && len(noti.Recipients.PartyAdmin) == 0 {
			logger.Debug("[CreateNotification] Notification has no recipents, ignoring: ", noti.TemplateType)
			continue
		}
		// 项目开启配置有配置额外发送人且GetExtraSendUser is True额外发送
		if extraUserIds, exist := contractUserIDMap[*noti.ContractID]; exist && noti.EnableExtraSendUser {
			noti.AddUserRecipient(extraUserIds...)
		}

		validNotifications = append(validNotifications, noti)
	}
	if len(validNotifications) == 0 {
		logger.Debug("[CreateNotification] No notifications to create.")
		return nil
	}
	logger.Debugf("[CreateNotification] Creating %d notification(s).", len(validNotifications))
	return createOneOrManyNotifications(tk, validNotifications)
}

func createOneOrManyNotifications(tk string, body interface{}) error {
	client := resty.New()
	// Set retries (base on legor requiremnt since duplicate notifications)
	// client.
	// 	SetRetryCount(5).
	// 	SetRetryWaitTime(5 * time.Minute).
	// 	SetRetryMaxWaitTime(10 * time.Minute).
	// 	AddRetryCondition(func(r *resty.Response, err error) bool {
	// 		if !r.IsSuccess() {
	// 			logger.Warnf("[CreateNotification]: API returns status code %d. Retrying...", r.StatusCode())
	// 			return true
	// 		}
	// 		return false
	// 	})
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
