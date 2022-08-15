package document

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"

	"github.com/go-resty/resty/v2"
)

const (
	apiNotificationMdlUrlBase = "apis.internal.document.module.url.base"
	generateSiteWalk          = "%s/documents/inspection/sitewalk/report/generate"
	generateRATSiteWalk       = "%s/documents/inspection/sitewalk/rat/generate"
)

func GenerateSiteWalk(tk string, siteWalk intstring.IntString) (string, error) {
	return generateReportSiteWalk(tk, generateSiteWalk, siteWalk, true)
}

func GenerateRAT(tk string, siteWalk intstring.IntString) (string, error) {
	return generateReportSiteWalk(tk, generateRATSiteWalk, siteWalk, true)
}

func generateReportSiteWalk(tk, apiPath string, id intstring.IntString, publish bool) (string, error) {
	client := resty.New()
	var resp struct {
		Payload struct {
			Url string `json:"url"`
		} `json:"payload"`
	}
	result, err := client.R().SetAuthToken(tk).SetBody(map[string]interface{}{
		"id":      id,
		"publish": publish,
	}).Post(fmt.Sprintf(apiPath, apis.V().GetString(apiNotificationMdlUrlBase)))
	if err != nil {
		return "", err
	}
	if result.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("[GenerateSiteWalk] status code not 201: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return "", err
	}
	return resp.Payload.Url, nil
}
