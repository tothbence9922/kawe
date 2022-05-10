package aggregator

import (
	"fmt"
	"strings"
	"sync"

	processorInterfaces "github.com/tothbence9922/kawe/internal/processor/interfaces"
)

type Aggregator struct {
	sync.Mutex
	Channel chan (processorInterfaces.IProcessedData)
	Results map[string](processorInterfaces.IProcessedData)
}

var (
	aggregatorInstance *Aggregator
	initLock           sync.Mutex
)

func GetInstance() *Aggregator {
	initLock.Lock()
	defer initLock.Unlock()

	if aggregatorInstance == nil {
		aggregatorInstance = new(Aggregator)

		aggregatorInstance.Channel = make(chan processorInterfaces.IProcessedData)
		aggregatorInstance.Results = make(map[string](processorInterfaces.IProcessedData))
	}

	return aggregatorInstance
}

func (a *Aggregator) GetResults() map[string](processorInterfaces.IProcessedData) {

	a.Lock()
	defer a.Unlock()
	m := make(map[string](processorInterfaces.IProcessedData), len(a.Results))
	for k, v := range a.Results {
		m[k] = v
	}
	return m
}

func (a *Aggregator) AddResult(newResult processorInterfaces.IProcessedData) {
	a.Lock()
	defer a.Unlock()

	newResultMetricName := strings.ReplaceAll(newResult.GetServiceName(), "-", "_")

	m := make(map[string](processorInterfaces.IProcessedData), len(a.Results))
	for k, v := range a.Results {
		m[k] = v
	}

	m[newResultMetricName] = newResult

	a.Results = m
}

func Start(wg *sync.WaitGroup) {

	wg.Add(1)
	go func() {
		defer wg.Done()

		for true {
			aggregator := GetInstance()
			newResult := <-aggregator.Channel
			aggregator.AddResult(newResult)
		}
	}()
	fmt.Println("Aggregator started")
}
