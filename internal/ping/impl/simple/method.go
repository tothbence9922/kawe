package simple

import (
	"errors"
	"fmt"
	"net"
	"time"

	interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type SimplePingerMethod struct {
	Target      string
	Method      string
	Timeout     int
	Periodicity int
}

func (spm SimplePingerMethod) GetTarget() string {
	return spm.Target
}

func (spm SimplePingerMethod) String() string {

	return fmt.Sprintf("Target\t\tMethod\t\tTimeout\t\tPeriodicity\n%s\t\t%s\t\t%d\t\t%d\n", spm.Target, spm.Method, spm.Timeout, spm.Periodicity)
}

func (spm SimplePingerMethod) GetPeriodicity() (Periodicity int) {

	return spm.Periodicity
}

func (spm SimplePingerMethod) Ping() (interfaces.IPingResponse, error) {

	ret := PingResponse{Target: spm.Target}

	if len(spm.Method) == 0 {
		return ret, errors.New("No applicable network options given.")
	} else if spm.Target == "" {
		return ret, errors.New("Address to ping is not given.")
	}

	duration, _ := time.ParseDuration(fmt.Sprintf("%dms", spm.Timeout))

	if spm.Timeout > 0 {
		conn, err := net.DialTimeout(spm.Method, spm.Target, time.Duration(duration))

		ret.Timestamp = time.Now()
		ret.Success = (conn != nil)

		if conn != nil {
			conn.Close()
		}

		if err != nil {
			ret.Error = err.Error()
		}

		return ret, nil
	} else {
		conn, err := net.Dial(spm.Method, spm.Target)

		ret.Timestamp = time.Now()
		ret.Success = (conn != nil)

		if conn != nil {
			conn.Close()
		}

		if err != nil {
			ret.Error = err.Error()
		}

		return ret, nil
	}
}
