package document

type FollowUpReportInfo struct {
	Title       string `json:"title"`
	DueDate     string `json:"dueDate"`
	Description string `json:"description"`
	Images      []struct {
		Before string `json:"before"`
		After  string `json:"after"`
	} `json:"images"`
}
