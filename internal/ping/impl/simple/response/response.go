package simple

import (
	"encoding/json"
	"fmt"
	"time"
)

type PingResponse struct {
	Success   bool
	Error     string
	Target    string
	Timestamp time.Time
}

func (pr PingResponse) GetTarget() string {
	return pr.Target
}

func (pr PingResponse) String() string {

	successText := "Successful"
	errorText := ""
	if !pr.Success {
		successText = "Failed"
		errorText = pr.Error
	}
	return fmt.Sprintf("%s\tSimplePing\t->\t%s\t\t%s.\t%s\n", pr.Timestamp.Format("2 Jan 2006 15:04:05"), pr.Target, successText, errorText)
}
func (pr PingResponse) Json() string {

	jsonString, _ := json.Marshal(pr)
	return string(jsonString)
}
