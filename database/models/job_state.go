package models

import "fmt"

type JobState int

const (
	Queued JobState = iota
	Running
	Retrying
	Failed
	Succeeded
)

func (js JobState) String() string {
	switch js {
	case 0:
		return "Queued"
	case 1:
		return "Running"
	case 2:
		return "Retrying"
	case 3:
		return "Failed"
	case 4:
		return "Succeeded"
	}
	return fmt.Sprintf("JobState(%v)", int(js))
}
