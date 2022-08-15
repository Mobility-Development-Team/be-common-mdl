package inspection

const (
	FindingTypeGeneral = "GENERAL"
	FindingTypeNC      = "NC"

	TaskPurposeRecord   = "RECORD"
	TaskPurposeFollowUP = "FOLLOWUP"

	TaskSearchTypeCreated  = "CREATED"
	TaskSearchTypeAssigned = "ASSIGNED"
	TaskSearchTypeAll      = "ALL"

	TaskStatusDraft              = "DRAFT"
	TaskStatusWorkInProgress     = "WORK_IN_PROGRESS"
	TaskStatusFurtherFollowUp    = "FURTHER_FOLLOW_UP"
	TaskStatusInAwaitingApproval = "AWAITING_APPROVAL"
	TaskStatusCompleted          = "COMPLETED"

	TaskActionStatusWorkInProgress     = "WORK_IN_PROGRESS"
	TaskActionStatusInAwaitingApproval = "AWAITING_APPROVAL"
	TaskActionStatusAccepted           = "ACCEPTED"
	TaskActionStatusRejected           = "REJECTED"
)
