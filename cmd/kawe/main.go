package main

import (
	"fmt"
	"sync"

	"github.com/tothbence9922/kawe/internal/configuration"
	simple "github.com/tothbence9922/kawe/internal/ping/simple/service"
)

var (
	wg sync.WaitGroup
)

func startServices() {

	for _, serviceConfig := range configuration.GetInstance().ServiceConfigs {
		wg.Add(1)

		fmt.Println("ASd")
		sampleService := new(simple.SimplePingerService)

		sampleService.Configure(serviceConfig)
		sampleService.StartMethods(wg)
	}
	wg.Wait()
}

func main() {

	startServices()

}
