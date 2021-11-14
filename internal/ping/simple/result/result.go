package simple

import (
	"encoding/json"
	"fmt"

	simple "github.com/tothbence9922/kawe/internal/ping/simple/response"
)

type SimplePingResult struct {
	Responses   map[string](simple.PingResponse)
	ServiceName string
}

func (spr SimplePingResult) GetResponses() map[string](simple.PingResponse) {

	return spr.Responses
}

func (spr SimplePingResult) GetServiceName() string {

	return spr.ServiceName
}

func (spr SimplePingResult) String() string {

	ret := ""

	for _, value := range spr.Responses {
		ret += value.String() + "\n"
	}

	return fmt.Sprintf("%s\t%s", spr.ServiceName, ret)
}

func (spr SimplePingResult) Json() string {

	jsonString, _ := json.Marshal(spr)
	return string(jsonString)
}
