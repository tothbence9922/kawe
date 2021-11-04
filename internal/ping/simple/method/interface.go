package simple

import simple "github.com/tothbence9922/kawe/internal/ping/simple/response"

type PingerMethod interface {
	String() string
	Ping() (simple.PingResponse, error)
	GetPeriodicity() int
}
