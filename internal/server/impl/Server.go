package server

import (
	"sync"

	//"github.com/tothbence9922/kawe/internal/configuration"
	prometheusServer "github.com/tothbence9922/kawe/internal/server/impl/prometheus"
	httpServer "github.com/tothbence9922/kawe/internal/server/impl/simple"
)

func Start(wg *sync.WaitGroup) {

	//for _, value := range configuration.GetInstance().ServerConfigs {
	//	switch value.Type {
	//	case "HTTP":
	httpServer := httpServer.HttpServer{Port: /*value.Port*/ 80}
	httpServer.Serve(wg)
	//		break
	//	case "PROMETHEUS":
	prometheusServer := prometheusServer.PrometheusServer{Port: /*value.Port*/ 80}
	prometheusServer.Serve(wg)
	//		break
	//	default:
	//		httpServer := httpServer.HttpServer{Port: 8080}
	//		httpServer.Serve(wg)
	//	}
	//}

}
