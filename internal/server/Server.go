package server

import (
	"sync"

	"github.com/tothbence9922/kawe/internal/configuration"
	httpServer "github.com/tothbence9922/kawe/internal/server/http"
	prometheusServer "github.com/tothbence9922/kawe/internal/server/prometheus"
)

func Start(wg *sync.WaitGroup) {

	for _, value := range configuration.GetInstance().ServerConfigs {
		switch value.Type {
		case "HTTP":
			httpServer := httpServer.HttpServer{Port: value.Port}
			httpServer.Serve(wg)
			break
		case "PROMETHEUS":
			prometheusServer := prometheusServer.PrometheusServer{Port: value.Port}
			prometheusServer.Serve(wg)
			break
		default:
			break
		}
	}

}
