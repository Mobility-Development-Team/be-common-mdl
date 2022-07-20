package util

import (
	logger "github.com/sirupsen/logrus"
	"net/url"
)

// GetUrlWithQuery Gets an url with GET parameters appended to path. Returns empty string if it fails or path is not valid
func GetUrlWithQuery(path string, params map[string]string) string {
	u, err := url.Parse(path)
	if err != nil {
		logger.Errorf("[GetUrlWithQuery] Error generating url %s [%s]", path, err.Error())
		return ""
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
