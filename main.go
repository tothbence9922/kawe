package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type PingResponse interface {
	String() string
}

type SimplePingResponse struct {
	success bool
	error   string
}

func (spr SimplePingResponse) String() string {
	successText := "successful"
	errorText := ""
	if !spr.success {
		successText = "failed"
		errorText = spr.error
	}
	return fmt.Sprintf("SimplePing %s.\t%s\n", successText, errorText)
}

type PingerMethod interface {
	Ping(address interface{}) interface{}
}

type SimplePingerMethod struct {
}

func (simpleMethod SimplePingerMethod) Ping(method string, address string, timeout int) (SimplePingResponse, error) {
	ret := SimplePingResponse{}
	if len(method) == 0 {
		return ret, errors.New("No applicable network options given.")
	} else if address == "" {
		return ret, errors.New("Address to ping is not given.")
	}
	duration, _ := time.ParseDuration(fmt.Sprintf("%dms", timeout))
	if timeout > 0 {
		conn, err := net.DialTimeout(method, address, time.Duration(duration))

		ret.success = (conn != nil)

		if conn != nil {
			conn.Close()
		}

		if err != nil {
			ret.error = err.Error()
		}

		return ret, nil
	} else {
		conn, err := net.Dial(method, address)

		ret.success = (conn != nil)

		if conn != nil {
			conn.Close()
		}

		if err != nil {
			ret.error = err.Error()
		}

		return ret, nil
	}
}

func main() {
	pingerMethod := SimplePingerMethod{}

	pingRespose, error := pingerMethod.Ping("tcp", "google.com:443", 5000)
	if error == nil && pingRespose != (SimplePingResponse{}) {
		fmt.Printf("%s\n", pingRespose.String())
	}
}
