package machine

import (
	"encoding/json"
	"fmt"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/go-resty/resty/v2"
)

const (
	apiMachineMdlUrlBase          = "apis.internal.machine.module.url.base"
	getOnePlantPermit             = "%s/permits/plantpermits/%s"
	getOneNCAPermit               = "%s/permits/nca/%s"
	getOneHotworkPermit           = "%s/permits/hw/%s"
	getOneEXPermit                = "%s/permits/ex/%s"
	getOneELPermit                = "%s/permits/el/%s"
	getOneLA                      = "%s/plant/equip/LA/detail"
	getAllPermits                 = "%s/permits/internal/all"
	getPITChecklist               = "%s/permits/pc/%s"
	getOneTaskRelatedPITChecklist = "%s/permits/pc/checklist/%s"
	getAllAppointmentsForInternal = "%s/permits/appt/internal/all"
	getOneCSPermit                = "%s/permits/cs/%s"
)

func GetOneLA(tk string, criteria LA, isSimple bool) (*LA, error) {
	resp := struct {
		Payload *LA `json:"payload"`
	}{}
	uri := getOneLA
	if isSimple {
		uri += "?isSimple=true"
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(criteria).Post(fmt.Sprintf(uri, apis.V().GetString(apiMachineMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOnePlantPermit(tk string, permitMasterId intstring.IntString) (*PlantPermit, error) {
	resp := struct {
		Payload *PlantPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOnePlantPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOneNCAPermit(tk string, permitMasterId intstring.IntString) (*NCAPermit, error) {
	resp := struct {
		Payload *NCAPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneNCAPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetAllPermits(tk string, userRefKey string, criteria PermitCriteria, opt GetAllPermitOps, preloadNames ...string) ([]*MasterPermit, error) {
	resp := struct {
		Payload []*MasterPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"criteria": criteria,
			"opts":     opt,
			"preloads": preloadNames,
		},
	).Post(
		fmt.Sprintf(getAllPermits, apis.V().GetString(apiMachineMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOneHotworkPermit(tk string, permitMasterId intstring.IntString) (*HotworkPermit, error) {
	resp := struct {
		Payload *HotworkPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneHotworkPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOnePermitToDig(tk string, permitMasterId intstring.IntString) (*EXPermit, error) {
	resp := struct {
		Payload *EXPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneEXPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOneELPermit(tk string, permitMasterId intstring.IntString) (*ELPermit, error) {
	resp := struct {
		Payload *ELPermit `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneELPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOnePITChecklist(tk string, permitMasterId intstring.IntString) (*PITChecklist, error) {
	resp := struct {
		Payload *PITChecklist `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getPITChecklist, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}

	return resp.Payload, err
}

func GetOneCSPermit(tk string, permitMasterId intstring.IntString) (*ConfinedSpacePermit, error) {
	resp := struct {
		Payload *ConfinedSpacePermit  `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(fmt.Sprintf(getOneCSPermit, apis.V().GetString(apiMachineMdlUrlBase), permitMasterId))
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetOneTaskRelatedPITChecklist(tk string, parentGroupId intstring.IntString) (interface{}, error) {
	resp := struct {
		Payload interface{} `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).Get(
		fmt.Sprintf(getOneTaskRelatedPITChecklist, apis.V().GetString(apiMachineMdlUrlBase), parentGroupId),
	)
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func GetAllAppointmentsForMyTask(tk string, criteria PermitApptCriteria) ([]PermitAppointment, error) {
	resp := struct {
		Payload []PermitAppointment `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"criteria": criteria,
		},
	).Post(
		fmt.Sprintf(getAllAppointmentsForInternal, apis.V().GetString(apiMachineMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("machine module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}
