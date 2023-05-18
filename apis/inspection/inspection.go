package inspection

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
)

const (
	apiInspectionMdlUrlBase    = "apis.internal.inspection.module.url.base"
	getUserPendingAppointments = "%s/inspection/tasks/appointments/pending/current"
	getSiteWalkInfo            = "%s/inspection/tasks/%s"
	registerAttachment         = "%s/inspection/tasks/attachments"
	getSiteWalkActivityLog     = "%s/inspection/tasks/activities/all"
	getAllTasks                = "%s/tasks/all"
	getSitePlanBySiteWalkId    = "%s/inspection/tasks/siteplans/latest"
	getFollowUpByParentRefIds  = "%s/tasks/followup/tasks/all/many"
	findManyTaskByParentId     = "%s/inspection/tasks/parents"
)

// Gets all appointments requiring the user's attention
//
// Setting isSimple to ture skip preloading of some fields, setting it to false ensures the appointment is fully populated.
// However, the nested fields inside the sitewalk object is never fully populated
func FindUserPendingAppointments(tk string, userRefKey string, isSimple bool) ([]Appointment, error) {
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
		Payload []Appointment `json:"payload"`
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

func GetSitePlanBySiteWalkId(tk string, siteWalkId intstring.IntString) (*SitePlanDisplay, error) {
	client := resty.New()
	result, err := client.R().
		SetAuthToken(tk).
		SetBody(map[string]interface{}{
			"siteWalkId": siteWalkId,
		}).
		Post(fmt.Sprintf(getSitePlanBySiteWalkId, apis.V().GetString(apiInspectionMdlUrlBase)))
	if err != nil {
		logger.Error("[GetSitePlan] err: ", err)
		return nil, err
	}
	var resp struct {
		response.Response
		Payload *SitePlanDisplay `json:"payload"`
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("inspection module returned status code: %d", result.StatusCode())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		logger.Error("[GetSitePlan] Unmarshal err:", err)
		return nil, err
	}
	resp.Payload.ShouldAddSystemFieldsFromDisplay()
	return resp.Payload, nil
}

func GetLatestFollowUpTasksByParentRefIds(tk string, taskParentRefIds ...intstring.IntString) (map[intstring.IntString]*FollowUpTaskDisplay, error) {
	if len(taskParentRefIds) == 0 {
		return map[intstring.IntString]*FollowUpTaskDisplay{}, nil
	}
	client := resty.New()
	result, err := client.R().
		SetAuthToken(tk).
		SetBody(map[string]interface{}{
			"taskParentRefIds": taskParentRefIds,
		}).
		Post(fmt.Sprintf(getFollowUpByParentRefIds, apis.V().GetString(apiInspectionMdlUrlBase)))
	if err != nil {
		logger.Error("[GetLatestTasksByParentRefIds] err: ", err)
		return nil, err
	}
	var resp struct {
		response.Response
		Payload map[intstring.IntString]*FollowUpTaskDisplay `json:"payload"`
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("inspection module returned status code: %d", result.StatusCode())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		logger.Error("[GetLatestTasksByParentRefIds] Unmarshal err:", err)
		return nil, err
	}
	for _, k := range resp.Payload {
		k.ShouldAddSystemFieldsFromDisplay()
	}
	return resp.Payload, nil
}

func GetAllTasks(tk string, cri GetAllTasksCriteria) ([]TaskDisplay, error) {
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
			Tasks      []TaskDisplay `json:"tasks"`
			TotalCount int           `json:"totalCount"`
		} `json:"payload"`
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("inspection module returned status code: %d", result.StatusCode())
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		logger.Error("[FindAllTasks] Unmarshal err:", err)
		return nil, err
	}
	for i := range resp.Payload.Tasks {
		resp.Payload.Tasks[i].ShouldAddSystemFieldsFromDisplay()
	}
	return resp.Payload.Tasks, nil
}

func GetSiteWalkDetail(tk string, siteWalkId intstring.IntString) (*SiteWalk, error) {
	resp := struct {
		Payload *SiteWalk `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getSiteWalkInfo, apis.V().GetString(apiInspectionMdlUrlBase), siteWalkId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("inspection module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func RegisterAttachment(tk string, attachment Attachment) (interface{}, error) {
	resp := struct {
		Payload interface{} `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).
		SetBody(attachment).
		Post(fmt.Sprintf(registerAttachment, apis.V().GetString(apiInspectionMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	if result.StatusCode() != 201 {
		return nil, fmt.Errorf("inspection module returned status code not 201: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetSiteWalkActivityLog(tk string, siteWalkId, checklistId *intstring.IntString) ([]ActivityLog, error) {
	resp := struct {
		Payload []ActivityLog `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().
		SetAuthToken(tk).
		SetBody(struct {
			SiteWalkId  *intstring.IntString `json:"siteWalkId,omitempty"`
			ChecklistId *intstring.IntString `json:"checklistId,omitempty"`
			Descending  bool                 `json:"descending"`
		}{
			SiteWalkId:  siteWalkId,
			ChecklistId: checklistId,
			Descending:  false,
		}).
		Post(fmt.Sprintf(getSiteWalkActivityLog, apis.V().GetString(apiInspectionMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("inspection module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func FindManyTaskByParentId(tk string, parentId, parentGroupId *intstring.IntString, parentType *string) (map[intstring.IntString]interface{}, error) {
	resp := struct {
		Payload map[intstring.IntString]interface{} `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().
		SetAuthToken(tk).
		SetBody(map[string]interface{}{
			"parentId":      parentId,
			"parentGroupId": parentGroupId,
			"parentType":    parentType,
		}).
		Post(fmt.Sprintf(findManyTaskByParentId, apis.V().GetString(apiInspectionMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("inspection module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}
