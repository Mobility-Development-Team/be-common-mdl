package hypath

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

const (
	apiHypathUrlBase             = "apis.external.hypath.url.base"
	authenticate                 = "%s/auth/authenticate"
	getCSByProjectCode           = "%s/confinedspace/ext_permit/confinedspace?projectcode=%s"
	getCSBySpaceIdAndProjectCode = "%s/confinedspace/ext_permit/confinedspace/%s?projectcode=%s"
	postCreateCSPermit           = "%s/confinedspace/ext_permit/permit/create"
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

func GetConfinedSpaceByProjectCode(tk, projectCode string) (result GetCSByProjectCodeResponse, err error) {
	var (
		client   = resty.New()
		resp     *resty.Response
		authResp HyPathAuthenResponse
	)
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
		fmt.Sprintf(getCSByProjectCode, projectCode, apis.V().GetString(apiHypathUrlBase)),
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
		fmt.Sprintf(getCSBySpaceIdAndProjectCode, spaceId, projectCode, apis.V().GetString(apiHypathUrlBase)),
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

func PostCreateCSPermit(tk string, request PostCreateCSPermitRequest) (result PostCreateCSPermitResponse, err error) {
	var (
		client   = resty.New()
		resp     *resty.Response
		authResp HyPathAuthenResponse
	)
	// Get Token if not provided
	if tk == "" {
		authResp, err = AuthenticateHyPath()
		if err != nil || len(authResp.Token) == 0 {
			return
		}
		tk = authResp.Token
	}
	resp, err = client.R().SetAuthToken(tk).SetBody(request).Post(
		fmt.Sprintf(postCreateCSPermit, apis.V().GetString(apiHypathUrlBase)),
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

func PostUpdateCSPermitWorkflow(tk, projectFormId, actionType, pdfUrl string) (result PostCommonCSPermitWorkflowResponse, err error) {
	var (
		client   = resty.New()
		resp     *resty.Response
		authResp HyPathAuthenResponse
	)
	// Get Token if not provided
	if tk == "" {
		// Get Token
		authResp, err = AuthenticateHyPath()
		if err != nil || len(authResp.Token) == 0 {
			return
		}
		tk = authResp.Token
	}
	resp, err = client.R().SetAuthToken(tk).SetBody(&PostCommonCSPermitWorkflowRequest{
		pdfUrl,
	}).Post(
		fmt.Sprintf(postUpdateCSPermitWorkflow, projectFormId, strings.ToUpper(actionType), apis.V().GetString(apiHypathUrlBase)),
	)
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
