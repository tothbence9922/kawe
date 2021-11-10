package simple

import (
	"sync"
	"time"

	"github.com/tothbence9922/kawe/internal/configuration"
	simpleMethod "github.com/tothbence9922/kawe/internal/ping/simple/method"
	simpleResult "github.com/tothbence9922/kawe/internal/ping/simple/result"
)

type SimplePingerService struct {
	methods []simpleMethod.PingerMethod
	Name    string
	Channel chan (simpleResult.PingResult)
}

func (sps SimplePingerService) String() string {

	var ret string

	for _, method := range sps.methods {
		ret = ret + method.String()
	}
	return ret
}

func (sps *SimplePingerService) Configure(config configuration.ServiceConfiguration, channel chan (simpleResult.PingResult)) {
	sps.Name = config.Name
	sps.Channel = channel
	for _, pingConfig := range config.PingConfigs {
		sps.methods = append(sps.methods, simpleMethod.SimplePingerMethod{Target: pingConfig.Target, Timeout: 5000, Method: "tcp", Periodicity: pingConfig.Periodicity})
	}
}

func (sps *SimplePingerService) StartMethod(wg sync.WaitGroup, method simpleMethod.PingerMethod) {

	go func(method simpleMethod.PingerMethod, outChannel chan<- (simpleResult.PingResult)) {
		defer wg.Done()
		for true {
			pingResponse, error := method.Ping()
			if error == nil {
				outChannel <- simpleResult.SimplePingResult{ServiceName: sps.Name, Response: pingResponse}
			}
			time.Sleep(time.Second * time.Duration(method.GetPeriodicity()))
		}
	}(method, sps.Channel)
}

func (sps *SimplePingerService) StartMethods(wg sync.WaitGroup) {

	for _, method := range sps.methods {
		sps.StartMethod(wg, method)
	}
}
