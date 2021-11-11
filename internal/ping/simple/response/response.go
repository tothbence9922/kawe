package simple

import (
	"encoding/json"
	"fmt"
	"time"
)

type SimplePingResponse struct {
	Success   bool
	Error     string
	Target    string
	Timestamp time.Time
}

func (spr SimplePingResponse) String() string {

	successText := "Successful"
	errorText := ""
	if !spr.Success {
		successText = "Failed"
		errorText = spr.Error
	}
	return fmt.Sprintf("%s\tSimplePing\t->\t%s\t\t%s.\t%s\n", spr.Timestamp.Format("2 Jan 2006 15:04:05"), spr.Target, successText, errorText)
}
func (spr SimplePingResponse) Json() string {

	jsonString, _ := json.Marshal(spr)
	return string(jsonString)
}
