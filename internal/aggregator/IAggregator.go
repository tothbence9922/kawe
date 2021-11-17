package aggregator

import interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"

type IAggregator interface {
	GetResults() interfaces.IPingResult
	AddResult(interfaces.IPingResult)
}
