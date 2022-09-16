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
)

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
		return nil, fmt.Errorf("system module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}
