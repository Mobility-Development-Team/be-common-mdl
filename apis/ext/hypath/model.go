package hypath

type (
	HyPathApiBase struct {
		Success    bool        `json:"success"`
		StatusCode int         `json:"statusCode"`
		Error      interface{} `json:"error"`
	}
	HyPathAuthenRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Scope    string `json:"scope"`
	}
	HyPathAuthenResponse struct {
		Success bool        `json:"success"`
		Error   interface{} `json:"error"`
		Token   string      `json:"token"`
	}

	GetProjectListResponse struct {
		HyPathApiBase
		Data []GetProjectListDetail `json:"data"`
	}
	GetProjectListDetail struct {
		ProjectCode   string `json:"projectCode"`
		DescriptionZH string `json:"descriptionZH"`
		DescriptionEN string `json:"descriptionEN"`
	}

	GetCSByProjectCodeResponse struct {
		HyPathApiBase
		Data []GetCSByProjectCodeResponseDetail `json:"data"`
	}
	GetCSByProjectCodeResponseDetail struct {
		Id           string  `json:"id"`
		ProjectCode  string  `json:"projectCode"`
		WorkLocation string  `json:"workLocation"`
		Northing     float64 `json:"northing"`
		Easting      float64 `json:"easting"`
	}

	GetCSBySpaceIdAndProjectCodeResponse struct {
		HyPathApiBase
		Data GetCSBySpaceIdAndProjectCodeResponseDetail `json:"data"`
	}
	GetCSBySpaceIdAndProjectCodeResponseDetail struct {
		Sensors []Sensor `json:"sensors"`
		GetCSByProjectCodeResponseDetail
	}

	Sensor struct {
		SensorType            string  `json:"sensorType"`
		Name                  string  `json:"name"`
		Value                 float64 `json:"value"`
		DisplayUnit           string  `json:"displayUnit"`
		SensorModel           string  `json:"sensorModel"`
		SensorSerialNo        string  `json:"sensorSerialNo"`
		ValueCollectDateTime  string  `json:"valueCollectDateTime"`
		CalibrationExpiryDate string  `json:"calibrationExpiryDate"`
		Location              string  `json:"location"`
	}
	PostCreateCSPermitRequest struct {
		ProjectCode     string   `json:"projectCode"`
		ConfinedSpaceId string   `json:"confinedSpaceId"`
		StartDateTime   string   `json:"startDateTime"` // could be time.Time
		EndDateTime     string   `json:"endDateTime"`   // could be time.Time
		PDFUrl          string   `json:"PDFUrl"`
		Workers         []Worker `json:"workers"`
	}
	PostCreateCSPermitResponse struct {
		HyPathApiBase
		Data PostCreateCSPermitResponseResponseDetail `json:"data"`
	}
	PostCreateCSPermitResponseResponseDetail struct {
		PermitFormId string `json:"permitFormId"`
	}
	Worker struct {
		WorkerType string `json:"workerType"`
		MappingKey string `json:"mappingKey"`
		IsSZWorker bool   `json:"isSZWorker"`
	}

	// PostUpdateCSPermitRequest identical with PostCreateCSPermitRequest
	// PostUpdateCSPermitResponse identical with PostCreateCSPermitResponse
	// PostUpdateCSPermitResponseResponseDetail identical with PostCreateCSPermitResponseResponseDetail
	// For EMAT-7791
	PostUpdateCSPermitRequest struct {
		ProjectCode     string   `json:"projectCode"`
		ConfinedSpaceId string   `json:"confinedSpaceId"`
		StartDateTime   string   `json:"startDateTime"` // could be time.Time
		EndDateTime     string   `json:"endDateTime"`   // could be time.Time
		PDFUrl          string   `json:"PDFUrl"`
		Workers         []Worker `json:"workers"`
	}
	PostUpdateCSPermitResponse struct {
		HyPathApiBase
		Data PostUpdateCSPermitResponseResponseDetail `json:"data"`
	}
	PostUpdateCSPermitResponseResponseDetail struct {
		PermitFormId string `json:"permitFormId"`
	}

	PostCommonCSPermitWorkflowRequest struct {
		PDFUrl string `json:"PDFUrl"`
	}
	PostCommonCSPermitWorkflowResponse struct {
		HyPathApiBase
		Data interface{} `json:"data"`
	}
)
