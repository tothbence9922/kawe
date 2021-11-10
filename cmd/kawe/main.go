package main

import (
	"github.com/tothbence9922/kawe/internal/configuration"
	simple "github.com/tothbence9922/kawe/internal/ping/simple/service"
)

func startServices() {

	for _, serviceConfig := range configuration.GetInstance().ServiceConfigs {
		sampleService := new(simple.SimplePingerService)

		sampleService.Configure(serviceConfig)

		sampleService.StartMethods()
	}

}

func main() {

	startServices()

}
