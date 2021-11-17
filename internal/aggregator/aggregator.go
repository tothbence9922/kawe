package aggregator

import (
	"fmt"
	"sync"

	interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type Aggregator struct {
	sync.Mutex
	Channel chan (interfaces.IPingResult)
	Results map[string](interfaces.IPingResult)
}

var aggregatorInstance *Aggregator

func GetInstance() *Aggregator {

	if aggregatorInstance == nil {
		aggregatorInstance = new(Aggregator)
		aggregatorInstance.Channel = make(chan interfaces.IPingResult)
		aggregatorInstance.Results = make(map[string](interfaces.IPingResult))
	}

	return aggregatorInstance
}

func (a *Aggregator) GetResults() map[string](interfaces.IPingResult) {

	return a.Results
}

func Start(wg *sync.WaitGroup) {

	wg.Add(1)
	go func(inChannel <-chan (interfaces.IPingResult)) {
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
