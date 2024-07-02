package hypath

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	apiHypathUrlBase             = "apis.external.hypath.url.base"
	authenticate                 = "%s/auth/authenticate"
	getProjectList               = "%s/confinedspace/ext_permit/projects"
	getCSByProjectCode           = "%s/confinedspace/ext_permit/confinedspace?projectcode=%s"
	getCSBySpaceIdAndProjectCode = "%s/confinedspace/ext_permit/confinedspace/%s?projectcode=%s"
	postCreateCSPermit           = "%s/confinedspace/ext_permit/permit/create"
	postUpdateCSPermit           = "%s/confinedspace/ext_permit/permit/update"
	postUpdateCSPermitWorkflow   = "%s/confinedspace/ext_permit/permit/%s/status/%s"
	// getForm              = "%s/permits/internal/all"
)

var (
	ErrHyPathInvalidCredential    = errors.New("invalid credential")
	ErrHyPathUnableToAuthenticate = errors.New("unable to authenticate")
	ErrHyPathInvalidApiCall       = errors.New("invalid API call")
	ErrHyPathInvalidParam         = errors.New("invalid param")
	ErrHyPathInvalidApiResponse   = errors.New("invalid API response")
)

func AuthenticateHyPath() (result HyPathAuthenResponse, err error) {
	var (
		client     = resty.New()
		un, pw, sp = apis.V().GetString("hypath.username"), apis.V().GetString("hypath.password"), "confinedspace"
		pwDecoded  []byte
	)
	// Handle password decode
	pwDecoded, err = base64.StdEncoding.DecodeString(pw)
	if err != nil {
		err = ErrHyPathInvalidCredential
		return
	}
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err := client.R().SetBody(HyPathAuthenRequest{
		Username: un,
		Password: fmt.Sprintf("%s", pwDecoded),
		Scope:    sp,
	}).Post(
		fmt.Sprintf(authenticate, apis.V().GetString(apiHypathUrlBase)),
	)
	if err != nil || resp.StatusCode() != http.StatusOK {
		err = ErrHyPathUnableToAuthenticate
		return
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		err = ErrHyPathUnableToAuthenticate
		return
	}
	return
}

func GetProjectList(tk string) (result GetProjectListResponse, err error) {
	var (
		client   = resty.New()
		resp     *resty.Response
		authResp HyPathAuthenResponse
	)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	// Get Token if not provided
	if tk == "" {
		authResp, err = AuthenticateHyPath()
		if err != nil || len(authResp.Token) == 0 {
			return
		}
		tk = authResp.Token
	}
	resp, err = client.R().SetAuthToken(tk).Get(
		fmt.Sprintf(getProjectList, apis.V().GetString(apiHypathUrlBase)),
	)
	if err != nil || resp.StatusCode() != http.StatusOK {
		err = ErrHyPathInvalidApiCall
		return
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		err = ErrHyPathInvalidApiResponse
		return
	}
	return
}

func GetConfinedSpaceByProjectCode(tk, projectCode string) (result GetCSByProjectCodeResponse, err error) {
	var (
		client   = resty.New()
		resp     *resty.Response
		authResp HyPathAuthenResponse
	)

	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	if projectCode == "" {
		err = ErrHyPathInvalidParam
		return
	}
	// Get Token if not provided
	if tk == "" {
		authResp, err = AuthenticateHyPath()
		if err != nil || len(authResp.Token) == 0 {
			return
		}
		tk = authResp.Token
	}
	resp, err = client.R().SetAuthToken(tk).Get(
		fmt.Sprintf(getCSByProjectCode, apis.V().GetString(apiHypathUrlBase), projectCode),
	)
	if err != nil || resp.StatusCode() != http.StatusOK {
		err = ErrHyPathInvalidApiCall
		return
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		err = ErrHyPathInvalidApiResponse
		return
	}
	return
}

func GetConfinedSpaceBySpaceIdAndProjectCode(tk, spaceId, projectCode string) (result GetCSBySpaceIdAndProjectCodeResponse, err error) {
	var (
		client   = resty.New()
		resp     *resty.Response
		authResp HyPathAuthenResponse
	)

	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	if spaceId == "" || projectCode == "" {
		err = ErrHyPathInvalidParam
		return
	}
	// Get Token if not provided
	if tk == "" {
		authResp, err = AuthenticateHyPath()
		if err != nil || len(authResp.Token) == 0 {
			return
		}
		tk = authResp.Token
	}
	resp, err = client.R().SetAuthToken(tk).Get(
		fmt.Sprintf(getCSBySpaceIdAndProjectCode, apis.V().GetString(apiHypathUrlBase), spaceId, projectCode),
	)
	if err != nil || resp.StatusCode() != http.StatusOK {
		err = ErrHyPathInvalidApiCall
		return
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		err = ErrHyPathInvalidApiResponse
		return
	}
	return
}

func PostCreateCSPermit(tk string, request PostCreateCSPermitRequest) (result PostCreateCSPermitResponse, url string, requestBody []byte, err error) {
	var (
		client   = resty.New()
		resp     *resty.Response
		authResp HyPathAuthenResponse
	)

	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	// Get Token if not provided
	if tk == "" {
		authResp, err = AuthenticateHyPath()
		if err != nil || len(authResp.Token) == 0 {
			return
		}
		tk = authResp.Token
	}
	// prepare post body for logging
	url = fmt.Sprintf(postCreateCSPermit, apis.V().GetString(apiHypathUrlBase))
	if b, err := json.Marshal(request); err != nil {
		logger.Error("[PublishOneConfinedSpace][PostUpdateCSPermitWorkflow] unable to marshal request")
	} else {
		requestBody = b
	}
	resp, err = client.R().SetAuthToken(tk).SetBody(request).Post(url)
	if err != nil || resp.StatusCode() != http.StatusOK {
		err = ErrHyPathInvalidApiCall
		return
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		err = ErrHyPathInvalidApiResponse
		return
	}
	return
}

func PostUpdateCSPermit(tk string, request PostUpdateCSPermitRequest) (result PostUpdateCSPermitResponse, url string, requestBody []byte, err error) {
	var (
		client   = resty.New()
		resp     *resty.Response
		authResp HyPathAuthenResponse
	)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	// Get Token if not provided
	if tk == "" {
		authResp, err = AuthenticateHyPath()
		if err != nil || len(authResp.Token) == 0 {
			return
		}
		tk = authResp.Token
	}
	// prepare post body for logging
	url = fmt.Sprintf(postUpdateCSPermit, apis.V().GetString(apiHypathUrlBase))
	if b, err := json.Marshal(request); err != nil {
		logger.Error("[UpdateOneConfinedSpaceDSD][PostUpdateCSPermit] unable to marshal request")
	} else {
		requestBody = b
	}
	resp, err = client.R().SetAuthToken(tk).SetBody(request).Post(url)
	if err != nil || resp.StatusCode() != http.StatusOK {
		err = ErrHyPathInvalidApiCall
		return
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		err = ErrHyPathInvalidApiResponse
		return
	}
	return
}

func PostUpdateCSPermitWorkflow(tk, projectFormId, actionType, pdfUrl string) (result PostCommonCSPermitWorkflowResponse, url string, requestBody []byte, err error) {
	var (
		client   = resty.New()
		req      = &PostCommonCSPermitWorkflowRequest{}
		resp     *resty.Response
		authResp HyPathAuthenResponse
	)

	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	// Get Token if not provided
	if tk == "" {
		// Get Token
		authResp, err = AuthenticateHyPath()
		if err != nil || len(authResp.Token) == 0 {
			return
		}
		tk = authResp.Token
	}
	if pdfUrl != "" {
		req.PDFUrl = pdfUrl
	}
	// prepare post body for logging
	url = fmt.Sprintf(postUpdateCSPermitWorkflow, apis.V().GetString(apiHypathUrlBase), projectFormId, strings.ToUpper(actionType))
	if b, err := json.Marshal(req); err != nil {
		logger.Error("[PublishOneConfinedSpace][PostUpdateCSPermitWorkflow] unable to marshal request")
	} else {
		requestBody = b
	}
	// Call hyPath
	resp, err = client.R().SetAuthToken(tk).SetBody(req).Post(url)
	if err != nil || resp.StatusCode() != http.StatusOK {
		// err = ErrHyPathInvalidApiCall
		return
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		err = ErrHyPathInvalidApiResponse
		return
	}
	return
}
