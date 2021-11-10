package simple

import (
	"fmt"
	"sync"
	"time"

	"github.com/tothbence9922/kawe/internal/configuration"
	simpleMethod "github.com/tothbence9922/kawe/internal/ping/simple/method"
	simpleResponse "github.com/tothbence9922/kawe/internal/ping/simple/response"
)

type SimplePingerService struct {
	methods []simpleMethod.PingerMethod
}

func (sps SimplePingerService) String() string {
	var ret string

	for _, method := range sps.methods {
		ret = ret + method.String()
	}
	return ret
}

func (sps *SimplePingerService) Configure(config configuration.ServiceConfiguration) {

	for _, pingConfig := range config.PingConfigs {
		sps.methods = append(sps.methods, simpleMethod.SimplePingerMethod{Target: pingConfig.Target, Timeout: 5000, Method: "tcp", Periodicity: pingConfig.Periodicity})
	}

}

func StartMethod(wg sync.WaitGroup, method simpleMethod.PingerMethod) {
	go func(method simpleMethod.PingerMethod) {
		defer wg.Done()
		for true {
			pingResponse, error := method.Ping()
			if error == nil && pingResponse != (simpleResponse.SimplePingResponse{}) {
				fmt.Printf("%s\n", pingResponse.String())
			}
			time.Sleep(time.Second * time.Duration(method.GetPeriodicity()))
		}
	}(method)
}

func (sps *SimplePingerService) StartMethods(wg sync.WaitGroup) {

	for _, method := range sps.methods {
		StartMethod(wg, method)
	}

}
