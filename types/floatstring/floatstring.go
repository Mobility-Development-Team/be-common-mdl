package floatstring

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type FloatString float64

func FromString(str string) FloatString {
	c, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return FloatString(c)
}

func (f FloatString) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f *FloatString) UnmarshalJSON(b []byte) error {
	var s string
	if json.Unmarshal(b, &s) == nil {
		// String unmarshal successful
		if s == "" {
			// Empty string unmarshals to default value
			*f = 0
			return nil
		}
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*f = FloatString(v)
		return nil
	}
	// String unmarshal unsuccessful, try int
	return json.Unmarshal(b, (*float64)(f))
}

func (FloatString) GormDataType() string {
	return "double"
}

func (FloatString) GormDBDataType(*gorm.DB, *schema.Field) string {
	return "double"
}

func (f FloatString) Value() (driver.Value, error) {
	return float64(f), nil
}

func (f *FloatString) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*f = 0
	case float64:
		*f = FloatString(v)
	case bool:
		if v {
			*f = 1
		} else {
			*f = 0
		}
	case []byte:
		c, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return err
		} else {
			*f = FloatString(c)
		}
	case string:
		c, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		} else {
			*f = FloatString(c)
		}
	default:
		return errors.New("unable to cast value to IntString, no matching type")
	}
	return nil
}

func (f FloatString) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}
