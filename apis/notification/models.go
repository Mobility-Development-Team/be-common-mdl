package notification

import (
	"fmt"
	"strings"

	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"github.com/Mobility-Development-Team/be-common-mdl/util/strutil"
	logger "github.com/sirupsen/logrus"
)

const (
	cloudMessageMagicDefault = "{{$DEFAULT}}"
)

type (
	NotificationParams []interface{}
	Notification       struct {
		NotificationType string               `json:"notificationType"`
		TemplateType     string               `json:"templateType"`
		ContractID       *intstring.IntString `json:"contractId,omitempty"`
		Recipients       Recipients           `json:"recipients"`
		Params           NotificationParams   `json:"params"`
		Actions          []Action             `json:"actions"`
		WithMail         *MailOptions         `json:"withMail"`
		WithPush         []*CloudMessage      `json:"withPush"`
		AutoPush         *AutoPushParams      `json:"autoPush"`
	}
	Recipients struct {
		Users  []intstring.IntString `json:"users"`
		Groups []intstring.IntString `json:"groups"`
	}
	Action struct {
		ActionID    string `json:"actionId"`
		ActionLabel string `json:"actionLabel"`
	}
	MailTemplate interface {
		Id() string
	}
	Mail struct {
		Recipients []string            `json:"recipients"`
		UserId     intstring.IntString `json:"userId"`
		Body       MailTemplate        `json:"body"`
		// The following values must be a struct
		// These language specific body payload will be merged to the base `body`
		// payload according to the user selected language or the default language
		BodyEn interface{} `json:"bodyEn"`
		BodyZh interface{} `json:"bodyZh"`
	}
	MailOptions struct {
		TemplateId string `json:"templateId"`
		Mails      []Mail `json:"mails"`
	}
	AutoPushParams struct {
		NoSelf  bool    `json:"noSelf"` // If set, the sender will be skipped
		Body    *string `json:"body"`
		BodyZh  *string `json:"bodyZh"`
		Title   *string `json:"title"`
		TitleZh *string `json:"titleZh"`
	}
	CloudMessage struct {
		UserId   intstring.IntString `json:"userId"`
		Topic    string              `json:"topic"`
		Priority string              `json:"priority"`
		Body     string              `json:"body"`
		BodyZh   string              `json:"bodyZh"`
		Title    string              `json:"title"`
		TitleZh  string              `json:"titleZh"`

		// Specify data for this message in addition to the notification itself.
		// Some data are added to the notification by default, and does not require
		// specifying again. These include:
		//	- "id"            /* Id of the notification */
		//	- "templateType"  /* Template type of the notification */
		//	- "actions"       /* Actions of the notification */
		Data map[string]string `json:"data"`
	}
)

func (n *Notification) AddUserRecipient(uid ...intstring.IntString) *Notification {
	n.Recipients.Users = append(n.Recipients.Users, uid...)
	return n
}

func (n *Notification) AddGroupRecipient(gid ...intstring.IntString) *Notification {
	n.Recipients.Groups = append(n.Recipients.Groups, gid...)
	return n
}

func (n *Notification) AddAction(actionId, actionLabel interface{}) *Notification {
	n.Actions = append(n.Actions, Action{
		ActionID:    strutil.StrOrEmptyFromInterface(actionId),
		ActionLabel: strutil.StrOrNotProvided(strutil.StrOrEmptyFromInterface(actionLabel)),
	})
	return n
}

// Attaches a mail payload to the notification.
// All mails MUST be of the same template id.
// Note that the notification module may decide not to send
// the email even though it is provided.
func (n *Notification) SetMail(mail ...Mail) *Notification {
	if len(mail) == 0 {
		return n
	}
	templateId := mail[0].Body.Id()
	for _, m := range mail[1:] {
		if m.Body.Id() != templateId {
			logger.Error("[SetMail] Sending mails with different templates within the same notification is not supported, mail ignored...")
			return n
		}
	}
	n.WithMail = &MailOptions{
		TemplateId: templateId,
		Mails:      mail,
	}
	return n
}

// Note: Use GeneratePush if possible unless the default behavior is not preffered
//
// Attaches a push notification payload to the notification.
// Note that the notification module may decide not to send
// the push notification even though it is provided.
func (n *Notification) AttachPushExtra(pushOptions ...*CloudMessage) *Notification {
	n.WithPush = append(n.WithPush, pushOptions...)
	return n
}

func (n *Notification) AttachPushAuto(autoPushParams AutoPushParams) *Notification {
	n.AutoPush = &autoPushParams
	return n
}

// A shortcut of AttachAutoPush for most default behaviour. `additionalBody`, if specified,
// will be appended below the original body of the notification.
func (n *Notification) AttachPushAutoDefault(additionalBody ...string) *Notification {
	n.AutoPush = &AutoPushParams{
		NoSelf: true,
	}
	if len(additionalBody) > 0 {
		mergedBody := fmt.Sprintf("%s\n%s", cloudMessageMagicDefault, strings.Join(additionalBody, ""))
		n.AutoPush.Body = strutil.NewPtr(mergedBody)
		n.AutoPush.BodyZh = strutil.NewPtr(mergedBody)
	}
	return n
}
