package main

import (
	"fmt"

	"github.com/tothbence9922/kawe/internal/configuration"
	simple "github.com/tothbence9922/kawe/internal/ping/simple/service"
)

func main() {
	configuration.GetInstance()

	samplePingConfig := configuration.PingConfiguration{Periodicity: 5, Target: "google.com:443"}
	samplePingConfigs := make([]configuration.PingConfiguration, 1)
	samplePingConfigs[0] = samplePingConfig

	sampleService := new(simple.SimplePingerService)

	sampleService.Configure()

	fmt.Println(sampleService.String())

	sampleService.StartMethods()
}
