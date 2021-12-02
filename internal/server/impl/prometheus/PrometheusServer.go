package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusServer struct {
	Port int
}

func (ps PrometheusServer) Serve(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Handle("/metrics", promhttp.Handler())

		portString := fmt.Sprintf(":%d", ps.Port)
		http.ListenAndServe(portString, nil)
	}()
	fmt.Println("Prometheus Server started")
}
