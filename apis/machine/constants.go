package machine

const (
	PreloadApprovalFlows = "APPROVAL_FLOWS"
	PreloadChecklists    = "CHECKLISTS"
	PreloadLogs          = "LOGS"
	PreloadParticipants  = "PARTICIPANTS"
	PreloadAttachments   = "ATTACHMENTS"

	PermitSearchTypeCreated    = "CREATED"
	PermitSearchTypeAssigned   = "ASSIGNED"
	PermitSearchActionRequired = "ACTION_REQUIRED"
	PermitSearchTypeAll        = "ALL"
	PermitSearchTypeSystem     = "SYSTEM"

	PermitStatusDraft                = "DRAFT"
	PermitStatusWorkInProgress       = "WORK_IN_PROGRESS"
	PermitStatusAwaitingApproval     = "AWAITING_APPROVAL"
	PermitStatusApproved             = "APPROVED"
	PermitStatusExpired              = "EXPIRED"
	PermitStatusAwaitingAccepted     = "ACCEPTED"
	PermitStatusValid                = "VALID"
	PermitStatusAwaitingAcceptance   = "AWAITING_ACCEPTANCE"
	PermitStatusAwaitingCancellation = "AWAITING_CANCELLATION"
)
