package strutil

// StrOrEmpty Returns the string pointed by a string pointer str,
// if str is nil, return empty string "" instead
func StrOrEmpty(str *string) string {
	if str != nil {
		return *str
	}
	return ""
}

// NewPtr Copies str, and returns a pointer
func NewPtr(str string) *string {
	return &str
}
