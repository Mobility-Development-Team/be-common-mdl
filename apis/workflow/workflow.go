package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/response"

	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
)

const (
	apiWorkflowMdlUrlBase = "apis.internal.workflow.module.url.base"
	createWorkflow        = "%s/workflows"
	getLatestWorkflow     = "%s/workflows/tasks/latest"
	submitWorkflowAction  = "%s/workflows/tasks"
)

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
