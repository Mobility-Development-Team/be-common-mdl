package boolutil

func False() *bool {
	value := false
	return &value
}

func True() *bool {
	value := true
	return &value
}

// GetBool Returns the bool pointed by a bool pointer b,
// if b is nil, return false instead
func GetBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}
