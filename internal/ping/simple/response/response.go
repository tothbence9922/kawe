package simple

import "fmt"

type SimplePingResponse struct {
	Success bool
	Error   string
}

func (spr SimplePingResponse) String() string {
	successText := "successful"
	errorText := ""
	if !spr.Success {
		successText = "failed"
		errorText = spr.Error
	}
	return fmt.Sprintf("SimplePing %s.\t\t%s\n", successText, errorText)
}
