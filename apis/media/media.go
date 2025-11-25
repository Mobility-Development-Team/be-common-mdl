package media

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Mobility-Development-Team/be-common-mdl/apis"
	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/response"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"

	"github.com/go-resty/resty/v2"
	logger "github.com/sirupsen/logrus"
)

const (
	apiMediaMdlUrlBase          = "apis.internal.media.module.url.base"
	getMediaMany                = "%s/media/all"
	getMediaManySimple          = "%s/media/many/simple"
	getNoAuthUsersFirebaseToken = "%s/fb/custom/token"
	getMediaManyByRefId         = "%s/media/many?showAsMap=true"
	getBatchMany                = "%s/media/batch/many"
	uploadSitePlanPicture       = "%s/file/upload/siteplan"
	cloneMediaToBatch           = "%s/media/batch/clone"
	sendCloudMessage            = "%s/fcm/messaging"
	getFileKeys                 = "%s/file/upload/worker-mgt/file/k"
	uploadUrlBase               = "%s/file/upload/%s/%s"
	uploadFileUrlBase           = "%s/file/upload/%s"

	previewFolderName   = "preview"
	publishedFolderName = "publish"
)

const (
	ReportTypeSiteWalk            = "siteWalk"
	ReportTypeSiteWalkAdmin       = "siteWalkAdmin"
	ReportTypeTaskFollowUp        = "taskFollowUpReport"
	ReportTypeRat                 = "rat"
	ReportTypePlantPermitCert     = "plantPermitCert"
	ReportTypePlantPermitReport   = "plantPermitReport"
	ReportTypeNCAPermitReport     = "ncaPermitReport"
	ReportTypeHotworkPermitReport = "hwPermitReport"
	ReportTypeEXPermitReport      = "exPermitReport"
	ReportTypeELPermitReport      = "elPermitReport"
	ReportTypePCPermitCert        = "pcPermitCert"
	ReportTypePCPermitReport      = "pcPermitReport"
	ReportTypeCSPermitReport      = "csPermitReport"
	ReportTypeCDPermitReport      = "cdPermitReport"
	ReportTypeCDV2PermitReport    = "cdv2PermitReport"
	ReportTypeLDPermitReport      = "ldPermitReport"
	ReportTypeLSPermitReport      = "lsPermitReport"
	ReportTypeEFPermitReport      = "efPermitReport"
	ReportTypeDocReport           = "docReport"
)

type CloneOpts struct {
	NoTruncate bool `json:"noTruncate"` // Do not delete any records
	NoClone    bool `json:"noClone"`    // Do not clone new records
	NoUpdate   bool `json:"noUpdate"`   // Do not update existing record
	ForReport  bool `json:"forReport"`
}

type Scope map[string]string
type CloneMediaParams struct {
	Params  []model.MediaParam
	Scope   Scope
	BatchId string
}

func GetManySimpleMedia(tk string, body map[string]interface{}) ([]model.SimpleMediaItems, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(getMediaManySimple, apis.V().GetString(apiMediaMdlUrlBase)))
	if err != nil {
		logger.Error("[GetManySimpleMedia]", "err:", err)
		return nil, err
	}
	type respType struct {
		response.Response
		Payload []model.SimpleMediaItems `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		logger.Error("[GetManySimpleMedia]", "Unmarshal err:", err)
		return nil, err
	}
	return resp.Payload, nil
}

func GetMedia(tk string, body map[string]interface{}) ([]model.MediaParam, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(body).Post(fmt.Sprintf(getMediaMany, apis.V().GetString(apiMediaMdlUrlBase)))
	if err != nil {
		logger.Error("[GetMedia]", "err:", err)
		return nil, err
	}
	type respType struct {
		response.Response
		Payload []model.MediaParam `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		logger.Error("[getMedia]", "Unmarshal err:", err)
		return nil, err
	}
	return resp.Payload, nil
}

func GetUsersFirebaseToken(tk string, body map[string]string) (*model.UsersFirebaseToken, error) {
	client := resty.New()
	url := apis.V().GetString(apiMediaMdlUrlBase)
	result, err := client.R().SetAuthToken(tk).SetQueryParams(body).Get(fmt.Sprintf(getNoAuthUsersFirebaseToken, url))
	if err != nil {
		logger.Error("[GetUsersFirebaseToken]", "err:", err)
		return nil, err
	}
	type respType struct {
		response.Response
		Payload *model.UsersFirebaseToken `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		logger.Error("[GetUsersFirebaseToken]", "Unmarshal err:", err)
		return nil, err
	}
	return resp.Payload, nil
}

func GetMediaByRefId(tk string, refId ...string) (map[string]model.MediaParam, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(map[string][]string{
		"ids": refId,
	}).Post(fmt.Sprintf(getMediaManyByRefId, apis.V().GetString(apiMediaMdlUrlBase)))
	if err != nil {
		logger.Error("[GetMediaBatches] err: ", err)
		return nil, err
	}
	type respType struct {
		response.Response
		Payload map[string]model.MediaParam `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		logger.Error("[GetMediaBatches] Unmarshal err: ", err)
		return nil, err
	}
	return resp.Payload, nil
}

func GetMediaBatches(tk string, batchId ...string) (map[string][]model.MediaParam, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(map[string][]string{
		"batchIds": batchId,
	}).Post(fmt.Sprintf(getBatchMany, apis.V().GetString(apiMediaMdlUrlBase)))
	if err != nil {
		logger.Error("[GetMediaBatches] err: ", err)
		return nil, err
	}
	type respType struct {
		response.Response
		Payload map[string][]model.MediaParam `json:"payload"`
	}
	var resp respType
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		logger.Error("[GetMediaBatches] Unmarshal err: ", err)
		return nil, err
	}
	return resp.Payload, nil
}

func GetMediaByBatchId(tk string, batchId string) ([]model.MediaParam, error) {
	return GetMedia(tk, map[string]interface{}{
		"batchId": batchId,
	})
}

// Media in different batches can refer to the same siteWalkId.
// This function would return all of them
func GetMediaBySiteWalkId(tk string, siteWalkId intstring.IntString) ([]model.MediaParam, error) {
	return GetMedia(tk, map[string]interface{}{
		"mediaRefInfo": map[string]interface{}{
			"siteWalkId": siteWalkId,
		},
	})
}

func GetMediaByNcId(tk string, ncFindingId intstring.IntString) ([]model.MediaParam, error) {
	return GetMedia(tk, map[string]interface{}{
		"mediaRefInfo": map[string]interface{}{
			"ncFindingId":  ncFindingId,
			"taskActionId": "",
		},
		"includeEmptyValues": true,
	})
}

func GetMediaByTaskId(tk string, taskId intstring.IntString) ([]model.MediaParam, error) {
	return GetMedia(tk, map[string]interface{}{
		"mediaRefInfo": map[string]interface{}{
			"taskId":       taskId,
			"taskActionId": "",
		},
		"includeEmptyValues": true,
	})
}

func GetMediaByTaskActionId(tk string, taskActionId intstring.IntString, taskActionType ...string) ([]model.MediaParam, error) {
	refInfo := map[string]interface{}{
		"taskActionId": taskActionId,
	}
	if len(taskActionType) > 0 {
		refInfo["taskActionType"] = taskActionType[0]
	}
	return GetMedia(tk, map[string]interface{}{
		"mediaRefInfo": refInfo,
	})
}

func ShouldGetMediaByTaskActionId(tk string, taskActionId intstring.IntString, taskActionType ...string) []model.MediaParam {
	m, err := GetMediaByTaskActionId(tk, taskActionId, taskActionType...)
	if err != nil {
		logger.Errorf("[ShouldGetMediaByTaskActionId] Unable to get media: action=%s err=%v", taskActionId, err)
	}
	if m == nil {
		return []model.MediaParam{}
	}
	return m
}

func GetMediaByGeneralFindingId(tk string, generalFindingId intstring.IntString) ([]model.MediaParam, error) {
	return GetMedia(tk, map[string]interface{}{
		"mediaRefInfo": map[string]interface{}{
			"generalFindingId": generalFindingId,
		},
	})
}

func GetMediaByChecklistId(tk string, checklistId intstring.IntString) ([]model.MediaParam, error) {
	return GetMedia(tk, map[string]interface{}{
		"mediaRefInfo": map[string]interface{}{
			"checklistId": checklistId,
		},
	})
}

func MapMediaByChecklistItemId(media []model.MediaParam, errOpt ...error) (map[intstring.IntString][]model.MediaParam, error) {
	var err error
	if len(errOpt) > 0 {
		err = errOpt[0]
	}
	return MapMediaById(media, err, func(m model.MediaParam) intstring.IntString {
		return m.ShouldGetRefInfo().ShouldGetIntString("checklistItemId")
	})
}

func MapMediaByGeneralFindingId(media []model.MediaParam, errOpt ...error) (map[intstring.IntString][]model.MediaParam, error) {
	var err error
	if len(errOpt) > 0 {
		err = errOpt[0]
	}
	return MapMediaById(media, err, func(m model.MediaParam) intstring.IntString {
		return m.ShouldGetRefInfo().ShouldGetIntString("generalFindingId")
	})
}

func MapMediaByNcFindingId(media []model.MediaParam, errOpt ...error) (map[intstring.IntString][]model.MediaParam, error) {
	var err error
	if len(errOpt) > 0 {
		err = errOpt[0]
	}
	return MapMediaById(media, err, func(m model.MediaParam) intstring.IntString {
		return m.ShouldGetRefInfo().ShouldGetIntString("ncFindingId")
	})
}

func MapMediaById(media []model.MediaParam, err error, mapfunc func(m model.MediaParam) intstring.IntString) (
	map[intstring.IntString][]model.MediaParam, error,
) {
	if err != nil {
		return nil, err
	}
	mediaMap := map[intstring.IntString][]model.MediaParam{}
	for _, m := range media {
		id := mapfunc(m)
		mediaMap[id] = append(mediaMap[id], m)
	}
	return mediaMap, nil
}

func MapMediaByUuid(media []model.MediaParam, err error, mapfunc func(m model.MediaParam) string) (
	map[string][]model.MediaParam, error,
) {
	if err != nil {
		return nil, err
	}
	mediaMap := map[string][]model.MediaParam{}
	for _, m := range media {
		id := mapfunc(m)
		mediaMap[id] = append(mediaMap[id], m)
	}
	return mediaMap, nil
}

// CloneMediaToBatch Clones the media into a specified batch, overwriting exsiting batch with truncation if possible.
//
// `scope` determines the subset of media to be affected, if scope is specified, only existing media inside the batch that
// have all matching key-value pair inside their refInfo will be updated/deleted as required. If a specified media is being
// put inside the batch (specified in `media`) but is currently outside the scope (belongs to a different batchId, or does
// not have the matching key-value pairs), the record duplicated and put inside the given `batchId“, leaving the original
// media record intact.
func CloneMediaToBatch(tk string, batchId string, media []model.MediaParam, scope Scope, optOpts ...CloneOpts) error {
	var opts *CloneOpts
	if len(optOpts) > 0 {
		opts = &optOpts[0]
	}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(struct {
		BatchId string             `json:"batchId"`
		Media   []model.MediaParam `json:"media"`
		Scope   map[string]string  `json:"scope"`
		Opts    *CloneOpts         `json:"opts,omitempty"`
	}{
		BatchId: batchId,
		Media:   media,
		Scope:   scope,
		Opts:    opts,
	}).Post(fmt.Sprintf(cloneMediaToBatch, apis.V().GetString(apiMediaMdlUrlBase)))
	if err != nil {
		logger.Error("[CloneMediaToBatch]", "err:", err)
		return err
	}
	if !result.IsSuccess() {
		logger.Error("[CloneMediaToBatch] media module returns status code: ", result.Status())
		return errors.New("media module returns status code: " + result.Status())
	}
	return nil
}

// A quick helper function for conveniently adding additional restriction to the scope
func (s Scope) AddRestriction(key, value string) Scope {
	s[key] = value
	return s
}

func (s Scope) AddTaskActionTypeRestriction(value string) Scope {
	return s.AddRestriction("taskActionType", value)
}

func UploadSitePlanPicture(tk string, fileName string, imgBytes []byte) (*string, error) {
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetHeader("Content-Type", "multipart/form-data;charset=UTF-8").
		SetFileReader("file", fileName, bytes.NewReader(imgBytes)).
		Post(fmt.Sprintf(uploadSitePlanPicture, apis.V().GetString(apiMediaMdlUrlBase)))
	if err != nil {
		return nil, err
	}
	var resp struct {
		response.Response
		Payload *string `json:"payload"`
	}
	if err = json.Unmarshal(result.Body(), &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// UploadReport Uploads a site walk report
// contractId must be specified
// If publish mode is false  fileName will be used as file name instead
// If publish mode is true, the fileName specified would be ignored by the media module
func UploadReport(tk string, file io.Reader, reportType string, contractId intstring.IntString, fileName string, publish bool) (string, error) {
	client := resty.New()
	var resp struct {
		Payload string `json:"payload"`
	}
	var reqUri string
	if publish {
		reqUri = fmt.Sprintf(uploadUrlBase, apis.V().GetString(apiMediaMdlUrlBase), reportType, publishedFolderName)
	} else {
		reqUri = fmt.Sprintf(uploadUrlBase, apis.V().GetString(apiMediaMdlUrlBase), reportType, previewFolderName)
	}
	// The filename specified here would only be used when it is in preview mode (publish == false)
	fileName = fmt.Sprintf("preview-file-%s.pdf", fileName)
	result, err := client.R().SetAuthToken(tk).
		SetFileReader("file", fileName, file).
		SetFormData(map[string]string{
			"contractId": contractId.String(),
		}).
		Post(reqUri)
	if err != nil {
		return "", err
	}
	if !result.IsSuccess() {
		return "", fmt.Errorf("media module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return "", err
	}
	return resp.Payload, err
}

// UploadFile Uploads permit reference doc
func UploadFile(tk string, fileBytes []byte, fileName string, reportType string, contractId intstring.IntString) (string, error) {
	client := resty.New()
	var resp struct {
		Payload string `json:"payload"`
	}
	// The filename specified here would only be used when it is in preview mode (publish == false)
	// fileName = fmt.Sprintf("preview-file-%s.pdf", fileName)
	result, err := client.R().SetAuthToken(tk).
		SetHeader("Content-Type", "multipart/form-data;charset=UTF-8").
		SetFileReader("file", fileName, bytes.NewReader(fileBytes)).
		SetFormData(map[string]string{
			"contractId": contractId.String(),
		}).
		Post(fmt.Sprintf(uploadFileUrlBase, apis.V().GetString(apiMediaMdlUrlBase), reportType))
	if err != nil {
		return "", err
	}
	if !result.IsSuccess() {
		return "", fmt.Errorf("media module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return "", err
	}
	return resp.Payload, err
}

// Get file Keys  permit reference doc
func GetFileKeys(tk string, urls []string) (map[string]string, error) {
	resp := struct {
		Payload map[string]string `json:"payload"`
	}{}
	client := resty.New()
	result, err := client.R().SetAuthToken(tk).SetBody(
		map[string]interface{}{
			"urls": urls,
		},
	).Post(
		fmt.Sprintf(getFileKeys, apis.V().GetString(apiMediaMdlUrlBase)),
	)
	if err != nil {
		return nil, err
	}
	if !result.IsSuccess() {
		return nil, fmt.Errorf("photo module returned status code: %d", result.StatusCode())
	}
	err = json.Unmarshal(result.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}
