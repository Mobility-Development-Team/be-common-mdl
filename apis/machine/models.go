package machine

import (
	"time"

	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
)

type (
	PlantPermit struct {
		model.Model
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
		// Shadowing fields
		CertExpiryDate string `json:"certExpiryDate"`
		// From master permit's fields
		PermitNo           string                `json:"permitNo"`
		PermitType         string                `json:"permitType"`
		PermitBatchRefUuid string                `json:"permitBatchRefUuid"`
		PermitStatus       string                `json:"permitStatus"`
		ContractRefId      intstring.IntString   `json:"contractId"`
		Checklists         []Checklist           `json:"checklists"`
		Participants       []Participant         `json:"participants"`
		ApprovalFlow       []PlantPermitApproval `json:"approvalFlow"`
		// Custom fields
		CurrentApprovalStage int `json:"currentApprovalStage"`
		// Additional preloading fields (placeholders)
		Logs        []interface{} `json:"logs"`
		Attachments []interface{} `json:"attachments"`
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

	Participant struct {
		model.UserInfo
		ParticipantType string               `json:"participantType"`
		PartyRefId      *intstring.IntString `json:"partyRefId"`
		UserSource      string               `json:"userSource"`
		Party           *model.PartyInfo     `json:"party"`
	}

	PlantPermitApproval struct {
		model.Model
		SubmittedBy           *model.UserInfo `json:"submittedBy"`
		SubmittedByActionType string          `json:"submittedByActionType"`
		IsCompleted           bool            `json:"isCompleted"`
		IsRejected            bool            `json:"isRejected"`
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
)
