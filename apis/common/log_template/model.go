package log_template

import "fmt"

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
