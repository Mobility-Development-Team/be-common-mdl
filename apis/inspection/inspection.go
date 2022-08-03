package inspection

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/apis/inspection/models"
	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
)

const (
	apiInspectionMdlUrlBase    = "apis.internal.inspection.module.url.base"
	getUserPendingAppointments = "%s/inspection/tasks/appointments/pending/current"
)

// Gets all appointments requiring the user's attention
//
// Setting isSimple to ture skip preloading of some fields, setting it to false ensures the appointment is fully populated.
// However, the nested fields inside the sitewalk object is never fully populated
func FindUserPendingAppointments(tk string, userRefKey string, isSimple bool) ([]models.Appointment, error) {
	client := resty.New()
	result, err := client.R().
		SetAuthToken(tk).
		SetQueryParam("isSimple", strconv.FormatBool(isSimple)).
		Get(fmt.Sprintf(getUserPendingAppointments, apis.V().GetString(apiInspectionMdlUrlBase)))
	if err != nil {
		logger.Error("[FindUserPendingAppointments] err: ", err)
		return nil, err
	}
	type respType struct {
		response.Response
		Payload []models.Appointment `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		logger.Error("[FindUserPendingAppointments] Unmarshal err:", err)
		return nil, err
	}
	return resp.Payload, nil
}
