package model

type TestCase struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Note        string `json:"note,omitempty"`
	Explanation string `json:"explanation,omitempty"`
}
