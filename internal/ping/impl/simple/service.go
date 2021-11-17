package simple

import (
	"fmt"
	"sync"
	"time"

	"github.com/tothbence9922/kawe/internal/aggregator"
	"github.com/tothbence9922/kawe/internal/configuration"
	interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type SimplePingerService struct {
	methods []interfaces.IPingMethod
	Name    string
	Channel chan (interfaces.IPingResult)
	Result  interfaces.IPingResult
}

func (sps SimplePingerService) String() string {

	var ret string

	for _, method := range sps.methods {
		ret = ret + method.String()
	}
	return ret
}

func (sps *SimplePingerService) Configure(config configuration.ServiceConfiguration, channel chan (interfaces.IPingResult)) {

	sps.Name = config.Name
	sps.Channel = channel
	sps.Result = &SimplePingResult{ServiceName: sps.Name, Responses: make(map[string](interfaces.IPingResponse))}
	for _, pingConfig := range config.PingConfigs {
		sps.methods = append(sps.methods, SimplePingerMethod{Target: pingConfig.Target, Timeout: 5000, Method: "tcp", Periodicity: pingConfig.Periodicity})
	}
}

func (sps *SimplePingerService) StartMethod(wg *sync.WaitGroup, method interfaces.IPingMethod) {

	go func(method interfaces.IPingMethod, outChannel chan<- (interfaces.IPingResult)) {

		defer wg.Done()
		for true {
			pingResponse, error := method.Ping()
			if error == nil {
				sps.Result.GetResponses()[pingResponse.GetTarget()] = pingResponse
				outChannel <- sps.Result
			}
			time.Sleep(time.Second * time.Duration(method.GetPeriodicity()))
		}
	}(method, sps.Channel)
}

func (sps *SimplePingerService) StartMethods(wg *sync.WaitGroup) {

	for _, method := range sps.methods {
		sps.StartMethod(wg, method)
	}
}

func Start(wg *sync.WaitGroup) {

	commonChannel := aggregator.GetInstance().Channel

	for _, serviceConfig := range configuration.GetInstance().ServiceConfigs {
		wg.Add(1)

		sampleService := new(SimplePingerService)

		sampleService.Configure(serviceConfig, commonChannel)
		sampleService.StartMethods(wg)
	}
	fmt.Println("Services started")
}
