package simple

import simple "github.com/tothbence9922/kawe/internal/ping/simple/response"

type PingResult interface {
	String() string
	Json() string
	GetServiceName() string
	GetResponse() simple.PingResponse
}
