package ping

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/tothbence9922/kawe/kawe/internal/configuration"
)

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
	return fmt.Sprintf("SimplePing %s.\t\t%s\n", successText, errorText)
}

type SimplePingerMethod struct {
	target      string
	method      string
	timeout     int
	periodicity int
}

func (spm SimplePingerMethod) String() string {
	return fmt.Sprintf("target\t\tmethod\t\ttimeout\t\tperiodicity\n%s\t\t%s\t\t%d\t\t%d\n", spm.target, spm.method, spm.timeout, spm.periodicity)
}

func (spm SimplePingerMethod) getPeriodicity() (periodicity int) {
	return spm.periodicity
}

func (spm SimplePingerMethod) ping() (PingResponse, error) {
	ret := SimplePingResponse{}
	if len(spm.method) == 0 {
		return ret, errors.New("No applicable network options given.")
	} else if spm.target == "" {
		return ret, errors.New("Address to ping is not given.")
	}
	duration, _ := time.ParseDuration(fmt.Sprintf("%dms", spm.timeout))
	if spm.timeout > 0 {
		conn, err := net.DialTimeout(spm.method, spm.target, time.Duration(duration))

		ret.success = (conn != nil)

		if conn != nil {
			conn.Close()
		}

		if err != nil {
			ret.error = err.Error()
		}

		return ret, nil
	} else {
		conn, err := net.Dial(spm.method, spm.target)

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

type SimplePingerService struct {
	methods []PingerMethod
}

func (sps SimplePingerService) String() string {
	var ret string

	for _, method := range sps.methods {
		ret = ret + method.String()
	}
	return ret
}

func (sps *SimplePingerService) Configure() {
	for _, pingConfig := range configuration.GetInstance().PingConfigs {
		sps.methods = append(sps.methods, SimplePingerMethod{target: pingConfig.Target, timeout: 5000, method: "tcp", periodicity: pingConfig.Periodicity})
	}
}

func StartMethod(currentMethod PingerMethod) {

}

func (sps *SimplePingerService) StartMethods() {
	var wg sync.WaitGroup

	for _, method := range sps.methods {
		wg.Add(1)
		go func(currentMethod PingerMethod) {
			defer wg.Done()
			for true {
				pingResponse, error := currentMethod.ping() // error should not be ignored
				if error == nil && pingResponse != (SimplePingResponse{}) {
					fmt.Printf("%s\n", pingResponse.String())
				}
				time.Sleep(time.Second * time.Duration(currentMethod.getPeriodicity()))
			}
		}(method)
	}
	wg.Wait()
}
