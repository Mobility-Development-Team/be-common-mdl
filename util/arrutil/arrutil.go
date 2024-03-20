package arrutil

import (
	"github.com/Mobility-Development-Team/be-common-mdl/types/intstring"
	"strings"
)

// DiffStrSlice check the difference between two given string slices
func DiffStrSlice(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// Contains check if the given intStr value exist in Slice intStrSlice
func Contains(intStrSlice []intstring.IntString, intStr intstring.IntString) bool {
	for _, v := range intStrSlice {
		if v != intStr {
			continue
		}
		return true
	}
	return false
}

// ContainsStr check if the given intStr value exist in Slice intStrSlice
func ContainsStr(strSlice []string, str string) bool {
	for _, v := range strSlice {
		if v != str {
			continue
		}
		return true
	}
	return false
}

// IsWordContainStr check if the given intStr value exist in Slice intStrSlice
func IsWordContainStr(str string, strWords ...string) bool {
	for _, v := range strWords {
		if strings.Contains(str, v) {
			continue
		}
		return true
	}
	return false
}

// Unique distinct the given IntString slice
func Unique(intSlice []intstring.IntString) []intstring.IntString {
	var (
		ls []intstring.IntString
	)
	keys := make(map[intstring.IntString]bool)
	ls = []intstring.IntString{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			ls = append(ls, entry)
		}
	}
	return ls
}

// UniqueStrSlice distinct the given String slice
func UniqueStrSlice(strSlice []string) []string {
	var (
		ls []string
	)
	keys := make(map[string]bool)
	ls = []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			ls = append(ls, entry)
		}
	}
	return ls
}
