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
