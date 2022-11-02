package dbutil

import (
	"errors"

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

// GetAutoIncrement gets the auto increment of a table.
// db must not be nil
func GetAutoIncrement(db *gorm.DB, tableName string) (int64, error) {
	if db == nil {
		return 0, errors.New("db must not be nil")
	}
	var ai int64
	if err := db.
		Table("INFORMATION_SCHEMA.TABLES").
		Select("AUTO_INCREMENT").
		Where("TABLE_SCHEMA = DATABASE()").
		Where("TABLE_NAME = ?", tableName).
		Scan(&ai).
		Error; err != nil {
		return 0, err
	}
	return ai, nil
}
