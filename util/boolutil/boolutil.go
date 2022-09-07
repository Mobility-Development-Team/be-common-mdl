package boolutil

func False() *bool {
	value := false
	return &value
}

func True() *bool {
	value := true
	return &value
}
