package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tothbence9922/kawe/internal/aggregator"
	processorInterfaces "github.com/tothbence9922/kawe/internal/processor/interfaces"
)

type HttpServer struct {
	Port int
}

func getQueriedServices(nameParam string, statusParam string) (map[string](processorInterfaces.IProcessedData), error) {

	ag := aggregator.GetInstance()
	ag.Lock()
	defer ag.Unlock()

	results := make(map[string](processorInterfaces.IProcessedData))

	statusQueried := (statusParam != "")
	var status bool
	var statusErr error
	if statusQueried {
		status, statusErr = strconv.ParseBool(statusParam)

		if statusErr != nil {
			return nil, fmt.Errorf("An error occured while parsing status query param. Please check if you entered it correctly.")
		}
	}

	nameQueried := (nameParam != "")

	for service, result := range ag.Results {
		// Filter for name
		if nameQueried {
			if nameParam == service {

				// Filter for status too
				if statusQueried {
					if status == result.GetAvailability() {
						results[service] = result
					}
					// Only filtered for name, add entry
				} else {
					results[service] = result
				}
			}
			// Do no filter for name
		} else {
			// Filter for status
			if statusQueried {
				if status == result.GetAvailability() {
					results[service] = result
				}
				// Did not filter for anything, add entry
			} else {
				results[service] = result
			}
		}
	}
	return results, nil
}

func handleQueryServices(req *http.Request) ([]byte, error) {

	statusQuery := req.URL.Query().Get("status")
	nameQuery := req.URL.Query().Get("name")

	results, err := getQueriedServices(nameQuery, statusQuery)

	if err != nil {
		return []byte(err.Error()), err
	}

	return json.Marshal(results)
}

func handleGetServices(w http.ResponseWriter, req *http.Request) {

	response, _ := handleQueryServices(req)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Allow", http.MethodGet)

	fmt.Fprintf(w, string(response))
}

func (hs HttpServer) Serve(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		server := new(http.Server)
		server.ReadTimeout = 5 * time.Second
		server.WriteTimeout = 5 * time.Second
		defer wg.Done()
		http.HandleFunc("/api/v1/services", handleGetServices)

		portString := fmt.Sprintf(":%d", hs.Port)
		http.ListenAndServe(portString, nil)
	}()
	fmt.Println("HTTP Server started")
}
