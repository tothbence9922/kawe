package simple

import (
	"encoding/json"
	"fmt"
	"sync"

	interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type PingResult struct {
	sync.RWMutex
	Responses   map[string](interfaces.IPingResponse)
	ServiceName string
}

func (spr *PingResult) GetResponses() map[string](interfaces.IPingResponse) {
	spr.Lock()
	defer spr.Unlock()
	m := make(map[string](interfaces.IPingResponse), len(spr.Responses))
	for k, v := range spr.Responses {
		m[k] = v
	}
	return m
}

func (spr *PingResult) AddResponse(newResponse interfaces.IPingResponse) {

	spr.Lock()
	defer spr.Unlock()

	m := make(map[string](interfaces.IPingResponse), len(spr.Responses))
	for k, v := range spr.Responses {
		m[k] = v
	}

	m[newResponse.GetTarget()] = newResponse
	spr.Responses = m
}

func (spr *PingResult) GetServiceName() string {
	ret := spr.ServiceName
	return ret
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
