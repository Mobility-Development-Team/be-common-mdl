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
		PlantType         string              `json:"plantType"`
		PlantOwner        string              `json:"plantOwner"`
		Manufacturer      string              `json:"manufacturer"`
		ModelNo           string              `json:"modelNo"`
		YearOfManufacture string              `json:"yearOfManufacture"`
		SerialNo          string              `json:"serialNo"`
		OwnerNo           string              `json:"ownerNo"`
		IsRental          bool                `json:"isRental"`
		RejectionReason   string              `json:"rejectionReason"`
		PermitMasterId    intstring.IntString `json:"permitMasterId"`
		CertExpiryDate    string              `json:"certExpiryDate"`
		// Custom fields
		CurrentApprovalStage int `json:"currentApprovalStage"`
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
	}
	Checklist struct {
		model.Model
		TemplateRefKey       string              `json:"templateRefKey"`
		ChecklistNameEn      string              `json:"checklistNameEn"`
		ChecklistNameZh      string              `json:"checklistNameZh"`
		IsCompleted          bool                `json:"isCompleted"`
		PermitMasterId       intstring.IntString `json:"permitMasterId"`
		TemplateRefOwnerType string              `json:"templateRefOwnerType"`
		Items                []ChecklistItem     `json:"items" gorm:"foreignKey:PermitChecklistId"`
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
		Seq                   int                 `json:"seq"`
		Response              string              `json:"response"`
		ItemNameEn            string              `json:"itemNameEn"`
		ItemNameZh            string              `json:"itemNameZh"`
		ResponsedBy           intstring.IntString `json:"responsedBy"`
		TemplateItemRefId     intstring.IntString `json:"templateItemRefId"`
		PermitChecklistId     intstring.IntString `json:"permitChecklistId"`
		ResponsedByUserRefKey string              `json:"responsedByUserRefKey"`
		Media                 []model.MediaParam  `json:"media" gorm:"-"`
	}

	Equipment struct {
		model.Model
		Uuid              string         `json:"uuid"`
		PlantType         string         `json:"plantType"`
		PlantOwner        string         `json:"plantOwner"`
		Manufacturer      string         `json:"manufacturer"`
		ModelNo           string         `json:"modelNo"`
		YearOfManufacture string         `json:"yearOfManufacture"`
		CertExpiryDate    *time.Time     `json:"certExpiryDate"`
		SerialNo          string         `json:"serialNo"`
		OwnerNo           string         `json:"ownerNo"`
		IsRental          *bool          `json:"isRental"`
		Permits           []MasterPermit `json:"permits"`
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
)
