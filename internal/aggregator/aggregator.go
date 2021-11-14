package aggregator

import (
	"fmt"
	"sync"

	simpleResult "github.com/tothbence9922/kawe/internal/ping/simple/result"
)

type Aggregator struct {
	Channel chan (simpleResult.PingResult)
	Results map[string](simpleResult.PingResult)
}

var aggregatorInstance *Aggregator

func GetInstance() *Aggregator {

	if aggregatorInstance == nil {
		aggregatorInstance = new(Aggregator)
		aggregatorInstance.Channel = make(chan simpleResult.PingResult)
		aggregatorInstance.Results = make(map[string](simpleResult.PingResult))
	}

	return aggregatorInstance
}

func (a Aggregator) GetResults() map[string](simpleResult.PingResult) {

	return a.Results
}

func Start(wg *sync.WaitGroup) {

	wg.Add(1)
	go func(inChannel <-chan (simpleResult.PingResult)) {
		defer wg.Done()

		for true {
			// Handling incoming data, setting the "state"
			newResult := <-inChannel
			GetInstance().Results[newResult.GetServiceName()] = newResult

			//outJson, _ := json.Marshal(GetInstance().Results) // Printing the state for debug...
			//fmt.Println(string(outJson))
		}
	}(GetInstance().Channel)
	fmt.Println("Aggregator started")
}
