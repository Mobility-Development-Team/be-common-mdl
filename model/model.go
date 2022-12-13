package model

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/Mobility-Development-Team/be-common-mdl/apis/auth"
	"github.com/Mobility-Development-Team/be-common-mdl/genericjson"
	"github.com/Mobility-Development-Team/be-common-mdl/types/floatstring"
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
		// Field level permission `<-:update` is removed to workaround a gorm issue for assoication update
		// See Model.BeforeCreate() below
		UpdatedBy        *string     `json:"-" gorm:"column:updated_by"`
		UpdatedByDisplay interface{} `json:"updatedBy" gorm:"-" `
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
	PartyInfo struct {
		Id            intstring.IntString `json:"id"`
		PartyName     string              `json:"partyName"`
		PartyNameZh   string              `json:"partyNameZh"`
		Address       string              `json:"address"`
		Email         string              `json:"email"`
		Br            string              `json:"br"`
		TradeCategory string              `json:"tradeCategory"`
		PartyIconUrl  string              `json:"partyIconUrl"`
		PartyPrefix   string              `json:"partyPrefix"`
		SubconRefId   string              `json:"subconRefId"`
	}
	GroupInfo struct {
		Model
		Uuid          string              `json:"uuid"`
		Name          string              `json:"name"`
		Status        string              `json:"status"`
		ContractRefId intstring.IntString `json:"contractId"`
		PartyRefId    intstring.IntString `json:"partyId"`
		IsSystemGroup *bool               `json:"isSystemGroup"`
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
		ProjectIdRef string  `json:"projectIdRef"`
	}
	MediaParam struct {
		Id                   intstring.IntString     `json:"id"`
		CreatedAt            string                  `json:"createdAt"`
		CreatedBy            json.RawMessage         `json:"createdBy"`
		UpdatedAt            string                  `json:"updatedAt"`
		UpdatedBy            json.RawMessage         `json:"updatedBy"`
		FbRefId              string                  `json:"fbRefId"`
		FbCreatedBy          string                  `json:"fbCreatedBy"`
		FbUpdatedBy          *string                 `json:"fbUpdatedBy"`
		BatchId              string                  `json:"batchId"`
		LocalPath            string                  `json:"localPath"`
		LocalPathThumbnail   string                  `json:"localPathThumbnail"`
		RefUserName          string                  `json:"userName"`
		Description          string                  `json:"description"`
		UploadStatus         string                  `json:"uploadStatus"`
		FirebaseUrl          string                  `json:"firebaseUrl"`
		FirebaseUrlThumbnail string                  `json:"firebaseUrlThumbnail"`
		Latitude             floatstring.FloatString `json:"latitude"`
		Longitude            floatstring.FloatString `json:"longitude"`
		DeviceType           string                  `json:"deviceType"`
		ContractId           string                  `json:"contractId"`
		MediaType            string                  `json:"mediaType"`
		MediaRefType         string                  `json:"mediaRefType"`
		FbCreatedAt          string                  `json:"fbCreatedAt"`
		FbUpdatedAt          string                  `json:"fbUpdatedAt"`
		Hashtags             json.RawMessage         `json:"hashtags"`
		MediaRefInfo         json.RawMessage         `json:"mediaRefInfo"`
	}
)

// GetIdFromInterface attempts to get the Id field of obj
//
// obj must be a struct with field ``Id`` of type intstring.IntString,
// embedding Model is not required as long as the requirement is met.
func GetIdFromInterface(obj interface{}) (intstring.IntString, error) {
	if obj == nil {
		return 0, errors.New("obj is nil")
	}
	objStruct := reflect.Indirect(reflect.ValueOf(obj))
	if objStruct.Kind() != reflect.Struct {
		return 0, errors.New("obj is not struct")
	}
	idField := objStruct.FieldByName("Id")
	if !idField.IsValid() {
		return 0, errors.New("no field named 'Id'")
	}
	if id, ok := reflect.Indirect(idField).Interface().(intstring.IntString); !ok {
		return 0, errors.New("field is not of type 'intstring.IntString'")
	} else {
		return id, nil
	}
}

// GetIdsFromRecords attempts to get all the ids from the given slice
//
// All slice values must be structs with field ``Id`` of type intstring.IntString.
// Any value failing to meet such requirement result in an error.
func GetIdsFromRecords(slice interface{}) ([]intstring.IntString, error) {
	if slice == nil {
		return nil, errors.New("slice is nil")
	}
	sliceValue := reflect.Indirect(reflect.ValueOf(slice))
	if sliceValue.Kind() != reflect.Slice && sliceValue.Kind() != reflect.Array {
		return nil, errors.New("obj is not a slice or array")
	}
	sliceLen := sliceValue.Len()
	ids := make([]intstring.IntString, sliceLen)
	for i := 0; i < sliceLen; i++ {
		id, err := GetIdFromInterface(sliceValue.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

// Attempts to set MediaRefInfo, ignores and logs a message if fails
func (mp *MediaParam) ShouldSetRefInfo(refType string, obj interface{}) *MediaParam {
	b, err := json.Marshal(obj)
	if err != nil {
		logger.Errorf("[MediaParam][ShouldSetRefInfo] Unable to set MediaRefInfo: %s", err.Error())
		return mp
	}
	mp.MediaRefType = refType
	mp.MediaRefInfo = b
	return mp
}

// Attempts to get MediaRefInfo, ignores if fails.
// The returned value is a JsonObject for quick retrieval.-
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

// BeforeCreate is a hook for updating the UpdatedBy field when update is called in a Create context.
//
// Assoication updates are called internally by gorm with INSERT OR UPDATE for upsert
// Therefore, field update permission will be ignored and all update hooks are not effective.
//
// To workaround this issue with setting UpdatedBy, we assume creating record with an Id is
// always an update and clears the UpdatedBy field if set. This hook might be shadowed by
// other model's custom hook if they also implement BeforeCreate.
func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Id == 0 {
		m.UpdatedBy = nil
	}
	return
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
	m.ShouldAddUpdatedBy(c)
	// This is now controlled by declared update permission
	// if m.Id != 0 {
	// 	m.ShouldAddUpdatedBy(c)
	// }
}

// Attempts to parse and extract userRefKey from CreatedByDisplay and UpdatedByDisplay
// and set to their corresponding fields. Useful after a json unmarshal during internal calls
func (m *Model) ShouldAddSystemFieldsFromDisplay() *Model {
	switch createdBy := m.CreatedByDisplay.(type) {
	case string:
		m.CreatedBy = createdBy
	case map[string]interface{}:
		m.CreatedBy = genericjson.Object(createdBy).ShouldGetString("userRefKey")
	}
	switch updatedBy := m.UpdatedByDisplay.(type) {
	case string:
		m.UpdatedBy = &updatedBy
	case map[string]interface{}:
		if value, ok := genericjson.Object(updatedBy).GetString("userRefKey"); ok {
			m.UpdatedBy = &value
		}
	}
	return m
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

// SanitizeForCreate clears all models inside a struct for creation.
// It looks inside nested structs and arrays
func SanitizeForCreate(mdl interface{}, creatorRefKey string) interface{} {
	return replaceModel(reflect.ValueOf(mdl), Model{
		CreatedBy: creatorRefKey,
	}).Interface()
}

func replaceModel(v reflect.Value, replacement Model) reflect.Value {
	if !v.IsValid() {
		return v
	}
	v = reflect.Indirect(v)
	if !v.IsValid() {
		return v
	}
	switch kind := v.Kind(); kind {
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(replacement) {
			if v.CanSet() {
				v.Set(reflect.ValueOf(replacement))
			}
			return v
		}
		for _, f := range reflect.VisibleFields(v.Type()) {
			if fv := v.FieldByIndex(f.Index); fv.IsValid() {
				replaceModel(fv, replacement)
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if iv := v.Index(i); iv.IsValid() {
				replaceModel(iv, replacement)
			}
		}
	}
	return v
}
