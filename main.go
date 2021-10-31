package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
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
	return fmt.Sprintf("SimplePing %s.\t\t%s\n", successText, errorText)
}

type PingerMethod interface {
	ping() (SimplePingResponse, error)
	getPeriodicity() int
	String() string
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

func (spm SimplePingerMethod) ping() (SimplePingResponse, error) {
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

type PingConfiguration struct {
	periodicity int
	method      PingerMethod
	target      string
}

func (pc PingConfiguration) String() string {
	return fmt.Sprintf("periodicity\t\ttarget\n%d\t\t%s\n", pc.periodicity, pc.target)
}

type Configuration struct {
	pingConfig []PingConfiguration
}

func (c Configuration) String() string {
	var ret string

	for _, curPingConfig := range c.pingConfig {
		ret = ret + curPingConfig.String()
	}

	return ret

}

type PingerService interface {
	configure(config Configuration)
	String() string
	startMethods()
}

type SimplePingerService struct {
	methods []PingerMethod
	config  Configuration
}

func (sps SimplePingerService) String() string {
	var ret string

	for _, method := range sps.methods {
		ret = ret + method.String()
	}

	ret += sps.config.String()
	return ret
}

func (sps *SimplePingerService) configure() {
	for _, pingConfig := range sps.config.pingConfig {
		sps.methods = append(sps.methods, SimplePingerMethod{target: pingConfig.target, timeout: 5000, method: "tcp", periodicity: pingConfig.periodicity})
	}
}

func startMethod(currentMethod PingerMethod) {

}

func (sps *SimplePingerService) startMethods() {
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
				time.Sleep(1000 * time.Duration(currentMethod.getPeriodicity()))
			}
		}(method)
	}
	wg.Wait()
}

func main() {

	samplePingConfig := PingConfiguration{periodicity: 10, target: "google.com:443"}
	samplePingConfigs := make([]PingConfiguration, 1)
	samplePingConfigs[0] = samplePingConfig

	sampleConfig := Configuration{pingConfig: samplePingConfigs}

	sampleService := new(SimplePingerService)
	sampleService.config = sampleConfig

	sampleService.configure()

	fmt.Println(sampleService.String())

	sampleService.startMethods()
}
