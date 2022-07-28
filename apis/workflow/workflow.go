package workflow

import (
	"encoding/json"
	"errors"
	"fmt"

	"time"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"

	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
)

const (
	apiWorkflowMdlUrlBase = "apis.internal.workflow.module.url.base"
	createWorkflow        = "%s/workflows"
	getLatestWorkflow     = "%s/workflows/tasks/latest"
	submitWorkflowAction  = "%s/workflows/tasks"
)

// The types declared here are for submitting requests
type (
	WorkFlowCreateParam struct {
		WorkflowTemplateKey string      `json:"workflowTemplateKey"`
		TaskParams          []TaskParam `json:"taskParams"`
	}
	Content struct {
		ContentRefID      intstring.IntString `json:"contentRefId"`
		ContentRefTblName string              `json:"contentRefTblName"`
	}
	Actor struct {
		ActorUserID     intstring.IntString `json:"actorUserId"`
		ActorUserRefKey string              `json:"actorUserRefKey"`
		ActorUserType   string              `json:"actorUserType"`
	}
	TaskParam struct {
		TaskUuid      *string   `json:"taskUuid,omitempty"`
		TaskStartDate *string   `json:"taskStartDate"`
		Contents      []Content `json:"contents"`
		Actor         Actor     `json:"actor"`
	}
)

// The types declared here are for parsing response
type (
	WorkflowView struct {
		model.Model
		UUID                    string       `json:"uuid"`
		WorkflowTemplateRefUUID string       `json:"workflowTemplateRefUuid"`
		WorkflowTemplateRefKey  string       `json:"workflowTemplateRefKey"`
		TemplateVersion         int          `json:"templateVersion"`
		Tasks                   []TaskView   `json:"tasks"`
		Actions                 []ActionView `json:"actions"`
	}
	ActionView struct {
		model.Model
		ActionKey           string              `json:"actionKey"`
		ActionGroupUUID     string              `json:"actionGroupUuid"`
		IsPrimary           bool                `json:"isPrimary"`
		IsSystem            bool                `json:"isSystem"`
		ActionTemplateRefID intstring.IntString `json:"actionTemplateRefId"`
		TaskID              intstring.IntString `json:"taskId"`
	}
	ActorView struct {
		model.Model
		ActorUserID     string               `json:"actorUserId"`
		ActorUserRefKey string               `json:"actorUserRefKey"`
		ActorUserType   string               `json:"actorUserType"`
		ActorGroupID    *intstring.IntString `json:"actorGroupId"`
		TaskID          intstring.IntString  `json:"taskId"`
	}
	ContentView struct {
		model.Model
		ContentRefID      string              `json:"contentRefId"`
		ContentRefTblName string              `json:"contentRefTblName"`
		TaskID            intstring.IntString `json:"taskId"`
	}
	TaskView struct {
		model.Model
		UUID                 string              `json:"uuid"`
		Status               string              `json:"status"`
		TaskKey              string              `json:"taskKey"`
		TaskTemplateRefID    string              `json:"taskTemplateRefId"`
		TaskTemplateRefUUID  string              `json:"taskTemplateRefUuid"`
		TaskStartDate        time.Time           `json:"taskStartDate"`
		TaskTargetEndDate    time.Time           `json:"taskTargetEndDate"`
		TaskActualEndDate    *time.Time          `json:"taskActualEndDate"`
		TaskEstimatedDayCost int                 `json:"taskEstimatedDayCost"`
		IsCurrentTask        bool                `json:"isCurrentTask"`
		IsFinished           bool                `json:"isFinished"`
		WorkflowID           intstring.IntString `json:"workflowId"`
		Actions              []ActionView        `json:"actions"`
		Actors               []ActorView         `json:"actors"`
		Contents             []ContentView       `json:"contents"`
		SubsequentTasks      []TaskView          `json:"subsequentTasks"`
	}
)

func (a *ActionView) UnmarshalJSON(b []byte) error {
	type alias ActionView
	var resultStruct alias
	if err := json.Unmarshal(b, &resultStruct); err == nil {
		// Unmarshal successful
		*a = ActionView(resultStruct)
		return nil
	}
	var resultString string
	if err := json.Unmarshal(b, &resultString); err != nil {
		return errors.New("not a struct nor a string")
	}
	*a = ActionView{
		ActionKey: resultString,
	}
	return nil
}

func CreateWorkflow(tk string, action WorkFlowCreateParam) (*WorkflowView, error) {
	type (
		respType struct {
			response.Response
			Payload *WorkflowView `json:"payload"`
		}
	)
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		action,
	).Post(fmt.Sprintf(createWorkflow, apis.V().GetString(apiWorkflowMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	var resp respType
	if err := json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("API returns status: %s message: %s %s", result.Status(), resp.MsgCode, resp.Message)
	}
	return resp.Payload, nil
}

func GetLatestWorkflowTask(tk, workflowUuid string) (*WorkflowView, error) {
	type (
		respType struct {
			response.Response
			Payload []WorkflowView `json:"payload"`
		}
	)
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(map[string]string{
		"workflowUuid": workflowUuid,
	}).Post(fmt.Sprintf(getLatestWorkflow, apis.V().GetString(apiWorkflowMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	var resp respType
	if err := json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("API returns status: %s message: %s %s", result.Status(), resp.MsgCode, resp.Message)
	}
	if len(resp.Payload) == 0 {
		return nil, nil
	}
	if len(resp.Payload) > 1 {
		logger.Warn("[GetLatestWorkflow] API returned more than 1 results, using first one: %+v", resp.Payload)
	}
	return &resp.Payload[0], nil
}

type WorkflowActionParam struct {
	TaskParam
	SelectedAction string `json:"selectedAction"`
}

func SubmitWorkflowAction(tk string, actions []WorkflowActionParam) (map[string][]ActionView, error) {
	var resp struct {
		response.Response
		Payload map[string][]ActionView `json:"payload"`
	}
	client := resty.New()
	result, err := client.R().
		SetAuthToken(tk).
		SetBody(actions).
		Post(fmt.Sprintf(submitWorkflowAction, apis.V().GetString(apiWorkflowMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("API returns status: %s message: %s %s", result.Status(), resp.MsgCode, resp.Message)
	}
	return resp.Payload, nil
}
