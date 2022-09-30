package strutil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// StrOrEmpty Returns the string pointed by a string pointer str,
// if str is nil, return empty string "" instead
func StrOrEmpty(str *string) string {
	return StrOrDefault(str, "")
}

// StrOrEmpty Returns the string pointed by a string pointer str,
// if str is nil, return defaultStr instead
func StrOrDefault(str *string, defaultStr string) string {
	if str != nil {
		return *str
	}
	return defaultStr
}

// EmptyOrDefault is the same as StrOrDefault but defaultStr is
// also returned when the string pointed by str is empty.
func EmptyOrDefault(str *string, defaultStr string) string {
	if str != nil && *str != "" {
		return *str
	}
	return defaultStr
}

// NewPtr Copies str, and returns a pointer
func NewPtr(str string) *string {
	return &str
}

func StrOrNotProvided(str string) string {
	if str == "" {
		return "Not provided"
	}
	return str
}

func StrOrEmptyFromInterface(obj interface{}) string {
	if obj == nil {
		return ""
	}
	value := reflect.Indirect(reflect.ValueOf(obj))
	if !value.IsValid() {
		// Value is invalid, e.g. pointer of nil
		return ""
	}
	return fmt.Sprint(value.Interface())
}

// A Less(i,j) function used by sort.Slice for ordering numbers before other strings
//  1. Number always go first and can be sorted in descending order provided by numberDesc
//  2. Other strings come after and can be sorted in desceding order provided by strDesc
//  3. If emptyStrLast is true, empty strings are put to the end of the list regardless of sort order
func CmpNumberFirst(str1 string, str2 string, numberDesc, strDesc, emptyStrLast bool) bool {
	num1Valid := false
	num1, err := strconv.Atoi(str1)
	if err == nil {
		num1Valid = true
	}
	num2Valid := false
	num2, err := strconv.Atoi(str2)
	if err == nil {
		num2Valid = true
	}
	switch {
	case num1Valid && num2Valid:
		// Both are numbers
		// Sort in numerical order
		return numberDesc != (num1 < num2)
	case !num1Valid && !num2Valid:
		// Both are string
		// Check if special handling needed for empty string
		if emptyStrLast && (str1 == "" || str2 == "") {
			// Empty string goes last
			return str1 != ""
		}
		// Sort in string order
		return strDesc != (str1 < str2)
	default:
		// Some of them is number
		// The one that is number goes first
		return num1Valid
	}
}

// ScreamCaseToLowerCamel converts SCREAM_CASE to screamCase
func ScreamCaseToLowerCamel(scream string) string {
	words := strings.Split(scream, "_")
	switch len(words) {
	case 0:
		return ""
	case 1:
		return strings.ToLower(words[0])
	default:
		words[0] = strings.ToLower(words[0])
		for i := 1; i < len(words); i++ {
			words[i] = strings.Title(strings.ToLower(words[i]))
		}
		return strings.Join(words, "")
	}
}

// ScreamCaseToTitle converts SCREAM_CASE to Scream Case
func ScreamCaseToTitle(scream string) string {
	words := strings.Split(scream, "_")
	if len(words) == 0 {
		return ""
	}
	for i := range words {
		words[i] = strings.Title(strings.ToLower(words[i]))
	}
	return strings.Join(words, " ")
}
