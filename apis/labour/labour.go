package labour

import (
	"encoding/json"
	"fmt"
	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/go-resty/resty/v2"
)

const (
	apiLabourMdlUrlBase         = "apis.internal.labour.module.url.base"
	getAllUnsafeCasesForMyTasks = "%s/unsafe-cases/all/internal"
)

func GetAllUnsafeCasesForMyTasks(tk string, criteria UnsafeCaseCriteria) ([]*UnsafeCase, error) {
	resp := struct {
		Payload []*UnsafeCase `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"criteria": criteria,
			// "opts":     opt,
			// "preloads": preloadNames,
		},
	).Post(
		fmt.Sprintf(getAllUnsafeCasesForMyTasks, apis.V().GetString(apiLabourMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("labour module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}
