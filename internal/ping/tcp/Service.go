package tcp

import (
	"fmt"
	"sync"
	"time"

	"github.com/tothbence9922/kawe/internal/configuration"
	configTypes "github.com/tothbence9922/kawe/internal/configuration/types"
	interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
	processorInterfaces "github.com/tothbence9922/kawe/internal/processor/interfaces"
	processor "github.com/tothbence9922/kawe/internal/processor/types"
)

type PingService struct {
	sync.Mutex

	methods     []interfaces.IPingMethod
	Name        string
	Annotations map[string]string
	Processor   processorInterfaces.IProcessor
	Result      interfaces.IPingResult
}

func (sps *PingService) GetMethods() []interfaces.IPingMethod {

	return sps.methods
}
func (sps *PingService) String() string {

	var ret string

	for _, method := range sps.methods {
		ret = ret + method.String()
	}
	return ret
}

func (sps *PingService) Configure(config configTypes.ServiceConfiguration, processor processorInterfaces.IProcessor) {

	sps.Name = config.Name
	sps.Annotations = config.Annotations
	sps.Processor = processor
	sps.Result = &PingResult{ServiceName: sps.Name, Annotations: config.Annotations, Responses: make(map[string](interfaces.IPingResponse))}
	for _, pingConfig := range config.Pods {
		sps.methods = append(sps.methods, PingMethod{Target: fmt.Sprintf("%s:%s", pingConfig.Address, pingConfig.Port), Name: pingConfig.Name, Timeout: pingConfig.Timeout, Periodicity: pingConfig.Periodicity})
	}
}

func (sps *PingService) StartMethod(wg *sync.WaitGroup, method interfaces.IPingMethod) {
	wg.Add(1)
	go func(method interfaces.IPingMethod, processor processorInterfaces.IProcessor) {
		defer wg.Done()
		for true {
			pingResponse, error := method.Ping()
			if error == nil {
				sps.Lock()
				sps.Result.AddResponse(pingResponse)
				processor.ProcessData(sps.Result)
				sps.Unlock()
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
	}
	fmt.Println("Processors started")
	fmt.Println("Services started")
}
