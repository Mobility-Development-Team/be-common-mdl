package util

// StrOrEmpty Returns the string pointed by a string pointer str,
// if str is nil, return empty string "" instead
func StrOrEmpty(str *string) string {
	if str != nil {
		return *str
	}
	return ""
}

// NewStrPtr Copies str, and returns a pointer
func NewStrPtr(str string) *string {
	return &str
}
