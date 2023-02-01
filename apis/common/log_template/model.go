package log_template

import "fmt"

// This common log template will be shared between Machine & Labour
const (
	ActivityTypeCreate   = "CREATE"
	ActivityTypePublish  = "PUBLISH"
	ActivityTypeUpdate   = "UPDATE"
	ActivityTypeComplete = "COMPLETE"
	// ActivityTypeADD For worker point
	ActivityTypeADD   = "ADD"
	ActivityTypeMINUS = "MINUS"
	ActivityTypeBEGIN = "BEGIN"
)

type ActivityLogTemplate struct {
	FormatChiStr string
	FormatEngStr string
	ActivityType string
}

func (t ActivityLogTemplate) FormatChinese(a ...interface{}) string {
	return fmt.Sprintf(t.FormatChiStr, a...)
}

func (t ActivityLogTemplate) FormatEnglish(a ...interface{}) string {
	return fmt.Sprintf(t.FormatEngStr, a...)
}

func (t ActivityLogTemplate) GetActivityType() string {
	return t.ActivityType
}
