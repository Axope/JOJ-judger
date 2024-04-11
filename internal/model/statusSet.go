package model

type StatusSet string

const (
	CE  StatusSet = "Compile Error"
	AC  StatusSet = "Accept"
	WA  StatusSet = "Wrong Answer"
	TLE StatusSet = "Time Limit Exceeded"
	MLE StatusSet = "Memory Limit Exceeded"
	RE  StatusSet = "Runtime Error"
	OLE StatusSet = "Output Limit Exceeded"
	UKE StatusSet = "Unknown Error"
)
