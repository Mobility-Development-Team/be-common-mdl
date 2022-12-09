package machine

import (
	"encoding/json"
	"fmt"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/go-resty/resty/v2"
)

const (
	apiMachineMdlUrlBase = "apis.internal.machine.module.url.base"
	getOnePlantPermit    = "%s/permits/plantpermits/%s"
	getOneNCAPermit      = "%s/permits/nca/%s"
	getOneHotworkPermit  = "%s/permits/hw/%s"
	getOneEXPermit       = "%s/permits/ex/%s"
	getOneELPermit       = "%s/permits/el/%s"
	getOneAsset          = "%s/permits/assets/internal/getone"
	getAllPermits        = "%s/permits/internal/all"
)

func GetOneAsset(tk string, criteria Equipment, isSimple bool) (*Equipment, error) {
	resp := struct {
		Payload *Equipment `json:"payload"`
	}{}
	uri := getOneAsset
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
