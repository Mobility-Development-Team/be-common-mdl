package auth

import (
	"time"

	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
)

type (
	AuthUserMaster struct {
		LastFailedLogin   *time.Time    `json:"lastFailedLogin"`
		LastSuccessLogin  *time.Time    `json:"lastSuccessLogin"`
		LoginAttempt      int           `json:"loginAttempt"`
		DeviceRegdAttempt bool          `json:"deviceRegdAttempt"`
		IsLocked          bool          `json:"isLocked"`
		IsApiAccount      bool          `json:"isApiAccount"`
		Identities        []interface{} `json:"identities"`
	}
	ValidateExternalResp struct {
		IsValid bool   `json:"isValid"`
		Uuid    string `json:"userKey"`
	}

	UserMaster struct {
		Uuid             string     `json:"userKey" binding:"required"`
		Status           string     `json:"status" binding:"required"`
		LastFailedLogin  *time.Time `json:"lastFailedLogin"`
		LastSuccessLogin *time.Time `json:"lastSuccessLogin"`
	}

	LoginHistory struct {
		Id            intstring.IntString `json:"id"`
		CreatedAt     string              `json:"createdAt"`
		CreatedBy     string              `json:"createdBy"`
		LogType       string              `json:"logType"`
		LogMessage    string              `json:"logMessage"`
		LogMessageZh  string              `json:"LogMessageZh"`
		LogDeviceType string              `json:"logDeviceType"`
		UserId        intstring.IntString `json:"userId"`
	}
)
