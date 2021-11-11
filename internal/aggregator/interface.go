package aggregator

import simple "github.com/tothbence9922/kawe/internal/ping/simple/result"

type IAggregator interface {
	GetResults() simple.PingResult
}
