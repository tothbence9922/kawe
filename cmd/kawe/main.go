package main

import (
	"fmt"

	"github.com/tothbence9922/kawe/kawe/internal/configuration"
	"github.com/tothbence9922/kawe/kawe/internal/ping"
)

func main() {
	configuration.GetInstance()

	samplePingConfig := configuration.PingConfiguration{Periodicity: 5, Target: "google.com:443"}
	samplePingConfigs := make([]configuration.PingConfiguration, 1)
	samplePingConfigs[0] = samplePingConfig

	sampleService := new(ping.SimplePingerService)

	sampleService.Configure()

	fmt.Println(sampleService.String())

	sampleService.StartMethods()
}
