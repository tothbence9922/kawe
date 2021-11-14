package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/tothbence9922/kawe/internal/aggregator"
)

type HttpServer struct {
	Port int
}

func api(w http.ResponseWriter, req *http.Request) {
	outJson, _ := json.Marshal(aggregator.GetInstance().Results)
	fmt.Fprintf(w, string(outJson))
}

func (hs HttpServer) Serve(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.HandleFunc("/api", api)

		portString := fmt.Sprintf(":%d", hs.Port)
		http.ListenAndServe(portString, nil)
	}()
	fmt.Println("HTTP Server started")
}
