package inspection

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/Mobility-Development-Team/be-common-mdl/util/strutil"
)

type (
	Appointment struct {
		model.Model
		AppointmentDateComponent *time.Time          `json:"-" gorm:"column:appointment_date"`
		AppointmentStatus        string              `json:"appointmentStatus"`
		SiteWalkId               intstring.IntString `json:"siteWalkId"`
		TimeSlots                []TimeSlot          `json:"timeSlots"`
		Responses                []Response          `json:"responses"`
		Invitees                 []Invitee           `json:"invitees"`
		SiteWalk                 *SiteWalk           `json:"siteWalk,omitempty" gorm:"->;foreignKey:SiteWalkId"`
	}
	TimeSlot struct {
		model.Model
		Status                           string              `json:"status"`
		Accepted                         bool                `json:"accepted"`
		ProposedAppointmentDate          string              `json:"proposedAppointmentDate" gorm:"-"`
		ProposedStartTime                string              `json:"proposedStartTime" gorm:"-"`
		ProposedEndTime                  string              `json:"proposedEndTime" gorm:"-"`
		ProposedAppointmentDateComponent *time.Time          `json:"-" gorm:"column:proposed_appointment_date"`
		ProposedStartTimeComponent       *time.Time          `json:"-" gorm:"column:proposed_start_time"`
		ProposedEndTimeComponent         *time.Time          `json:"-" gorm:"column:proposed_end_time"`
		ProposedPeriod                   string              `json:"proposedPeriod"`
		ProposedGroup                    string              `json:"proposedGroup"`
		IsPrimary                        bool                `json:"isPrimary"`
		Seq                              int                 `json:"seq"`
		AppointmentId                    intstring.IntString `json:"appointmentId"`
	}
	Response struct {
		model.Model
		Response      string              `json:"response"`
		Remark        *string             `json:"remark"`
		AppointmentId intstring.IntString `json:"appointmentId"`
	}
	Invitee struct {
		model.Model
		InviteeUserId          intstring.IntString  `json:"inviteeUserId"`
		InviteeUserRefKey      string               `json:"inviteeUserRefKey"`
		InviteeLastName        *string              `gorm:"-" json:"inviteeLastName"`
		InviteeFirstName       *string              `gorm:"-" json:"inviteeFirstName"`
		InviteeDisplayName     *string              `gorm:"-" json:"inviteeDisplayName"`
		InviteeProfileIconUrl  *string              `gorm:"-" json:"inviteeProfileIconUrl"`
		InviteeEmail           *string              `gorm:"-" json:"inviteeEmail"`
		InviteePosition        *string              `gorm:"-" json:"inviteePosition"`
		InviteeParticipantType *string              `json:"inviteeParticipantType"`
		InviteeAttendeeType    *string              `json:"inviteeAttendeeType"`
		AppointmentId          *intstring.IntString `json:"appointmentId"`
	}

	// Some fields are commented out due to not being used
	// uncomment the fields and add the corresponding structure if necessary
	TaskDisplay struct {
		model.Model
		Title               string                       `json:"title"`
		Purpose             string                       `json:"purpose"`
		TaskType            string                       `json:"taskType"`
		TaskStatus          string                       `json:"taskStatus"`
		ContentOfTask       string                       `json:"contentOfTask"`
		SuggestedFollowUp   string                       `json:"suggestedFollowUp"`
		AnticipatedCompDate string                       `json:"anticipatedCompDate"`
		TaskDueDate         string                       `json:"taskDueDate"`
		DaysToDue           int                          `json:"daysToDue"`
		HasOverdue          bool                         `json:"hasOverdue"`
		IsManualCompletion  bool                         `json:"isManualCompletion"`
		DueType             string                       `json:"dueType"`
		FollowUpCompDate    string                       `json:"followUpCompDate"`
		MediaBatchRefUuid   string                       `json:"mediaBatchRefUuid"`
		SiteWalkId          *intstring.IntString         `json:"siteWalkId"`
		ContractRefId       *intstring.IntString         `json:"contractId"`
		Media               []model.MediaParam           `json:"media" gorm:"-"` // Media accepts both array of string or struct during json unmarshal
		FollowUpUsers       []FindingFollowUpUserDisplay `json:"followUpUsers"`
		// Hashtags            HashtagDisplay               `json:"hashtags"`
	}
)

type (
	SiteWalkParticipant struct {
		model.Model
		ParticipantUserId         intstring.IntString  `json:"participantUserId"`
		ParticipantUserUuid       *string              `json:"participantUserUuid"`
		ParticipantLastName       *string              `json:"participantLastName"`
		ParticipantFirstName      *string              `json:"participantFirstName"`
		ParticipantDisplayName    *string              `json:"participantDisplayName"`
		ParticipantProfileIconUrl *string              `json:"participantProfileIconUrl"`
		ParticipantEmail          *string              `json:"participantEmail"`
		ParticipantUserRefKey     string               `json:"participantUserRefKey"`
		ParticipantType           *string              `json:"participantType"`
		SiteWalkId                *intstring.IntString `json:"siteWalkId"`
		ChecklistId               *intstring.IntString `json:"checklistId"`

		// To be filled by SiteWalk.GetDetailedParticipants, uncomment to use
		// ParticipantUser model.UserInfo `json:"participantUser"`
	}
	TemplateItemSimple struct {
		Id         intstring.IntString `json:"id"`
		ItemNo     string              `json:"itemNo"`
		ItemNameEn string              `json:"itemNameEn"`
		ItemNameZh string              `json:"itemNameZh"`
		Seq        int                 `json:"seq"`
	}
	TemplateSimple struct {
		Id               intstring.IntString `json:"id"`
		InspectionNameEn string              `json:"checklistTitleEn"`
		InspectionNameZh string              `json:"checklistTitleZh"`
		ParentTemplate   *TemplateSimple     `json:"parentTemplate,omitempty"`
		SubTemplate      *TemplateSimple     `json:"subTemplate,omitempty"`
	}
	SiteWalkLocation struct {
		model.Model
		LocationRefId *intstring.IntString `json:"locationId"`
		SiteWalkId    *intstring.IntString `json:"siteWalkId"`
		ChecklistId   *intstring.IntString `json:"checklistId"`
	}
	ChecklistItem struct {
		model.Model
		Uuid                      string               `json:"uuid"`
		TemplateItemRefId         *intstring.IntString `json:"templateItemId"`
		TemplateItem              *TemplateItemSimple  `json:"templateItem" gorm:"foreignKey:TemplateItemRefId"`
		Response                  *string              `json:"response"`
		ResponsedByUserRefKey     *string              `json:"responsedByUserRefKey"`
		ResponsedByDisplay        *model.UserInfo      `json:"responsedBy" `
		Remark                    *string              `json:"remark"`
		ChecklistId               *intstring.IntString `json:"checklistId"`
		ChecklistItemMediaDisplay []model.MediaParam   `json:"checklistItemMedia"`
	}
	Checklist struct {
		model.Model
		Uuid                     string                 `json:"uuid"`
		InspectionDate           string                 `json:"inspectionDate"`
		InspectionStartTime      string                 `json:"inspectionStartTime"`
		InspectionEndTime        string                 `json:"inspectionEndTime"`
		InspectionDueDate        *time.Time             `json:"inspectionDueDate"`
		InspectionCompletionDate *time.Time             `json:"inspectionCompletionDate"`
		FollowUp                 bool                   `json:"followUp"`
		ReviewedAt               *time.Time             `json:"reviewedAt"`
		ReviewedBy               *intstring.IntString   `json:"reviewedBy"`
		TemplateRefId            *intstring.IntString   `json:"templateId"`
		Template                 TemplateSimple         `json:"template"`
		SiteWalkId               intstring.IntString    `json:"siteWalkId"`
		Participants             []*SiteWalkParticipant `json:"participants,omitempty"`
		ChecklistItems           []ChecklistItem        `json:"checklistItems,omitempty"`
		LocationsDisplay         []*model.Location      `json:"locations"`
		LocationIds              []intstring.IntString  `json:"locationIds"`
		TotalItemSize            int                    `json:"totalItemSize"`
		OkItemSize               int                    `json:"okItemSize"`
		FailedItemSize           int                    `json:"failedItemSize"`
		NaItemSize               int                    `json:"naItemSize"`
	}
	SiteWalk struct {
		model.Model
		Uuid                string                  `json:"uuid"`
		InspectionTaskType  string                  `json:"inspectionTaskType"`
		InspectionStartDate string                  `json:"inspectionStartDate"`
		InspectionStartTime string                  `json:"inspectionStartTime"`
		InspectionEndTime   string                  `json:"inspectionEndTime"`
		InspectionStatus    string                  `json:"inspectionStatus"`
		InspectionName      *string                 `json:"inspectionName"`
		ContractRefId       *intstring.IntString    `json:"contractId"`
		PartyRefId          *intstring.IntString    `json:"partyId"`
		MediaBatchRefUuid   string                  `json:"mediaBatchUuid"`
		Checklists          []Checklist             `json:"checklists,omitempty"`
		TotalItemSize       int                     `json:"totalItemSize"`
		OkItemSize          int                     `json:"okItemSize"`
		FailedItemSize      int                     `json:"failedItemSize"`
		NaItemSize          int                     `json:"naItemSize"`
		Participants        []*SiteWalkParticipant  `json:"participants"`
		LocationsDisplay    []*model.Location       `json:"locations"`
		LocationIds         []intstring.IntString   `json:"locationIds"`
		GeneralFindings     []GeneralFindingDisplay `json:"generalFindings"`
		NcFindings          []NcFindingDisplay      `json:"ncFindings"`
		RatChecklist        *RatChecklistDisplay    `json:"ratChecklist"`
		Contract            ContractInfoDisplay     `json:"contract"`
		SiteWalkTypes       []SiteWalkType          `json:"siteWalkTypes"`
	}
	ContractInfoDisplay struct {
		ContractNo   string `json:"contractNo"`
		ContractDesc string `json:"contractDesc"`
	}
)

type (
	RatChecklist struct {
		model.Model
		Uuid       string              `json:"uuid"`
		SiteWalkId intstring.IntString `json:"siteWalkId"`
	}
	RatChecklistDisplay struct {
		RatChecklist
		PartOne        interface{}                        `json:"partOne"`
		PartTwoItems   []RatChecklistPartTwoItemDisplay   `json:"partTwoItems"`
		PartThreeItems []RatChecklistPartOtherItemDisplay `json:"partThreeItems"`
		PartFourItems  []RatChecklistPartOtherItemDisplay `json:"partFourItems"`
	}

	RatChecklistItem struct {
		model.Model
		Uuid               string              `json:"uuid"`
		PartNo             intstring.IntString `json:"partNo"`
		Description        *string             `json:"description"`
		IsAvailable        *bool               `json:"isAvailable"`
		ApprovalStatus     *bool               `json:"approvalStatus"`
		IsShownOn          *bool               `json:"isShownOn"`
		IsCompApproved     *bool               `json:"isCompApproved"`
		ApprovalFromEngine *bool               `json:"approvalFromEngine"`
		ApprovalFromIce    *bool               `json:"approvalFromIce"`
		ApprovalFromTechD  *bool               `json:"approvalFromTechD"`
		Remark             *string             `json:"remark"`
		RemedialAction     *string             `json:"remedialAction"`
		RatChecklistId     intstring.IntString `json:"ratChecklistId"`
	}
	RatChecklistPartTwoItemDisplay struct {
		model.Model
		Uuid           string              `json:"uuid"`
		PartNo         intstring.IntString `json:"partNo"`
		Description    *string             `json:"description"`
		IsAvailable    *bool               `json:"isAvailable"`
		ApprovalStatus *bool               `json:"approvalStatus"`
		Remark         *string             `json:"remark"`
		RemedialAction *string             `json:"remedialAction"`
		Name           string              `json:"name"`
		// RatChecklistId intstring.IntString `json:"ratChecklistId"`
	}
	RatChecklistPartOtherItemDisplay struct {
		model.Model
		Uuid                string              `json:"uuid"`
		PartNo              intstring.IntString `json:"partNo"`
		Description         *string             `json:"description"`
		IsShownOn           *bool               `json:"isShownOn"`
		ApprovalFromEngine  *bool               `json:"approvalFromEngine"`
		ApprovalFromIce     *bool               `json:"approvalFromIce"`
		ApprovalFromTechD   *bool               `json:"approvalFromTechD"`
		IsCompApproved      *bool               `json:"isCompApproved"`
		Remark              *string             `json:"remark"`
		RemedialAction      *string             `json:"remedialAction"`
		AnticipatedCompDate string              `json:"anticipatedCompDate"`
		// RatChecklistId     intstring.IntString `json:"ratChecklistId"`
	}
)

type (
	GeneralFindingDisplay struct {
		model.Model
		Title      string              `json:"title" gorm:"not null"`
		Remark     *string             `json:"remark"`
		SiteWalkId intstring.IntString `json:"siteWalkId"`
		Media      []model.MediaParam  `json:"media" gorm:"-"`
	}
	NcFindingDisplay struct {
		model.Model
		Title               string                       `json:"title"`
		Purpose             string                       `json:"purpose"`
		TaskType            string                       `json:"taskType"` // NC_FINDING, MEDIA, CHECKLIST
		TaskStatus          string                       `json:"taskStatus"`
		ContentOfTask       string                       `json:"contentOfTask"`
		SuggestedFollowUp   string                       `json:"suggestedFollowUp"`
		AnticipatedCompDate string                       `json:"anticipatedCompDate"`
		TaskDueDate         string                       `json:"taskDueDate"`
		DaysToDue           int                          `json:"daysToDue"`
		DueType             string                       `json:"dueType"`
		FollowUpCompDate    string                       `json:"followUpCompDate"`
		FollowUpUsers       []FindingFollowUpUserDisplay `json:"followUpUsers"`
		MediaBatchRefUuid   string                       `json:"mediaBatchRefUuid"`
		SiteWalkId          *intstring.IntString         `json:"siteWalkId"`
		ContractRefId       *intstring.IntString         `json:"contractId"`
		Media               []model.MediaParam           `json:"media" gorm:"-"`
		// Hashtags            HashtagDisplay               `json:"hashtags"`
	}
	FindingFollowUpUserDisplay struct {
		model.Model
		Status                     string              `json:"status"`
		FollowUpUserId             intstring.IntString `json:"followUpUserId"`
		FollowUpUserRefKey         string              `json:"followUpUserRefKey"`
		FollowUpUserLastName       *string             `json:"followUpLastName"`
		FollowUpUserFirstName      *string             `json:"followUpFirstName"`
		FollowUpUserDisplayName    *string             `json:"followUpDisplayName"`
		FollowUpUserProfileIconUrl *string             `json:"followUpProfileIconUrl"`
		FollowUpUserEmail          *string             `json:"followUpEmail"`
		FollowUpUserPosition       *string             `json:"followUpUserPosition"`
	}
)

func (ud FindingFollowUpUserDisplay) GetAvailableName() string {
	return model.GetAvailableName(
		strutil.StrOrEmpty(ud.FollowUpUserDisplayName),
		strutil.StrOrEmpty(ud.FollowUpUserFirstName),
		strutil.StrOrEmpty(ud.FollowUpUserLastName),
		strutil.StrOrEmpty(ud.FollowUpUserEmail),
	)
}

type (
	Attachment struct {
		AttachmentName      string              `json:"attachmentName"`
		AttachmentUrl       string              `json:"attachmentUrl"`
		AttachmentType      string              `json:"attachmentType"`
		AttachmentExtension string              `json:"attachmentExtension"`
		AttachmentMineType  string              `json:"attachmentMineType"`
		SiteWalkId          intstring.IntString `json:"siteWalkId"`
	}
	ActivityLog struct {
		Id           intstring.IntString  `json:"id"`
		CreatedAt    time.Time            `json:"createdAt"`
		CreatedBy    string               `json:"createdBy"`
		Message      *string              `json:"message"`
		MessageZh    *string              `json:"messageZh"`
		SiteWalkId   *intstring.IntString `json:"siteWalkId"`
		ChecklistId  *intstring.IntString `json:"checklistId"`
		ActivityType string               `json:"activityType"`
	}
)

// Some API calls returns SiteWalkType as a string. Some fields might not be available
type SiteWalkType struct {
	model.Model
	ContractRefId intstring.IntString `json:"contractRefId"`
	SiteWalkType  string              `json:"siteWalkType"`
	SiteWalkId    intstring.IntString `json:"siteWalkId"`
}

func (st *SiteWalkType) UnmarshalJSON(b []byte) error {
	type alias SiteWalkType
	var resultStruct alias
	if err := json.Unmarshal(b, &resultStruct); err == nil {
		// Unmarshal successful
		*st = SiteWalkType(resultStruct)
		return nil
	}
	var resultString string
	if err := json.Unmarshal(b, &resultString); err != nil {
		return errors.New("not a struct nor a string")
	}
	*st = SiteWalkType{
		SiteWalkType: resultString,
	}
	return nil
}

type (
	SitePlanDisplay struct {
		model.Model
		ImagePreview     string                `json:"imagePreview"`
		PrevImagePreview string                `json:"prevImagePreview"`
		OriginalImageUrl string                `json:"sitePlanPictureUrl"`
		SiteWalkId       intstring.IntString   `json:"siteWalkId"`
		General          SitePlanDetailDisplay `json:"general"`
		NonCompliance    SitePlanDetailDisplay `json:"nc"`
	}

	SitePlanDetailDisplay struct {
		Arrows []ArrowDisplay `json:"arrows"`
		Pins   []PinDisplay   `json:"points"`
	}

	ArrowDisplay struct {
		Id            intstring.IntString  `json:"id,omitempty"`
		RefUuid       string               `json:"_id"`
		AxisX         float64              `json:"x"`
		AxisY         float64              `json:"y"`
		IsDisplay     *bool                `json:"isDisplay"`
		Rotation      float64              `json:"rotation"`
		Fill          string               `json:"fill"`
		Stroke        string               `json:"stroke"`
		Points        []int                `json:"points"`
		PointerLength float64              `json:"pointerLength"`
		PointerWidth  float64              `json:"pointerWidth"`
		StrokeWidth   int                  `json:"strokeWidth"`
		Name          string               `json:"name"`
		Draggable     bool                 `json:"draggable"`
		FindingRefId  *intstring.IntString `json:"findingRefId"`
		SitePlanId    intstring.IntString  `json:"sitePlanId"`
	}

	PinDisplay struct {
		Id                  intstring.IntString  `json:"id,omitempty"`
		AxisX               float64              `json:"x"`
		AxisY               float64              `json:"y"`
		PinKey              string               `json:"key"`
		IsDisplay           *bool                `json:"isDisplay"`
		RefUuid             string               `json:"_id"`
		Fill                string               `json:"fill"`
		Stroke              string               `json:"stroke"`
		Width               float64              `json:"width"`
		Height              float64              `json:"height"`
		Title               string               `json:"title"`
		Media               []model.MediaParam   `json:"media"`
		Content             string               `json:"content"`
		Suggested           string               `json:"suggested"`
		Purpose             string               `json:"purpose"`
		AnticipatedCompDate string               `json:"anticipatedCompletionDate"`
		Remark              string               `json:"remark"`
		Status              string               `json:"status"`
		FindingType         string               `json:"findingType"`
		FindingRefId        *intstring.IntString `json:"findingRefId"`
	}
)

type (
	FollowUpTaskDisplay struct {
		FollowUpTaskBase
		TaskActions      []FollowUpTaskActionDisplay  `json:"actions"`
		TaskComments     []FollowUpTaskCommentDisplay `json:"comments"`
		TaskActivityLogs []TaskActivityLogDisplay     `json:"logs"`
		AllowNewAction   bool                         `json:"allowNewAction"`
	}
	FollowUpTaskBase struct {
		model.Model
		TaskParentType              *string              `json:"taskParentType"`
		TaskParentRefId             *intstring.IntString `json:"taskParentRefId"`
		TaskParentMediaBatchRefUuid string               `json:"taskParentMediaBatchRefUuid"`
	}
	FollowUpTaskActionDisplay struct {
		FollowUpTaskActionBase
		Seq                 int             `json:"seq"`
		AnticipatedCompDate string          `json:"anticipatedCompDate"`
		SubmittedByUser     *model.UserInfo `json:"submittedByUser"`
		// Media               []model.MediaParam `json:"media"`
		FollowUpUserMedia []model.MediaParam `json:"followUpUserMedia"`
		ApproverMedia     []model.MediaParam `json:"approverMedia"`
	}
	FollowUpTaskActionBase struct {
		model.Model
		ActionStatus          string              `json:"actionStatus"`
		SubmittedByUserId     intstring.IntString `json:"submittedByUserId"`
		SubmittedByUserRefKey string              `json:"submittedByUserRefKey"`
		TaskDueDate           *time.Time          `json:"taskDueDate"`
		IsOverdue             bool                `json:"isOverdue"`
		FollowUpRemark        string              `json:"followUpRemark"`
		ApprovalRemark        *string             `json:"approvalRemark"`
		FollowUpTaskId        intstring.IntString `json:"followUpTaskId"`
		WorkflowRefUuid       *string             `json:"workflowRefUuid"`
	}
	FollowUpTaskComment struct {
		model.Model
		CommentByUserId     intstring.IntString `json:"commentByUserId"`
		CommentByUserRefKey string              `json:"commentByUserRefKey"`
		CommentMessage      string              `json:"commentMessage"`
		FollowUpTaskId      intstring.IntString `json:"followUpTaskId"`
	}
	FollowUpTaskCommentDisplay struct {
		FollowUpTaskComment
		CommentByUser *model.UserInfo `json:"commentByUser"`
	}
	TaskActivityLogDisplay struct {
		Id              intstring.IntString  `gorm:"primaryKey" json:"id,omitempty"`
		CreatedAt       time.Time            `json:"createdAt"`
		CreatedBy       string               `json:"createdBy" gorm:"column:created_by"`
		ActorUserId     *intstring.IntString `json:"-"`
		ActorUserRefKey *string              `json:"-"`
		Actor           *model.UserInfo      `json:"actor"`
		Message         *string              `json:"message"`
		MessageZh       *string              `json:"messageZh"`
		ActivityType    string               `json:"activityType"`
		TaskId          intstring.IntString  `json:"taskId"`
	}
)
