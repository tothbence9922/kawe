package main

import (
	"fmt"
	"sync"

	aggregator "github.com/tothbence9922/kawe/internal/aggregator"
	"github.com/tothbence9922/kawe/internal/configuration"
	simpleResult "github.com/tothbence9922/kawe/internal/ping/simple/result"
	simple "github.com/tothbence9922/kawe/internal/ping/simple/service"
)

var (
	wg sync.WaitGroup
)

func startServices() {

	commonChannel := aggregator.GetInstance().Channel

	for _, serviceConfig := range configuration.GetInstance().ServiceConfigs {
		wg.Add(1)

		sampleService := new(simple.SimplePingerService)

		sampleService.Configure(serviceConfig, commonChannel)
		sampleService.StartMethods(wg)
	}
}

func startAggregator() {
	wg.Add(1)
	go func(inChannel <-chan (simpleResult.PingResult)) {
		defer wg.Done()

		for true {
			output := <-inChannel
			fmt.Println(output.String())
		}
	}(aggregator.GetInstance().Channel)

}

func main() {

	configuration.GetInstance()

	startServices()
	startAggregator()

	wg.Wait()
}
