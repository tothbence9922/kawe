package tcp

import (
	"fmt"
	"reflect"
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
	KillChannel chan (bool)
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
	sps.KillChannel = make(chan bool)

	sps.Result = &PingResult{ServiceName: sps.Name, Annotations: config.Annotations, ProcessorType: config.ProcessorConfig.Type, Responses: make(map[string](interfaces.IPingResponse))}
	for _, pingConfig := range config.Pods {
		sps.methods = append(sps.methods, PingMethod{Target: fmt.Sprintf("%s:%s", pingConfig.Address, pingConfig.Port), Name: pingConfig.Name, Timeout: pingConfig.Timeout, Periodicity: pingConfig.Periodicity})
	}
}

func (sps *PingService) StartMethod(wg *sync.WaitGroup, method interfaces.IPingMethod) {
	wg.Add(1)
	go func(method interfaces.IPingMethod, processor processorInterfaces.IProcessor) {
		defer wg.Done()
		for {
			select {
			case <-sps.KillChannel:
				fmt.Println("Killing ", sps.Name)
				return
			default:
				pingResponse, error := method.Ping()
				if error == nil {
					sps.Lock()
					sps.Result.AddResponse(pingResponse)
					processor.ProcessData(sps.Result)
					sps.Unlock()
				}
				time.Sleep(time.Second * time.Duration(method.GetPeriodicity()))
				break
			}
		}
	}(method, sps.Processor)
}

func (sps *PingService) StartMethods(wg *sync.WaitGroup) {
	for _, currentMethod := range sps.methods {
		go func(method interfaces.IPingMethod) {
			sps.StartMethod(wg, method)
		}(currentMethod)
	}
}

func killServices(services []*PingService) {
	fmt.Println("Killing services")

	for _, currentService := range services {
		go func(service *PingService) {
			service.KillChannel <- true
		}(currentService)
	}
}

func Start(wg *sync.WaitGroup) {
	wg.Add(1)

	var lastConfig []configTypes.NamespaceConfiguration

	var services []*PingService

	fmt.Println("Starting Services")
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		for true {
			fmt.Println("Looking for changes")

			newConfig := configuration.GetInstance().EndpointConfigs.Namespaces

			if !reflect.DeepEqual(lastConfig, newConfig) {
				fmt.Println("Changes found")
				lastConfig = newConfig
				killServices(services)

				for _, namespaceConfig := range newConfig {
					for _, currentServiceConfig := range namespaceConfig.Services {
						go func(serviceConfig configTypes.ServiceConfiguration) {
							service := new(PingService)
							curProcessor := processor.GetProcessor(serviceConfig.ProcessorConfig)
							service.Configure(serviceConfig, curProcessor)

							services = append(services, service)

							service.StartMethods(wg)
						}(currentServiceConfig)
					}
				}
			}
			time.Sleep(time.Second * time.Duration(10))
		}
	}(wg)

	fmt.Println("Processors started")
	fmt.Println("Services started")
}
