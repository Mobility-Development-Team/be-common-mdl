package inspection

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/apis/inspection/models"
	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
)

const (
	apiInspectionMdlUrlBase    = "apis.internal.inspection.module.url.base"
	getUserPendingAppointments = "%s/inspection/tasks/appointments/pending/current"
	getAllTasks                = "%s/tasks/all"
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

// This structure does not cover all optional paramters that can be passed to GetAllTasks
// It may be expanded if necessary
type GetAllTasksCriteria struct {
	SiteWalkId   *intstring.IntString `json:"siteWalkId"`
	ContractId   *intstring.IntString `json:"contractId"`
	SearchType   string               `json:"searchType"`
	TaskStatuses []string             `json:"taskStatuses"`
}

func GetAllTasks(tk string, cri GetAllTasksCriteria) ([]models.TaskDisplay, error) {
	if cri.SiteWalkId == nil && cri.ContractId == nil && cri.SearchType == "" {
		return nil, errors.New("invalid parameters: no search constraint")
	}
	client := resty.New()
	result, err := client.R().
		SetAuthToken(tk).
		SetBody(cri).
		Post(fmt.Sprintf(getAllTasks, apis.V().GetString(apiInspectionMdlUrlBase)))
	if err != nil {
		logger.Error("[GetAllTasks] err: ", err)
		return nil, err
	}
	var resp struct {
		response.Response
		Payload struct {
			Tasks      []models.TaskDisplay `json:"tasks"`
			TotalCount int                  `json:"totalCount"`
		} `json:"payload"`
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		logger.Error("[FindAllTasks] Unmarshal err:", err)
		return nil, err
	}
	return resp.Payload.Tasks, nil
}
