package main

import (
	"sync"

	aggregator "github.com/tothbence9922/kawe/internal/aggregator"
	simple "github.com/tothbence9922/kawe/internal/ping/impl/simple"
	prometheusServer "github.com/tothbence9922/kawe/internal/server/impl/prometheus"
	httpServer "github.com/tothbence9922/kawe/internal/server/impl/simple"
)

var (
	wg sync.WaitGroup
)

func main() {
	wgPtr := &wg

	aggregator.Start(wgPtr)

	simple.Start(wgPtr)

	httpServer := httpServer.HttpServer{Port: 8080}
	httpServer.Serve(wgPtr)

	prometheusServer := prometheusServer.PrometheusServer{Port: 8080}
	prometheusServer.Serve(wgPtr)

	wg.Wait()
}
