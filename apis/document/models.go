package document

type FollowUpReportInfo struct {
	ContractName string `json:"contractName"`
	Title        string `json:"title"`
	DueDate      string `json:"dueDate"`
	Description  string `json:"description"`
	Images       []struct {
		Before string `json:"before"`
		After  string `json:"after"`
	}
}
