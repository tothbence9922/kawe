package simple

import (
	"fmt"
	"sync"
	"time"

	"github.com/tothbence9922/kawe/internal/aggregator"
	"github.com/tothbence9922/kawe/internal/configuration"
	interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type PingService struct {
	methods []interfaces.IPingMethod
	Name    string
	Channel chan (interfaces.IPingResult)
	Result  interfaces.IPingResult
}

func (sps PingService) String() string {

	var ret string

	for _, method := range sps.methods {
		ret = ret + method.String()
	}
	return ret
}

func (sps *PingService) Configure(config configuration.ServiceConfiguration, channel chan (interfaces.IPingResult)) {

	sps.Name = config.Name
	sps.Channel = channel
	sps.Result = &PingResult{ServiceName: sps.Name, Responses: make(map[string](interfaces.IPingResponse))}
	for _, pingConfig := range config.PingConfigs {
		sps.methods = append(sps.methods, PingMethod{Target: pingConfig.Target, Timeout: 5000, Method: "tcp", Periodicity: pingConfig.Periodicity})
	}
}

func (sps *PingService) StartMethod(wg *sync.WaitGroup, method interfaces.IPingMethod) {

	go func(method interfaces.IPingMethod, outChannel chan<- (interfaces.IPingResult)) {

		defer wg.Done()
		for true {
			pingResponse, error := method.Ping()
			if error == nil {
				sps.Result.AddResponse(pingResponse)
				outChannel <- sps.Result
			}
			time.Sleep(time.Second * time.Duration(method.GetPeriodicity()))
		}
	}(method, sps.Channel)
}

func (sps *PingService) StartMethods(wg *sync.WaitGroup) {

	for _, method := range sps.methods {
		sps.StartMethod(wg, method)
	}
}

func Start(wg *sync.WaitGroup) {

	commonChannel := aggregator.GetInstance().Channel

	for _, serviceConfig := range configuration.GetInstance().ServiceConfigs {
		wg.Add(1)

		sampleService := new(PingService)

		sampleService.Configure(serviceConfig, commonChannel)
		sampleService.StartMethods(wg)
	}
	fmt.Println("Services started")
}
