package main

import (
	"sync"

	aggregator "github.com/tothbence9922/kawe/internal/aggregator"
	simpleService "github.com/tothbence9922/kawe/internal/ping/impl/simple/service"
	server "github.com/tothbence9922/kawe/internal/server"
)

var (
	wg sync.WaitGroup
)

func main() {
	wgPtr := &wg

	aggregator.Start(wgPtr)

	simpleService.Start(wgPtr)

	httpServer := server.HttpServer{Port: 8080}
	httpServer.Serve(wgPtr)

	wg.Wait()
}
