package intstring

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type IntString int

func FromString(str string) IntString {
	c, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return IntString(c)
}

func (i IntString) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *IntString) UnmarshalJSON(b []byte) error {
	var s string
	if json.Unmarshal(b, &s) == nil {
		// String unmarshal successful
		if s == "" {
			// Empty string unmarshals to default value
			*i = 0
			return nil
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*i = IntString(v)
		return nil
	}
	// String unmarshal unsuccessful, try int
	return json.Unmarshal(b, (*int)(i))
}

func (IntString) GormDataType() string {
	return "int"
}

func (IntString) GormDBDataType(*gorm.DB, *schema.Field) string {
	return "int"
}

func (i IntString) Value() (driver.Value, error) {
	return int64(i), nil
}

func (i *IntString) ShouldScan(value interface{}) *IntString {
	i.Scan(value)
	return i
}

func (i *IntString) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*i = 0
	case int64:
		*i = IntString(v)
	case bool:
		if v {
			*i = 1
		} else {
			*i = 0
		}
	case []byte:
		c, err := strconv.Atoi(string(v))
		if err != nil {
			return err
		} else {
			*i = IntString(c)
		}
	case string:
		c, err := strconv.Atoi(v)
		if err != nil {
			return err
		} else {
			*i = IntString(c)
		}
	default:
		return errors.New("unable to cast value to IntString, no matching type")
	}
	return nil
}

func (i IntString) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// A quick function for compatibility of legacy code
func ToIntSlice(values []IntString) []int {
	result := make([]int, len(values))
	for i, value := range values {
		result[i] = int(value)
	}
	return result
}
