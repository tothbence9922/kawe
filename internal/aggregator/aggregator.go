package aggregator

import (
	"fmt"
	"sync"

	processorInterfaces "github.com/tothbence9922/kawe/internal/processor/interfaces"
)

type Aggregator struct {
	sync.Mutex
	Channel chan (processorInterfaces.IProcessedData)
	Results map[string](processorInterfaces.IProcessedData)
}

var aggregatorInstance *Aggregator

func GetInstance() *Aggregator {

	if aggregatorInstance == nil {
		aggregatorInstance = new(Aggregator)
		aggregatorInstance.Channel = make(chan processorInterfaces.IProcessedData)
		aggregatorInstance.Results = make(map[string](processorInterfaces.IProcessedData))
	}

	return aggregatorInstance
}

func (a *Aggregator) GetResults() map[string](processorInterfaces.IProcessedData) {

	return a.Results
}

func (a *Aggregator) AddResult(newResult processorInterfaces.IProcessedData) {

	a.Lock()
	defer a.Unlock()

	GetInstance().Results[newResult.GetServiceName()] = newResult
}

func Start(wg *sync.WaitGroup) {

	wg.Add(1)
	go func(inChannel <-chan (processorInterfaces.IProcessedData)) {
		defer wg.Done()

		for true {
			newResult := <-inChannel
			GetInstance().AddResult(newResult)
		}
	}(GetInstance().Channel)
	fmt.Println("Aggregator started")
}
