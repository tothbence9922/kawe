package aggregator

import (
	"fmt"
	"sync"

	simpleResponse "github.com/tothbence9922/kawe/internal/ping/simple/response"
	simpleResult "github.com/tothbence9922/kawe/internal/ping/simple/result"
)

type Aggregator struct {
	Channel chan (simpleResult.PingResult)
	Results map[string](simpleResponse.PingResponse)
}

var aggregatorInstance *Aggregator

func GetInstance() *Aggregator {

	if aggregatorInstance == nil {
		aggregatorInstance = new(Aggregator)
		aggregatorInstance.Channel = make(chan simpleResult.PingResult)
		aggregatorInstance.Results = make(map[string]simpleResponse.PingResponse)
	}

	return aggregatorInstance
}

func (a Aggregator) GetResults() map[string](simpleResponse.PingResponse) {

	return a.Results
}

func Start(wg *sync.WaitGroup) {

	wg.Add(1)
	go func(inChannel <-chan (simpleResult.PingResult)) {
		defer wg.Done()

		for true {
			// Handling incoming data, setting the "state"
			output := <-inChannel
			GetInstance().Results[output.GetServiceName()] = output.GetResponse()
			fmt.Println(GetInstance().Results)
		}
	}(GetInstance().Channel)
	fmt.Println("Aggregator started")
}
