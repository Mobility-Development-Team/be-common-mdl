package models

import (
	"time"

	"github.com/Mobility-Development-Team/be-common-mdl/model"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
)

type Appointment struct {
	model.Model
	AppointmentDateComponent *time.Time          `json:"-" gorm:"column:appointment_date"`
	AppointmentStatus        string              `json:"appointmentStatus"`
	SiteWalkId               intstring.IntString `json:"siteWalkId"`
	TimeSlots                []TimeSlot          `json:"timeSlots"`
	Responses                []Response          `json:"responses"`
	Invitees                 []Invitee           `json:"invitees"`
	SiteWalk                 *SiteWalk           `json:"siteWalk,omitempty" gorm:"->;foreignKey:SiteWalkId"`
}
type TimeSlot struct {
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
type Response struct {
	model.Model
	Response      string              `json:"response"`
	Remark        *string             `json:"remark"`
	AppointmentId intstring.IntString `json:"appointmentId"`
}
type Invitee struct {
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

// Note: This structure does not currently support preloading,
// uncomment the fields and add the corresponding structure if necessary
type SiteWalk struct {
	model.Model
	Uuid                         string                `json:"uuid"`
	InspectionTaskType           string                `json:"inspectionTaskType"`
	InspectionDateComponent      *time.Time            `json:"-" gorm:"column:inspection_date"`
	InspectionStartTimeComponent *time.Time            `json:"-" gorm:"column:inspection_start_time"`
	InspectionEndTimeComponent   *time.Time            `json:"-" gorm:"column:inspection_end_time"`
	InspectionStartDate          string                `json:"inspectionStartDate" gorm:"-"`
	InspectionStartTime          string                `json:"inspectionStartTime" gorm:"-"`
	InspectionEndTime            string                `json:"inspectionEndTime" gorm:"-"`
	InspectionStatus             string                `json:"inspectionStatus"`
	InspectionName               *string               `json:"inspectionName"`
	ContractRefId                *intstring.IntString  `json:"contractId"`
	PartyRefId                   *intstring.IntString  `json:"partyId"`
	MediaBatchRefUuid            string                `json:"mediaBatchUuid"`
	TotalItemSize                int                   `json:"totalItemSize" gorm:"-"`
	OkItemSize                   int                   `json:"okItemSize" gorm:"-"`
	FailedItemSize               int                   `json:"failedItemSize" gorm:"-"`
	NaItemSize                   int                   `json:"naItemSize" gorm:"-"`
	LocationsDisplay             []*model.Location     `json:"locations" gorm:"-"`
	LocationIds                  []intstring.IntString `json:"locationIds" gorm:"-"`
	CalendarDownloadUrl          string                `json:"calendarDownloadUrl" gorm:"-"`
	// Checklists                   []Checklist                       `json:"checklists"`
	// GeneralFindings              []general_finding.GeneralFinding  `json:"generalFindings"`
	// NcFindings                   []task.Task                       `json:"ncFindings" gorm:"foreignKey:SiteWalkId"`
	// RatChecklist                 *rat_checklist.RatChecklist       `json:"-" gorm:"foreignKey:SiteWalkId"`
	// Participants                 []participant.SiteWalkParticipant `json:"participants"`
	// Locations                    []location.SiteWalkLocation       `json:"-" gorm:"foreignKey:SiteWalkId"`
	// SiteWalkTypes                []SiteWalkType                    `json:"siteWalkTypes" gorm:"foreignKey:SiteWalkId"`
}

// Some fields are commented out due to not being used
// uncomment the fields and add the corresponding structure if necessary
type TaskDisplay struct {
	model.Model
	Title               string               `json:"title"`
	Purpose             string               `json:"purpose"`
	TaskType            string               `json:"taskType"`
	TaskStatus          string               `json:"taskStatus"`
	ContentOfTask       string               `json:"contentOfTask"`
	SuggestedFollowUp   string               `json:"suggestedFollowUp"`
	AnticipatedCompDate string               `json:"anticipatedCompDate"`
	TaskDueDate         string               `json:"taskDueDate"`
	DaysToDue           int                  `json:"daysToDue"`
	HasOverdue          bool                 `json:"hasOverdue"`
	IsManualCompletion  bool                 `json:"isManualCompletion"`
	DueType             string               `json:"dueType"`
	FollowUpCompDate    string               `json:"followUpCompDate"`
	MediaBatchRefUuid   string               `json:"mediaBatchRefUuid"`
	SiteWalkId          *intstring.IntString `json:"siteWalkId"`
	ContractRefId       *intstring.IntString `json:"contractId"`
	Media               []model.MediaParam   `json:"media" gorm:"-"` // Media accepts both array of string or struct during json unmarshal
	// FollowUpUsers       []FindingFollowUpUserDisplay `json:"followUpUsers"`
	// Hashtags            HashtagDisplay               `json:"hashtags"`
}
