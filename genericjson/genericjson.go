package genericjson

import "github.com/Mobility-Development-Team/be-common-mdl/types/intstring"

// Object A simple wrapper type that provides methods to access a generic
// map[string]interface{} result obtained by encoding/json unmarshal
type Object map[string]interface{}

func (j Object) HasKey(key string) bool {
	_, ok := j[key]
	return ok
}

func (j Object) ShouldGetString(key string) string {
	value, _ := j.GetString(key)
	return value
}

func (j Object) GetString(key string) (string, bool) {
	obj, ok := j[key]
	if !ok {
		return "", false
	}
	result, ok := obj.(string)
	return result, ok
}

func (j Object) ShouldGetBool(key string) bool {
	value, _ := j.GetBool(key)
	return value
}

func (j Object) GetBool(key string) (bool, bool) {
	obj, ok := j[key]
	if !ok {
		return false, false
	}
	result, ok := obj.(bool)
	return result, ok
}

func (j Object) ShouldGetNumber(key string) float64 {
	value, _ := j.GetNumber(key)
	return value
}

func (j Object) GetNumber(key string) (float64, bool) {
	obj, ok := j[key]
	if !ok {
		return 0, false
	}
	result, ok := obj.(float64)
	return result, ok
}

func (j Object) ShouldGetIntString(key string) intstring.IntString {
	value, _ := j.GetIntString(key)
	return value
}

func (j Object) GetIntString(key string) (intstring.IntString, bool) {
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

func (j Object) ShouldGetArray(key string) Array {
	value, _ := j.GetArray(key)
	return value
}

func (j Object) GetArray(key string) (jsonArr Array, success bool) {
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

func (j Object) ShouldGetObj(key string) Object {
	value, _ := j.GetObj(key)
	return value
}

func (j Object) GetObj(key string) (jsonObj Object, success bool) {
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

// Array A simple wrapper type that provides methods to access a generic
// []interface{} result obtained by encoding/json unmarshal
// The array must have no mixed type
type Array []interface{}

func (j Array) GetStringArray() (arr []string, success bool) {
	for _, obj := range j {
		result, ok := obj.(string)
		if !ok {
			return
		}
		arr = append(arr, result)
	}
	return arr, true
}

func (j Array) GetBoolArray() (arr []bool, success bool) {
	for _, obj := range j {
		result, ok := obj.(bool)
		if !ok {
			return
		}
		arr = append(arr, result)
	}
	return arr, true
}

func (j Array) GetNumberArray() (arr []float64, success bool) {
	for _, obj := range j {
		result, ok := obj.(float64)
		if !ok {
			return
		}
		arr = append(arr, result)
	}
	return arr, true
}

func (j Array) GetIntStringArray() (arr []intstring.IntString, success bool) {
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

func (j Array) GetArrayArray() (arr []Array, success bool) {
	for _, obj := range j {
		result, ok := obj.([]interface{})
		if !ok {
			return
		}
		arr = append(arr, result)
	}
	return arr, true
}

func (j Array) GetObjectArray() (arr []Object, success bool) {
	for _, obj := range j {
		result, ok := obj.(map[string]interface{})
		if !ok {
			return
		}
		arr = append(arr, result)
	}
	return arr, true
}