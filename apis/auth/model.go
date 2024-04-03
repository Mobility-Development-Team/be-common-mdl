package auth

import "time"

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
		Uuid              string                      `json:"userKey" binding:"required"`
		Status            string                      `json:"status" binding:"required"`
		LastFailedLogin   *time.Time                  `json:"lastFailedLogin"`
		LastSuccessLogin  *time.Time                  `json:"lastSuccessLogin"`
		LoginAttempt      int                         `json:"loginAttempt"`
		DeviceRegdAttempt bool                        `json:"deviceRegdAttempt"`
		HasProfile        bool                        `json:"hasProfile"`
		IsLocked          bool                        `json:"isLocked"`
		IsApiAccount      bool                        `json:"isApiAccount"`
	}

)
