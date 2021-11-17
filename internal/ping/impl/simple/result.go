package simple

import (
	"encoding/json"
	"fmt"
	"sync"

	interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type PingResult struct {
	sync.Mutex
	Responses   map[string](interfaces.IPingResponse)
	ServiceName string
}

func (spr *PingResult) GetResponses() map[string](interfaces.IPingResponse) {

	return spr.Responses
}

func (spr *PingResult) AddResponse(newResponse interfaces.IPingResponse) {

	spr.Lock()
	defer spr.Unlock()

	spr.Responses[newResponse.GetTarget()] = newResponse
}

func (spr *PingResult) GetServiceName() string {

	return spr.ServiceName
}

func (spr *PingResult) String() string {

	ret := ""

	for _, value := range spr.Responses {
		ret += value.String() + "\n"
	}

	return fmt.Sprintf("%s\t%s", spr.ServiceName, ret)
}

func (spr *PingResult) Json() string {

	jsonString, _ := json.Marshal(spr)
	return string(jsonString)
}
