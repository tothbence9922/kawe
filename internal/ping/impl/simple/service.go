package simple

import (
	"fmt"
	"sync"
	"time"

	"github.com/tothbence9922/kawe/internal/configuration"
	interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
	processor "github.com/tothbence9922/kawe/internal/processor/impl"
	processorInterfaces "github.com/tothbence9922/kawe/internal/processor/interfaces"
)

type PingService struct {
	methods   []interfaces.IPingMethod
	Name      string
	Processor processorInterfaces.IProcessor
	Result    interfaces.IPingResult
}

func (sps PingService) String() string {

	var ret string

	for _, method := range sps.methods {
		ret = ret + method.String()
	}
	return ret
}

func (sps *PingService) Configure(config configuration.ServiceConfiguration, processor processorInterfaces.IProcessor) {

	sps.Name = config.Name
	sps.Processor = processor
	sps.Result = &PingResult{ServiceName: sps.Name, Responses: make(map[string](interfaces.IPingResponse))}
	for _, pingConfig := range config.Pods {
		sps.methods = append(sps.methods, PingMethod{Target: pingConfig.Address + ":80", Timeout: pingConfig.Timeout, Method: "tcp", Periodicity: pingConfig.Periodicity})
	}
}

func (sps *PingService) StartMethod(wg *sync.WaitGroup, method interfaces.IPingMethod) {
	wg.Add(1)
	go func(method interfaces.IPingMethod, processor processorInterfaces.IProcessor) {
		defer wg.Done()
		for true {
			pingResponse, error := method.Ping()
			if error == nil {
				sps.Result.AddResponse(pingResponse)
				processor.ProcessData(sps.Result)
			}
			time.Sleep(time.Second * time.Duration(method.GetPeriodicity()))
		}
	}(method, sps.Processor)
}

func (sps *PingService) StartMethods(wg *sync.WaitGroup) {

	for _, method := range sps.methods {
		sps.StartMethod(wg, method)
	}
}

func Start(wg *sync.WaitGroup) {

	for _, namespaceConfig := range configuration.GetInstance().EndpointConfigs.Namespaces {
		for _, serviceConfig := range namespaceConfig.Services {

			service := new(PingService)
			curProcessor := processor.GetProcessor(serviceConfig.ProcessorConfig)

			service.Configure(serviceConfig, curProcessor)
			service.StartMethods(wg)
		}
		// TODO ping pods without services too
	}
	fmt.Println("Processors started")
	fmt.Println("Services started")
}
