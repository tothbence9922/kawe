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

func (sps *SimplePingerService) Configure() {
	for _, pingConfig := range configuration.GetInstance().PingConfigs {
		sps.methods = append(sps.methods, simpleMethod.SimplePingerMethod{Target: pingConfig.Target, Timeout: 5000, Method: "tcp", Periodicity: pingConfig.Periodicity})
	}
}

func StartMethod(currentMethod simpleMethod.PingerMethod) {

}

func (sps *SimplePingerService) StartMethods() {
	var wg sync.WaitGroup

	for _, method := range sps.methods {
		wg.Add(1)
		go func(currentMethod simpleMethod.PingerMethod) {
			defer wg.Done()
			for true {
				pingResponse, error := currentMethod.Ping()
				if error == nil && pingResponse != (simpleResponse.SimplePingResponse{}) {
					fmt.Printf("%s\n", pingResponse.String())
				}
				time.Sleep(time.Second * time.Duration(currentMethod.GetPeriodicity()))
			}
		}(method)
	}
	wg.Wait()
}
