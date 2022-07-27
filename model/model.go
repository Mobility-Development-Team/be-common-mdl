package model

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/Mobility-Development-Team/be-common-mdl/apis/auth"
	"github.com/Mobility-Development-Team/be-common-mdl/genericjson"
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	Model struct {
		Id               intstring.IntString `gorm:"primaryKey" json:"id,omitempty"`
		CreatedAt        time.Time           `json:"createdAt" gorm:"<-:create"`
		CreatedBy        string              `json:"-" gorm:"column:created_by;<-:create"`
		CreatedByDisplay interface{}         `json:"createdBy" gorm:"-" `
		UpdatedAt        time.Time           `json:"updatedAt"`
		UpdatedBy        *string             `json:"-" gorm:"column:updated_by"`
		UpdatedByDisplay interface{}         `json:"updatedBy" gorm:"-" `
	}
	UserInfo struct {
		Model
		Uuid           string  `json:"uuid"`
		FirstName      string  `json:"firstName"`
		LastName       string  `json:"lastName"`
		DisplayName    string  `json:"displayName"`
		Email          string  `json:"email"`
		Position       *string `json:"position"`
		ContactNo      *string `json:"contactNo"`
		DefaultLang    *string `json:"defaultLang"`
		Status         string  `json:"status" binding:"required"`
		ProfileIconUrl *string `json:"profileIconUrl"`
		UserRefKey     string  `json:"userRefKey"`
	}
	UserDisplay struct {
		FirstName      string  `json:"firstName"`
		LastName       string  `json:"lastName"`
		DisplayName    string  `json:"displayName"`
		ProfileIconUrl *string `json:"profileIconUrl"`
	}
	Location struct {
		Id            intstring.IntString  `json:"id"`
		Uuid          string               `json:"uuid"`
		Name          string               `json:"name"`
		NameZh        *string              `json:"nameZh"`
		Status        string               `json:"status"`
		LocationType  *string              `json:"locationType"`
		Latitude      float64              `json:"latitude"`
		Longitude     float64              `json:"longitude"`
		LocationId    *intstring.IntString `json:"locationId"`
		ContractRefId *intstring.IntString `json:"contractId"`
	}
	Contract struct {
		ContractNo   string  `json:"contractNo"`
		ContractDesc *string `json:"contractDesc"`
	}
	MediaParam struct {
		Id                   intstring.IntString `json:"id"`
		FbRefId              string              `json:"fbRefId"`
		FirebaseUrlThumbnail string              `json:"firebaseUrlThumbnail"`
		CreatedAt            string              `json:"createdAt"`
		MediaRefType         string              `json:"mediaRefType"`
		MediaRefInfo         json.RawMessage     `json:"mediaRefInfo"`
	}
)

// Attempts to set MediaRefInfo, ignores and logs a message if fails
func (mp *MediaParam) ShouldSetRefInfo(refType string, obj interface{}) *MediaParam {
	b, err := json.Marshal(obj)
	if err != nil {
		logger.Errorf("[MediaParam][ShouldSetRefInfo] Unable to set MediaRefInfo: %s Value: %v", err.Error())
		return mp
	}
	mp.MediaRefType = refType
	mp.MediaRefInfo = b
	return mp
}

// Attempts to get MediaRefInfo, ignores if fails.
// The returned value is a JsonObject for quick retrieval.
// To unmarshal into struct, use json.Unmarshal()
func (mp MediaParam) ShouldGetRefInfo() (result genericjson.Object) {
	if err := json.Unmarshal(mp.MediaRefInfo, &result); err != nil {
		logger.Debugf("[MediaParam][ShouldGetRefInfo] No valid refInfo for media: %s (%s)", mp.Id, err.Error())
	}
	return result
}

func GetUserFromMap(userRefKey string, m map[string]UserInfo) *UserInfo {
	if userRefKey == "" {
		return nil
	}
	u, ok := m[userRefKey]
	if !ok {
		return nil
	}
	return &u
}

func (m *Model) AfterFind(tx *gorm.DB) (err error) {
	m.CreatedByDisplay = m.CreatedBy
	m.UpdatedByDisplay = m.UpdatedBy
	return
}

func (m *Model) SetCreatedBy(refKey string) {
	m.CreatedBy = refKey
	m.CreatedByDisplay = &m.CreatedBy
}

func (m *Model) ShouldAddCreatedBy(c *gin.Context) {
	m.CreatedBy = auth.GetUserRefKeyFromContext(c)
	m.CreatedByDisplay = &m.CreatedBy
}

func (m *Model) ShouldAddUpdatedBy(c *gin.Context) {
	ub := auth.GetUserRefKeyFromContext(c)
	m.UpdatedBy = &ub
	m.UpdatedByDisplay = &m.UpdatedBy
}

func (m *Model) ShouldAddSystemFields(c *gin.Context) {
	m.ShouldAddCreatedBy(c)
	if m.Id != 0 {
		m.ShouldAddUpdatedBy(c)
	}
}

func (ui *UserInfo) IsEmpty() bool {
	return reflect.DeepEqual(ui, UserInfo{})
}

// Gets the name of the user, returns the display name if available,
// otherwise, returns first name + last name. If all not available, use email.
// If the user is nil, returns [Unknown User] instead.
func (ui *UserInfo) GetAvailableName() string {
	if ui == nil {
		return "[Unknown User]"
	}
	return GetAvailableName(ui.DisplayName, ui.FirstName, ui.LastName, ui.Email)
}

func GetAvailableName(displayName, firstName, lastName, placeholder string) string {
	if displayName != "" {
		return displayName
	}
	if firstName != "" || lastName != "" {
		return strings.Join([]string{firstName, lastName}, " ")
	} else {
		return placeholder
	}
}

func LocationSliceToString(locationList []*Location) string {
	locations := map[string]struct{}{}
	for _, l := range locationList {
		if l == nil {
			continue
		}
		name := l.Name
		if l.NameZh != nil {
			name += " " + *l.NameZh
		}
		locations[name] = struct{}{}
	}
	locationStr := make([]string, 0, len(locations))
	for k := range locations {
		locationStr = append(locationStr, k)
	}
	if len(locationStr) == 0 {
		return "<Not given>"
	}
	return strings.Join(locationStr, "; ")
}

func (m *MediaParam) UnmarshalJSON(b []byte) error {
	type origMediaParam MediaParam
	var origM origMediaParam // Create new type so it doesn't reuse this unmarshaller
	if json.Unmarshal(b, &origM) == nil {
		// Unmarshal successful just by struct
		a := MediaParam(origM)
		*m = a
		return nil
	}
	// Struct unmarshal unsuccessful, it must be a string
	var str string
	if json.Unmarshal(b, &str) == nil {
		// Try parse as intstring
		if id := intstring.FromString(str); id != 0 {
			*m = MediaParam{
				Id: id,
			}
		}
		// If not, treat as refId
		*m = MediaParam{
			FbRefId: str,
		}
		return nil
	}
	return errors.New("cannot unmarshal MediaParam, not a struct nor a string of id")
}