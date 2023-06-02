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

type LabelQueryResponse struct {
	Results   map[string](processorInterfaces.IProcessedData)
	Available bool
}

func getLabelQueriedServices(labelParam string, aggregatorResults map[string]processorInterfaces.IProcessedData) (*LabelQueryResponse, error) {

	results := make(map[string](processorInterfaces.IProcessedData))

	available := true
	for service, result := range aggregatorResults {
		if labelParam == result.GetServiceLabel() {
			results[service] = result

			if result.GetAvailability() == false {
				available = false

			}
		}
	}

	return &LabelQueryResponse{Results: results, Available: available}, nil
}

func getQueriedServices(nameParam string, statusParam string, aggregatorResults map[string]processorInterfaces.IProcessedData) (map[string](processorInterfaces.IProcessedData), error) {

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

	for service, result := range aggregatorResults {
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
	labelQuery := req.URL.Query().Get("label")
	labelQueried := (labelQuery != "")

	aggregatorResults := aggregator.GetInstance().GetResults()

	if len(aggregatorResults) != 0 {
		if labelQueried {
			results, err := getLabelQueriedServices(labelQuery, aggregatorResults)
			if err != nil {
				return []byte(err.Error()), err
			}

			return json.Marshal(results)
		} else {
			results, err := getQueriedServices(nameQuery, statusQuery, aggregatorResults)
			if err != nil {
				return []byte(err.Error()), err
			}

			return json.Marshal(results)
		}
	} else {
		return json.Marshal(make(map[string]processorInterfaces.IProcessedData))
	}
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
