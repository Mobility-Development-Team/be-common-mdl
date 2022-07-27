package wfutil

import (
	"regexp"
)

type UpdateStatusParams map[string]*string

const StatusParent = "PARENT"

var regexParseAction = regexp.MustCompile(`(.*?)_?STATUS_CHANGE_(.*)`)

// NewDefaultUpdateStatusParams Register a UpdateStatusParams to parse multiple STATUS_CHANGE actions
// The following format is supported:
//	- {RULE}_STATUS_CHANGE_{TARGET_STATUS}
//	- STATUS_CHANGE_{TARGET_STATUS}
// If RULE is not specified, the pointer to defaultUpdates will be updated
// with TARGET_STATUS. If RULE is specified and a matching rule is added
// via AddRule() that matching content of the pointer will be updated.
func NewDefaultUpdateStatusParams(defaultUpdates *string) UpdateStatusParams {
	return map[string]*string{
		"": defaultUpdates,
	}
}

func (u UpdateStatusParams) AddRule(rule string, updates *string) UpdateStatusParams {
	u[rule] = updates
	return u
}

func (u UpdateStatusParams) UpdateStatus(status string) {
	matches := regexParseAction.FindStringSubmatch(status)
	if matches == nil {
		return
	}
	rule, value := matches[1], matches[2]
	if toUpdate, ok := u[rule]; ok && toUpdate != nil {
		*toUpdate = value
	}
}
