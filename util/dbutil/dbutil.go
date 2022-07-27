package dbutil

import (
	"gorm.io/gorm"
)

// PreloadHandler handle gorm DB preload
func PreloadHandler(query *gorm.DB, preloads ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	if len(preloads) == 0 {
		return query
	}
	q := query
	for _, preload := range preloads {
		if preload == nil {
			// Just in case nil is passed, skip it
			continue
		}
		q = q.Scopes(preload)
	}
	return q
}
