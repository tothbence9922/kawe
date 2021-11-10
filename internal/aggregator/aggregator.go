package aggregator

import simpleResult "github.com/tothbence9922/kawe/internal/ping/simple/result"

type Aggregator struct {
	Channel chan (simpleResult.PingResult)
	Results [](simpleResult.PingResult)
}

var aggregatorInstance *Aggregator

func GetInstance() *Aggregator {

	if aggregatorInstance == nil {
		aggregatorInstance = new(Aggregator)
		aggregatorInstance.Channel = make(chan simpleResult.PingResult)

	}
	return aggregatorInstance
}

func (a Aggregator) GetResults() [](simpleResult.PingResult) {
	return a.Results
}
