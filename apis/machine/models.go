package machine

import (
	"time"

	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
)

type (
	PlantPermit struct {
		model.Model
		MasterPermit
		PlantType            string              `json:"plantType"`
		PlantOwner           *string             `json:"plantOwner"`
		PlantOwnerName       *string             `json:"plantOwnerName"`
		Manufacturer         *string             `json:"manufacturer"`
		ModelNo              *string             `json:"modelNo"`
		YearOfManufacture    *string             `json:"yearOfManufacture"`
		CertExpiryDate       *string             `json:"certExpiryDate"`
		CertActualExpiryDate *string             `json:"certActualExpiryDate"`
		SerialNo             *string             `json:"serialNo"`
		OwnerNo              *string             `json:"ownerNo"`
		IsRental             *bool               `json:"isRental"`
		RejectionReason      *string             `json:"rejectionReason"`
		PermitMasterId       intstring.IntString `json:"permitMasterId"`
		// Custom fields
		CurrentApprovalStage int `json:"currentApprovalStage"`
	}
	NCAPermit struct {
		MasterPermit
		// WorkDate string                        `json:"workDate"`
		// WorkPeriodFromTime string              `json:"workPeriodFromTime"`
		// WorkPeriodToTime   string              `json:"workPeriodToTime"`
		Applicant          *ApplicantDisplay   `json:"applicant"`
		WorkLocation       *string             `json:"workLocation"`
		PermitMasterId     intstring.IntString `json:"permitMasterId"`
		WorkPeriodsDisplay []WorkPeriodDisplay `json:"workPeriods"`
		ConstructionWorks  []ConstructionWork  `json:"constructionWorks"`
		PrescribedWorks    []PrescribedWork    `json:"prescribedWorks"`
		ProjectPermitInfo  []ProjectPermitInfo `json:"projectPermitInfo"`
		MechEquipments     []MechEquipment     `json:"mechEquipments"`
		Workers            []Worker            `json:"workers"`
	}

	PITChecklist struct {
		MasterPermit
		PlantType            string              `json:"plantType"`
		PlantOwner           *string             `json:"plantOwner"`
		PlantOwnerName       *string             `json:"plantOwnerName"`
		Manufacturer         *string             `json:"manufacturer"`
		ModelNo              *string             `json:"modelNo"`
		YearOfManufacture    *string             `json:"yearOfManufacture"`
		CertExpiryDate       *string             `json:"certExpiryDate"`
		CertActualExpiryDate *string             `json:"certActualExpiryDate"`
		SerialNo             *string             `json:"serialNo"`
		OwnerNo              *string             `json:"ownerNo"`
		IsRental             *bool               `json:"isRental"`
		IsRpe                *bool               `json:"isRpe"`
		PermitMasterId       intstring.IntString `json:"permitMasterId"`
	}

	WorkPeriodDisplay struct {
		model.Model
		WorkDate       string               `json:"workDate"`
		WorkPeriodFrom string               `json:"workPeriodFrom"`
		WorkPeriodTo   string               `json:"workPeriodTo"`
		NoiseControlId *intstring.IntString `json:"noiseControlId"`
	}

	HotworkPermit struct {
		MasterPermit
		WorkDurationFromDate string              `json:"workDurationFromDate"`
		WorkDurationFromTime string              `json:"workDurationFromTime"`
		WorkDurationToDate   string              `json:"workDurationToDate"`
		WorkDurationToTime   string              `json:"workDurationToTime"`
		Applicant            *ApplicantDisplay   `json:"applicant"`
		HotWorkType          string              `json:"hotWorkType"`
		WorkLocation         *string             `json:"workLocation"`
		PermitMasterId       intstring.IntString `json:"permitMasterId"`
		CancelMedia          []model.MediaParam  `json:"cancelMedia"`
	}

	ConfinedSpacePermit struct {
		MasterPermit
		RraReportNo    string              `json:"rraReportNo"`
		Workers        []CSWorker          `json:"workers"`
		WorkActivity   string              `json:"workActivity"`
		WorkDate       string              `json:"workDate"`
		WorkLocation   *string             `json:"workLocation"`
		WorkPeriodFrom string              `json:"workPeriodFrom"`
		WorkPeriodTo   string              `json:"workPeriodTo"`
		PermitMasterId intstring.IntString `json:"permitMasterId"`
	}

	EXPermit struct {
		MasterPermit
		WorkDurationFromDate string              `json:"workDurationFromDate"`
		WorkDurationFromTime string              `json:"workDurationFromTime"`
		WorkDurationToDate   string              `json:"workDurationToDate"`
		WorkDurationToTime   string              `json:"workDurationToTime"`
		Applicant            *ApplicantDisplay   `json:"applicant"`
		PermitToDigType      string              `json:"permitToDigType"`
		WorkLocation         *string             `json:"workLocation"`
		PermitMasterId       intstring.IntString `json:"permitMasterId"`
	}
	ELPermit struct {
		MasterPermit
		CertExpiryDate       string            `json:"certExpiryDate"`
		LiftingWorkers       []LiftingWorker   `json:"liftingWorkers"`
		LiftingGears         []LiftingGear     `json:"liftingGears"`
		WorkDurationFromDate string            `json:"workDurationFromDate"`
		WorkDurationFromTime string            `json:"workDurationFromTime"`
		WorkDurationToDate   string            `json:"workDurationToDate"`
		WorkDurationToTime   string            `json:"workDurationToTime"`
		WorkLocation         *string           `json:"workLocation"`
		CraneType            *string           `json:"craneType"`
		CraneTypeRemark      *string           `json:"craneTypeRemark"`
		CraneSerialNo        *string           `json:"craneSerialNo"`
		Applicant            *ApplicantDisplay `json:"applicant"`
	}

	EFPermit struct {
		MasterPermit
		LiftingWorkers        []EFLiftingWorker `json:"liftingWorkers"`
		LiftingGears          []EFLiftingGear   `json:"liftingGears"`
		WorkDurationFromDate  string            `json:"workDurationFromDate"`
		WorkDurationFromTime  string            `json:"workDurationFromTime"`
		WorkDurationToDate    string            `json:"workDurationToDate"`
		WorkDurationToTime    string            `json:"workDurationToTime"`
		WorkLocation          *string           `json:"workLocation"`
		Weather               *string           `json:"weather"`
		WorkDescription       *string           `json:"workDescription"`
		PlantType             *string           `json:"plantType"`
		PlantTypeRemark       *string           `json:"plantTypeRemark"`
		ModelNo               *string           `json:"modelNo"`
		SerialNo              *string           `json:"serialNo"`
		MethodStatementRef    *string           `json:"methodStatementRef"`
		JhaRefNo              *string           `json:"jhaRefNo"`
		MaxLoadWeightKg       *string           `json:"maxLoadWeightKg"`
		MaxLoadWeightM        *string           `json:"maxLoadWeightM"`
		MaxRadiusWeightKg     *string           `json:"maxRadiusWeightKg"`
		MaxRadiusWeightM      *string           `json:"maxRadiusWeightM"`
		IsCapableToLift       bool              `json:"isCapableToLift"`
		IsDistApplicable      bool              `json:"isDistApplicable"`
		IsDistApplicableM     *string           `json:"isDistApplicableM"`
		LaCraneValidGreenTick bool              `json:"laCraneValidGreenTick"`
		LaCraneSitePlanNo     *string           `json:"laCraneSitePlanNo"`
		LaCraneCertExpiryDate string            `json:"laCraneCertExpiryDate"`
		LaWoGtCertNo          *string           `json:"laWoGtCertNo"`
		LaWoGtCertExpiryDate  string            `json:"laWoGtCertExpiryDate"`
		LaCraneLorryType      *string           `json:"laCraneLorryType"`
		Applicant             *ApplicantDisplay `json:"applicant"`
	}

	LDPermit struct {
		MasterPermit
		WorkLocation         *string             `json:"workLocation"`
		WorkDurationFromDate string              `json:"workDurationFromDate"`
		WorkDurationFromTime string              `json:"workDurationFromTime"`
		WorkDurationToDate   string              `json:"workDurationToDate"`
		WorkDurationToTime   string              `json:"workDurationToTime"`
		PermitMasterId       intstring.IntString `json:"permitMasterId"`
	}

	LSPermit struct {
		MasterPermit
		WorkActivity         string              `json:"workActivity"`
		WorkLocation         *string             `json:"workLocation"`
		WorkFloor            string              `json:"workFloor"`
		WorkDurationFromDate string              `json:"workDurationFromDate"`
		WorkDurationFromTime string              `json:"workDurationFromTime"`
		WorkDurationToDate   string              `json:"workDurationToDate"`
		WorkDurationToTime   string              `json:"workDurationToTime"`
		PermitMasterId       intstring.IntString `json:"permitMasterId"`
		Workers              []LSWorker          `json:"workers"`
	}

	CDPermit struct {
		MasterPermit
		PitDepth                  string               `json:"pitDepth"`
		WorkLocation              string               `json:"workLocation"`
		WorkDate                  string               `json:"workDate"`
		GasDetectorNo             string               `json:"gasDetectorNo"`
		Detectorvaliduntil        string               `json:"detectorvaliduntil"`
		GasSerialNo               string               `json:"gasSerialNo"`
		GasExpiryDate             *string              `json:"gasExpiryDate"`
		SafetyAdvise              string               `json:"safetyAdvise"`
		InvolveUgPipeWkPtOne      bool                 `json:"involveUgPipeWkPtOne"`
		EquipBreathApptPtOne      bool                 `json:"equipBreathApptPtOne"`
		RequireSafetyMeasurePtOne bool                 `json:"requireSafetyMeasurePtOne"`
		Safetymeasureverified     bool                 `json:"safetymeasureverified"`
		NoOfIdKeptPtOne           intstring.IntString  `json:"noOfIdKeptPtOne"`
		NoOfIdKeptPtTwo           intstring.IntString  `json:"noOfIdKeptPtTwo"`
		WorkNature                string               `json:"workNature"`
		WorkerCompany             string               `json:"workerCompany"`
		InvolveUgPipeWkPtTwo      bool                 `json:"involveUgPipeWkPtTwo"`
		EquipBreathApptPtTwo      bool                 `json:"equipBreathApptPtTwo"`
		RequireSafetyMeasurePtTwo bool                 `json:"requireSafetyMeasurePtTwo"`
		RequireNoWindCoalPipeWork bool                 `json:"requireNoWindCoalPipeWork"`
		Windcoldpipeworkverified  bool                 `json:"windcoldpipeworkverified"`
		SuperviseDate             *string              `json:"superviseDate"`
		CertValidFrom             string               `json:"certValidFrom"`
		CertValidTo               string               `json:"certValidTo"`
		CertValidHour             intstring.IntString  `json:"certValidHour"`
		DcwEnterTime              string               `json:"dcwEnterTime"`
		DcwDepartTime             string               `json:"dcwDepartTime"`
		DisplayStatus             string               `json:"displayStatus"`
		IsAcknowledged            bool                 `json:"isAcknowledged"`
		AcknowledgedBy            *string              `json:"acknowledgedBy"`
		AcknowledgedAt            *string              `json:"acknowledgedAt"`
		// CmpConnected              bool                 `json:"cmpConnected"`
		// CmpPermitRefId            *intstring.IntString `json:"cmpPermitRefId"`
		// CmpSpaceRefId             *intstring.IntString `json:"cmpSpaceRefId"`
		ContractId                intstring.IntString  `json:"contractId"`
		EligiblePersons           []EigiblePersons     `json:"eligiblePersons"`
		DetectiveReports          []DetectiveReports   `json:"detectiveReports"`
		EmergencyContacts         []EmergencyContacts  `json:"emergencyContacts"`
		PermitMasterId            intstring.IntString  `json:"permitMasterId"`
		IcItems                   []IcItems            `json:"icItems"`
		Inspectors                []Inspectors         `json:"inspectors"`
		Workers                   []DSDWorker          `json:"workers"`
		Media                     []model.MediaParam   `json:"media"`
	}

	PermitAppointment struct {
		model.Model
		ApptStatus      string              `json:"apptStatus"`
		ReceiverType    string              `json:"receiverType"`
		Remark          *string             `json:"remark"`
		WorkflowRefUuid *string             `json:"workflowRefUuid"`
		PermitMasterId  intstring.IntString `json:"permitMasterId"`
		TimeSlots       []interface{}       `json:"timeSlots"`
		Participants    []*interface{}      `json:"participants"`
		PermitNo        string              `json:"permitNo"`
		PermitStatus    string              `json:"permitStatus"`
		PlantType       string              `json:"plantType"`
		ContractRefId   intstring.IntString `json:"contractRefId"`
	}
	MasterPermit struct {
		model.Model
		PermitNo           string              `json:"permitNo"`
		PermitType         string              `json:"permitType"`
		PermitBatchRefUuid string              `json:"permitBatchRefUuid"`
		PermitStatus       string              `json:"permitStatus"`
		ContractRefId      intstring.IntString `json:"contractId"`
		WorkflowRefUuid    *string             `json:"workflowRefUuid"`
		Checklists         []Checklist         `json:"checklists"`
		Participants       []Participant       `json:"participants"`
		ApprovalFlow       []PermitApproval    `json:"approvalFlow"`
		Logs               []ActivityLog       `json:"logs"`
		Attachments        []Attachment        `json:"attachments"`
		ApprovalStage      *string             `json:"approvalStage"`
		SuppInfos          []PermitSuppInfo    `json:"suppInfos"`
	}
	Checklist struct {
		model.Model
		TemplateRefKey       string              `json:"templateRefKey"`
		ChecklistNameEn      string              `json:"checklistNameEn"`
		ChecklistNameZh      string              `json:"checklistNameZh"`
		IsCompleted          bool                `json:"isCompleted"`
		PermitMasterId       intstring.IntString `json:"permitMasterId"`
		TemplateRefOwnerType string              `json:"templateRefOwnerType"`
		Items                []ChecklistItem     `json:"items"`
	}
	Attachment struct {
		model.Model
		AttachmentName      string              `json:"attachmentName"`
		AttachmentUrl       string              `json:"attachmentUrl"`
		AttachmentType      string              `json:"attachmentType"`
		AttachmentExtension string              `json:"attachmentExtension"`
		AttachmentMimeType  string              `json:"attachmentMimeType"`
		AttachmentGroupUuid string              `json:"attachmentGroupUuid"`
		Version             int                 `json:"version"`
		PermitMasterId      intstring.IntString `json:"permitMasterId"`
	}
	ActivityLog struct {
		Id              intstring.IntString  `gorm:"primaryKey" json:"id"`
		CreatedAt       time.Time            `json:"createdAt"`
		CreatedBy       string               `json:"createdBy" gorm:"column:created_by"`
		ActorUserId     *intstring.IntString `json:"-"`
		ActorUserRefKey *string              `json:"-"`
		Actor           *model.UserInfo      `json:"actor"`
		Message         *string              `json:"message"`
		MessageZh       *string              `json:"messageZh"`
		ActivityType    string               `json:"activityType"`
		PermitMasterId  intstring.IntString  `json:"permitMasterId"`
	}
	Participant struct {
		model.UserInfo
		ParticipantType string               `json:"participantType"`
		PartyRefId      *intstring.IntString `json:"partyRefId"`
		UserSource      string               `json:"userSource"`
		Party           *model.PartyInfo     `json:"party"`
		SignatureB64    *string              `json:"signature,omitempty"`
	}
	PermitApproval struct {
		model.Model
		SubmittedBy           *model.UserInfo `json:"submittedBy"`
		SubmittedByActionType string          `json:"submittedByActionType"`
		IsCompleted           bool            `json:"isCompleted"`
		IsRejected            bool            `json:"isRejected"`
		IsCurrent             bool            `json:"isCurrent"`
		SubmittedAt           *time.Time      `json:"submittedAt"`
		Seq                   int             `json:"seq"`
	}
	ChecklistItem struct {
		model.Model
		Seq                   int                        `json:"seq"`
		Response              string                     `json:"response"`
		ItemNameEn            string                     `json:"itemNameEn"`
		ItemNameZh            string                     `json:"itemNameZh"`
		ResponsedBy           intstring.IntString        `json:"responsedBy"`
		TemplateItemRefId     intstring.IntString        `json:"templateItemRefId"`
		PermitChecklistId     intstring.IntString        `json:"permitChecklistId"`
		ResponsedByUserRefKey string                     `json:"responsedByUserRefKey"`
		HasRemark             bool                       `json:"hasRemark"`
		HasMedia              bool                       `json:"hasMedia"`
		ItemRemarkNameEn      string                     `json:"itemRemarkNameEn"`
		ItemRemarkNameZh      string                     `json:"itemRemarkNameZh"`
		ItemGrpNameEn         string                     `json:"itemGrpNameEn"`
		ItemGrpNameZh         string                     `json:"itemGrpNameZh"`
		IsMandatory           bool                       `json:"isMandatory"`
		Remark                string                     `json:"remark"`
		Media                 []model.MediaParam         `json:"media" gorm:"-"`
		ChecklistSuppInfos    []ChecklistItemSupportInfo `json:"suppInfos"`
	}
	PermitSuppInfo struct {
		model.Model
		SuppInfoKey      string              `json:"suppInfoKey"`
		SuppInfoVal      *string             `json:"suppInfoVal"`
		SuppInfoDataType string              `json:"suppInfoDataType"`
		PermitMasterId   intstring.IntString `json:"permitMasterId"`
		// Media is hidden via omitempty unless explicitly set, hence a pointer here
		Media *[]model.MediaParam `json:"media,omitempty"`
	}
	ChecklistItemSupportInfo struct {
		model.Model
		SuppInfoKey       string              `json:"suppInfoKey"`
		SuppInfoVal       string              `json:"suppInfoVal"`
		SuppInfoKeyNameEn *string             `json:"suppInfoKeyNameEn"`
		SuppInfoKeyNameZh *string             `json:"suppInfoKeyNameZh"`
		SuppInfoDataType  string              `json:"suppInfoDataType"`
		ChklItemId        intstring.IntString `json:"chklItemId"`
		ChecklistItem     *ChecklistItem      `json:"checklistItem,omitempty"`
	}

	LA struct {
		model.Model
		Uuid              string                        `json:"uuid"`
		AssetNo           string                        `json:"assetNo"`
		Status            string                        `json:"status"`
		Manufacturer      *string                       `json:"manufacturer"`
		YearOfManufacture *string                       `json:"yearOfManufacture"`
		CraneSerialNo     *string                       `json:"craneSerialNo"`
		ModelNo           *string                       `json:"modelNo"`
		PlantGroup        string                        `json:"plantGroup"`
		PlantType         string                        `json:"plantType"`
		SerialNo          *string                       `json:"serialNo"`
		CertNo            *string                       `json:"certNo"`
		PlantOwner        *string                       `json:"plantOwner"`
		IsRental          *bool                         `json:"isRental"`
		MaxSafeLiftCap    *string                       `json:"maxSafeLiftCap"`
		Description       *string                       `json:"description"`
		ContractRefId     *intstring.IntString          `json:"contractId"`
		CertExamDate      *string                       `json:"certExamDate"`
		CertValidFrom     *string                       `json:"certValidFrom"`
		CertValidTo       *string                       `json:"certValidTo"`
		Logs              []LALog                       `json:"logs"`
		Permits           []MasterPermit                `json:"permits"`
		Media             map[string][]model.MediaParam `json:"media"`
	}
	LALog struct {
		Id             intstring.IntString `json:"id"`
		CreatedAt      time.Time           `json:"createdAt"`
		CreatedBy      string              `json:"createdBy"`
		Actor          *model.UserInfo     `json:"actor"`
		Message        *string             `json:"message"`
		MessageZh      *string             `json:"messageZh"`
		ActivityType   string              `json:"activityType"`
		PlantEquipLaId intstring.IntString `json:"plantEquipLaId"`
	}
	GetAllPermitOps struct {
		GetApprovalStage bool `json:"getApprovalStage"`
		// Additional filtering options that have less general uses
		// Filter by multiple contract ids
		ContractRefIds []intstring.IntString `json:"contractIds"`
		// Filter by flow isCurrent and action_type
		CurrentFlowActionTypes []string `json:"currentFlowActionType"`
	}
	PermitCriteria struct {
		MasterPermit
		ParticipantUserRefKeys []string            `json:"participantUserRefKeys"`
		SearchType             string              `json:"searchType"`
		PermitStatuses         []string            `json:"permitStatuses"`
		ContractRefId          intstring.IntString `json:"contractId"`
	}
	PermitApptCriteria struct {
		ParticipantUserRefKeys []string `json:"participantUserRefKeys"`
		SearchType             string   `json:"searchType"`
		ApptStatuses           []string `json:"apptStatuses"`
	}
	ApplicantDisplay struct {
		DisplayName  string  `json:"displayName"`
		Position     *string `json:"position"`
		ContactNo    *string `json:"contactNo"`
		PartyName    string  `json:"partyName"`
		SignatureB64 *string `json:"signature"`
	}
	ConstructionWork struct {
		model.Model
		ConstWork      *string              `json:"constWork"`
		WorkLocation   *string              `json:"workLocation"`
		NoiseControlId *intstring.IntString `json:"noiseControlId"`
	}
	MechEquipment struct {
		model.Model
		CnpRefNo       *string              `json:"cnpRefNo"`
		EquipGroupName *string              `json:"equipGroupName"`
		IdCode         *string              `json:"idCode"`
		Pme            *string              `json:"pme"`
		EquipQty       *int                 `json:"equipQty"`
		WorkLocation   *string              `json:"workLocation"`
		NoiseControlId *intstring.IntString `json:"noiseControlId"`
	}
	PrescribedWork struct {
		model.Model
		Seq            int                  `json:"seq"`
		IsSelected     *bool                `json:"isSelected"`
		IdCode         string               `json:"idCode"`
		DescriptionEn  *string              `json:"descriptionEn"`
		DescriptionZh  *string              `json:"descriptionZh"`
		NoiseControlId *intstring.IntString `json:"noiseControlId"`
	}
	ProjectPermitInfo struct {
		model.Model
		ProjectName    *string              `json:"projectName"`
		PermitRefId    *intstring.IntString `json:"permitRefId"`
		PermitRef      *Reference           `json:"permitRef"`
		NoiseControlId *intstring.IntString `json:"noiseControlId"`
	}
	Worker struct {
		model.Model
		WorkerName     *string              `json:"workerName"`
		GreenCardNo    *string              `json:"greenCardNo"`
		NoiseControlId *intstring.IntString `json:"noiseControlId"`
	}
	Reference struct {
		model.Model
		RefNo             string              `json:"refNo"`
		RefType           string              `json:"refType"`
		FullAddress       string              `json:"fullAddress"`
		Status            string              `json:"status"`
		DocumentNo        string              `json:"documentNo"`
		DocumentName      *string             `json:"documentName"`
		DocumentUrl       *string             `json:"documentUrl"`
		DocumentExtension *string             `json:"documentExtension"`
		DocumentMimeType  *string             `json:"documentMimeType"`
		ContractRefId     intstring.IntString `json:"contractRefId"`
		ValidFrom         time.Time           `json:"validFrom"`
		ValidTo           time.Time           `json:"validTo"`
		ConstPeriods      []ConstPeriod       `json:"constPeriods"`
		PrescribedWorks   []PrescribedWork    `json:"prescribedWorks"`
	}
	ConstPeriod struct {
		model.Model
		PeriodType   string              `json:"periodType"`
		DurationFrom string              `json:"durationFrom"`
		DurationTo   string              `json:"durationTo"`
		PermitRefId  intstring.IntString `json:"permitRefId"`
	}
	LiftingWorker struct {
		model.Model
		WorkerType           string              `json:"workerType"`
		WorkerName           *string             `json:"workerName"`
		WorkerCertNo         *string             `json:"workerCertNo"`
		WorkerCertExpiryDate *string             `json:"workerCertExpiryDate"`
		WorkerRiggerPhotoUrl *string             `json:"workerRiggerPhotoUrl"`
		PermitLiftId         intstring.IntString `json:"permitLiftId"`
	}

	EFLiftingWorker struct {
		model.Model
		WorkerType           string              `json:"workerType"`
		WorkerName           *string             `json:"workerName"`
		WorkerSitePassNo     *string             `json:"workerSitePassNo"`
		WorkerCertNo         *string             `json:"workerCertNo"`
		WorkerCertExpiryDate *string             `json:"workerCertExpiryDate"`
		WorkerSignature      *string             `json:"workerSignature"`
		PermitLiftId         intstring.IntString `json:"permitLiftId"`
	}

	LiftingGear struct {
		model.Model
		LgType         *string             `json:"lgType"`
		LgTypeRemark   *string             `json:"lgTypeRemark"` // TODO to be removed
		LgMark         *string             `json:"lgMark"`
		OwnerId        *string             `json:"ownerId"`
		CertExpiryDate *string             `json:"certExpiryDate"`
		PermitLiftId   intstring.IntString `json:"permitLiftId"`
	}

	EFLiftingGear struct {
		model.Model
		LgNo            *string             `json:"lgNo"`
		From6           *string             `json:"formSix"`
		From7           *string             `json:"formSeven"`
		CertExpiryDate  *string             `json:"certExpiryDate"`
		SafeWorkingLoad *string             `json:"safeWorkingLoad"`
		PermitLiftId    intstring.IntString `json:"permitLiftId"`
	}

	CSWorker struct {
		model.Model
		WorkerType           string              `json:"workerType"`
		WorkerName           *string             `json:"workerName"`
		WorkerCompanyName    *string             `json:"workerCompanyName"`
		WorkerCertNo         *string             `json:"workerCertNo"`
		WorkerCertExpiryDate *string             `json:"workerCertExpiryDate"`
		ConfinedSpaceId      intstring.IntString `json:"confinedSpaceId"`
	}

	LSWorker struct {
		model.Model
		WorkerType        string              `json:"workerType"`
		WorkerName        *string             `json:"workerName"`
		WorkerCompanyName *string             `json:"workerCompanyName"`
		LiftShaftId       intstring.IntString `json:"liftShaftId"`
	}

	EigiblePersons struct {
		model.Model
		SubmitterType       string              `json:"submitterType"`
		SubmitterName       string              `json:"SubmitterName"`
		SubmittedAt         string              `json:"submittedAt"`
		SubmitterCertNo     string              `json:"submitterCertNo"`
		SubmitterCertExpiry *string             `json:"submitterCertExpiry"`
		SubmitterPosition   *string             `json:"submitterPosition"`
		SubmitterSignature  *string             `json:"submitterSignature"`
		HasVerifiedAbove    bool                `json:"hasVerifiedAbove"`
		Seq                 intstring.IntString `json:"seq"`
		ConfinedSpaceId     intstring.IntString `json:"confinedSpaceId"`
	}

	DetectiveReports struct {
		model.Model
		LocationDepth     string              `json:"locationDepth"`
		ReportDate        string              `json:"reportDate"`
		OxygenVal         string              `json:"oxygenVal"`
		OxygenUuid        string              `json:"oxygenUuid"`
		FlammableGasVal   string              `json:"flammableGasVal"`
		FlammableGasUuid  string              `json:"flammableGasUuid"`
		H2sVal            string              `json:"h2sVal"`
		H2sUuid           string              `json:"h2sUuid"`
		CoVal             string              `json:"coVal"`
		CoUuid            string              `json:"coUuid"`
		OxygenMedia       []model.MediaParam  `json:"oxygenMedia"`
		FlammableGasMedia []model.MediaParam  `json:"flammableGasMedia"`
		H2sMedia          []model.MediaParam  `json:"h2sMedia"`
		CoMedia           []model.MediaParam  `json:"coMedia"`
		ConfinedSpaceId   intstring.IntString `json:"confinedSpaceId"`
		ReportTime        *string             `json:"reportTime"`
	}

	EmergencyContacts struct {
		model.Model
		ContactType     string              `json:"contactType"`
		ContactName     string              `json:"contactName"`
		ContactNo       string              `json:"contactNo"`
		Seq             string              `json:"seq"`
		ConfinedSpaceId intstring.IntString `json:"confinedSpaceId"`
	}

	DSDWorker struct {
		model.Model
		WorkerType                 string              `json:"workerType"`
		WorkerName                 *string             `json:"workerName"`
		WorkerCompanyName          *string             `json:"workerCompanyName"`
		WorkerCertNo               *string             `json:"workerCertNo"`
		WorkerSitePassNo           *string             `json:"workerSitePassNo"`
		WorkerSignature            *string             `json:"workerSignature"`
		ConfinedSpaceId            intstring.IntString `json:"confinedSpaceId"`
		WorkerCertExpiryDate       *string             `json:"workerCertExpiryDate"`
		WorkerSitePassNoExpiryDate *string             `json:"workerSitePassNoExpiryDate"`
	}

	IcItems struct {
		model.Model
		InspectorName      string              `json:"inspectorName"`
		InspectorPosition  string              `json:"inspectorPosition"`
		InspectionDate     string              `json:"inspectionDate"`
		InspectorSignature *string             `json:"inspectorSignature"`
		ConfinedSpaceId    intstring.IntString `json:"confinedSpaceId"`
		Media              []model.MediaParam  `json:"media"`
	}

	Inspectors struct {
		model.Model
		InspectorName   string              `json:"inspectorName"`
		CsCertNo        string              `json:"csCertNo"`
		ApprWorkerNo    string              `json:"apprWorkerNo"`
		CompanyName     string              `json:"companyName"`
		Remark          string              `json:"remark"`
		Signature       *string             `json:"signature"`
		ConfinedSpaceId intstring.IntString `json:"confinedSpaceId"`
		ExpiryDate      *string             `json:"expiryDate"`
	}
)
