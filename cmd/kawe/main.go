package main

import (
	"sync"

	aggregator "github.com/tothbence9922/kawe/internal/aggregator"
	"github.com/tothbence9922/kawe/internal/configuration"
	service "github.com/tothbence9922/kawe/internal/ping/simple/service"
)

var (
	wg sync.WaitGroup
)

func main() {
	wgPtr := &wg

	configuration.GetInstance()

	aggregator.Start(wgPtr)
	service.Start(wgPtr)

	wg.Wait()
}
