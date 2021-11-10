package simple

import (
	"encoding/json"
	"fmt"

	simple "github.com/tothbence9922/kawe/internal/ping/simple/response"
)

type SimplePingResult struct {
	Response    simple.PingResponse
	ServiceName string
}

func (spr SimplePingResult) String() string {

	return fmt.Sprintf("%s\t%s", spr.ServiceName, spr.Response.String())
}
func (spr SimplePingResult) Json() string {

	jsonString, _ := json.Marshal(spr)
	return string(jsonString)
}
