package main

import (
	simple "github.com/tothbence9922/kawe/internal/ping/simple/service"
)

func main() {

	sampleService := new(simple.SimplePingerService)

	sampleService.Configure()

	sampleService.StartMethods()
}
