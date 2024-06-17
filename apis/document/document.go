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
	urlBase                  = "apis.internal.document.module.url.base"
	generateSiteWalk         = "%s/documents/inspection/sitewalk/report/generate"
	generateAdminSiteWalk    = "%s/documents/inspection/sitewalk/adminReport/generate"
	generateRATSiteWalk      = "%s/documents/inspection/sitewalk/rat/generate"
	generateFollowUpReport   = "%s/documents/inspection/tasks/followup/generate"
	generatePlantCertificate = "%s/documents/machine/permits/plantpermits/cert/generate"
	generatePlantReport      = "%s/documents/machine/permits/plantpermits/report/generate"
	generateNCAReport        = "%s/documents/machine/permits/nca/report/generate"
	generateHWReport         = "%s/documents/machine/permits/hw/report/generate"
	generateEXReport         = "%s/documents/machine/permits/ex/report/generate"
	generateELReport         = "%s/documents/machine/permits/el/report/generate"
	generatePCReport         = "%s/documents/machine/permits/pc/report/generate"
	generatePCCertificate    = "%s/documents/machine/permits/pc/cert/generate"
	generateCSReport         = "%s/documents/machine/permits/cs/report/generate"
	generateLDReport         = "%s/documents/machine/permits/ld/report/generate"
	generateEFReport         = "%s/documents/machine/permits/ef/report/generate"
	generateLSReport         = "%s/documents/machine/permits/ls/report/generate"
	generateCDReport         = "%s/documents/machine/permits/cd/report/generate"
	generateDocReport        = "%s/reports/generate"
)

func GenerateSiteWalk(tk string, siteWalkId intstring.IntString) (string, error) {
	return generateReportSiteWalk(tk, generateSiteWalk, siteWalkId, true)
}

func GenerateRAT(tk string, siteWalkId intstring.IntString) (string, error) {
	return generateReportSiteWalk(tk, generateRATSiteWalk, siteWalkId, true)
}

func GenerateSiteWalkAdmin(tk string, siteWalkId intstring.IntString) (string, error) {
	return generateReportSiteWalk(tk, generateAdminSiteWalk, siteWalkId, true)
}

func GenerateTaskFollowUpReport(tk string, params FollowUpReportInfo, taskId intstring.IntString, contractId intstring.IntString) (string, error) {
	client := resty.New()
	var resp struct {
		Payload struct {
			Url string `json:"url"`
		} `json:"payload"`
	}
	result, err := client.R().SetAuthToken(tk).SetBody(struct {
		FollowUpReportInfo
		TaskId     intstring.IntString `json:"taskId"`
		ContractId intstring.IntString `json:"contractId"`
	}{
		FollowUpReportInfo: params,
		TaskId:             taskId,
		ContractId:         contractId,
	}).Post(
		fmt.Sprintf(generateFollowUpReport, apis.V().GetString(urlBase)),
	)
	if err != nil {
		return "", err
	}
	if result.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf(
			"[GenerateTaskFollowUpReport] status code not 201: %d", result.StatusCode(),
		)
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return "", err
	}
	return resp.Payload.Url, nil
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
	}).Post(fmt.Sprintf(apiPath, apis.V().GetString(urlBase)))
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

func GeneratePermitCertificate(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generatePlantCertificate, permitMasterId, true)
}

func GeneratePCCertificate(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generatePCCertificate, permitMasterId, true)
}

func GeneratePlantReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generatePlantReport, permitMasterId, true)
}

func GenerateNCAReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generateNCAReport, permitMasterId, true)
}

func GenerateHWReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generateHWReport, permitMasterId, true)
}

func GenerateEXReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generateEXReport, permitMasterId, true)
}

func GenerateELReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generateELReport, permitMasterId, true)
}

func GeneratePCReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generatePCReport, permitMasterId, true)
}

func GenerateCSReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generateCSReport, permitMasterId, true)
}

func GenerateLDReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generateLDReport, permitMasterId, true)
}

func GenerateEFReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generateEFReport, permitMasterId, true)
}

func GenerateLSReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generateLSReport, permitMasterId, true)
}


func GenerateCDReport(tk string, permitMasterId intstring.IntString) (string, error) {
	return generatePermitType(tk, generateCDReport, permitMasterId, true)
}


func generatePermitType(tk string, apiPath string, permitMasterId intstring.IntString, publish bool) (string, error) {
	client := resty.New()
	var resp struct {
		Payload struct {
			Url string `json:"url"`
		} `json:"payload"`
	}
	result, err := client.R().SetAuthToken(tk).SetBody(map[string]interface{}{
		"permitMasterId": permitMasterId,
		"publish":        publish,
	}).Post(fmt.Sprintf(apiPath, apis.V().GetString(urlBase)))
	if err != nil {
		return "", err
	}
	if result.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("[generatePlantPermitType] status code not 201: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return "", err
	}
	return resp.Payload.Url, nil
}

func GenerateDocReport(tk string, reportId intstring.IntString) (string, error) {
	return generateDoc(tk, generateDocReport, reportId, true)
}

func generateDoc(tk, apiPath string, reportId intstring.IntString, publish bool) (string, error) {
	client := resty.New()
	var resp struct {
		Payload struct {
			Url string `json:"url"`
		} `json:"payload"`
	}
	result, err := client.R().SetAuthToken(tk).SetBody(map[string]interface{}{
		"reportId": reportId,
		"publish":  publish,
	}).Post(fmt.Sprintf(apiPath, apis.V().GetString(urlBase)))
	if err != nil {
		return "", err
	}
	if result.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("[GenerateDocReport] status code not 201: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return "", err
	}
	return resp.Payload.Url, nil
}
