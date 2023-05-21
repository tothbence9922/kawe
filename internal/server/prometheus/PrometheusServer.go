package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tothbence9922/kawe/internal/aggregator"
)

type PrometheusServer struct {
	Port    int
	Metrics map[string]prometheus.Gauge
}

func (ps *PrometheusServer) Init() {
	ps.Metrics = make(map[string](prometheus.Gauge))
}

func (ps *PrometheusServer) CalcMetrics() {
	ag := aggregator.GetInstance()
	ag.Lock()
	defer ag.Unlock()
	processedData := ag.GetResults()
	if processedData != nil {
		for key, value := range processedData {
			if ps.Metrics[key] == nil {
				ps.Metrics[key] = promauto.NewGauge(
					prometheus.GaugeOpts{
						Name: key,
						Help: "Availability of the given service represented by 0 - False - Unavailable and 1 - True - Available values.",
					},
				)
			}
			if value.GetAvailability() {
				ps.Metrics[key].Set(1)
			} else {
				ps.Metrics[key].Set(0)
			}
		}
	}
}

func (ps *PrometheusServer) RecordMetrics() {

	go func() {
		for {
			ps.CalcMetrics()
			time.Sleep(2 * time.Second)
		}
	}()
}

func (ps PrometheusServer) Serve(wg *sync.WaitGroup) {

	wg.Add(1)

	go func() {
		defer wg.Done()

		ps.Init()

		ps.RecordMetrics()

		http.Handle("/metrics", promhttp.Handler())

		portString := fmt.Sprintf(":%d", ps.Port)
		http.ListenAndServe(portString, nil)
	}()
	fmt.Println("Prometheus Server started")
}
