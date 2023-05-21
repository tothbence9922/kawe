package aggregator

import processorInterfaces "github.com/tothbence9922/kawe/internal/processor/interfaces"

type IAggregator interface {
	GetResults() processorInterfaces.IProcessedData
	AddResult(processorInterfaces.IProcessedData)
	ClearResults()
}
