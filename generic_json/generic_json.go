package generic_json

import "fours/be-common-mdl/custom_types/intstring"

// JsonObject A simple wrapper type that provides methods to access a generic
// map[string]interface{} result obtained by encoding/json unmarshal
type JsonObject map[string]interface{}

func (j JsonObject) HasKey(key string) bool {
	_, ok := j[key]
	return ok
}

func (j JsonObject) ShouldGetString(key string) string {
	value, _ := j.GetString(key)
	return value
}

func (j JsonObject) GetString(key string) (string, bool) {
	obj, ok := j[key]
	if !ok {
		return "", false
	}
	result, ok := obj.(string)
	return result, ok
}

func (j JsonObject) ShouldGetBool(key string) bool {
	value, _ := j.GetBool(key)
	return value
}

func (j JsonObject) GetBool(key string) (bool, bool) {
	obj, ok := j[key]
	if !ok {
		return false, false
	}
	result, ok := obj.(bool)
	return result, ok
}

func (j JsonObject) ShouldGetNumber(key string) float64 {
	value, _ := j.GetNumber(key)
	return value
}

func (j JsonObject) GetNumber(key string) (float64, bool) {
	obj, ok := j[key]
	if !ok {
		return 0, false
	}
	result, ok := obj.(float64)
	return result, ok
}

func (j JsonObject) ShouldGetIntString(key string) intstring.IntString {
	value, _ := j.GetIntString(key)
	return value
}

func (j JsonObject) GetIntString(key string) (intstring.IntString, bool) {
	var result intstring.IntString
	obj, ok := j[key]
	if !ok {
		return result, false
	}
	s, ok := obj.(string)
	if ok {
		result.Scan(s)
		return result, true
	}
	n, ok := obj.(float64)
	if ok {
		result.Scan(int64(n))
		return result, true
	}
	return result, ok
}

func (j JsonObject) ShouldGetArray(key string) JsonArray {
	value, _ := j.GetArray(key)
	return value
}

func (j JsonObject) GetArray(key string) (jsonArr JsonArray, success bool) {
	obj, ok := j[key]
	if !ok {
		return
	}
	result, ok := obj.([]interface{})
	if !ok {
		return
	}
	return result, ok
}

func (j JsonObject) ShouldGetObj(key string) JsonObject {
	value, _ := j.GetObj(key)
	return value
}

func (j JsonObject) GetObj(key string) (jsonObj JsonObject, success bool) {
	obj, ok := j[key]
	if !ok {
		return
	}
	result, ok := obj.(map[string]interface{})
	if !ok {
		return
	}
	return result, true
}

// JsonArray A simple wrapper type that provides methods to access a generic
// []interface{} result obtained by encoding/json unmarshal
// The array must have no mixed type
type JsonArray []interface{}

func (j JsonArray) GetStringArray() (arr []string, success bool) {
	for _, obj := range j {
		result, ok := obj.(string)
		if !ok {
			return
		}
		arr = append(arr, result)
	}
	return arr, true
}

func (j JsonArray) GetBoolArray() (arr []bool, success bool) {
	for _, obj := range j {
		result, ok := obj.(bool)
		if !ok {
			return
		}
		arr = append(arr, result)
	}
	return arr, true
}

func (j JsonArray) GetNumberArray() (arr []float64, success bool) {
	for _, obj := range j {
		result, ok := obj.(float64)
		if !ok {
			return
		}
		arr = append(arr, result)
	}
	return arr, true
}

func (j JsonArray) GetIntStringArray() (arr []intstring.IntString, success bool) {
	for _, obj := range j {
		var result intstring.IntString
		s, ok := obj.(string)
		if ok {
			result.Scan(s)
			arr = append(arr, result)
			continue
		}
		n, ok := obj.(float64)
		if ok {
			result.Scan(int64(n))
			arr = append(arr, result)
			continue
		}
		return
	}
	return arr, true
}

func (j JsonArray) GetArrayArray() (arr []JsonArray, success bool) {
	for _, obj := range j {
		result, ok := obj.([]interface{})
		if !ok {
			return
		}
		arr = append(arr, result)
	}
	return arr, true
}

func (j JsonArray) GetObjectArray() (arr []JsonObject, success bool) {
	for _, obj := range j {
		result, ok := obj.(map[string]interface{})
		if !ok {
			return
		}
		arr = append(arr, result)
	}
	return arr, true
}
