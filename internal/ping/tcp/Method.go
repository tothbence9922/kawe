package tcp

import (
	"errors"
	"fmt"
	"net"
	"time"

	interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type PingMethod struct {
	Target      string
	Name        string
	Annotations map[string]string
	Timeout     int
	Periodicity int
}

func (spm PingMethod) GetTarget() string {

	return spm.Target
}

func (spm PingMethod) String() string {

	return fmt.Sprintf("Target\t\tTimeout\t\tPeriodicity\n%s\t\t%s\t\t%d\t\t%d\n", spm.Target, spm.Timeout, spm.Periodicity)
}

func (spm PingMethod) GetPeriodicity() (Periodicity int) {

	return spm.Periodicity
}

func (spm PingMethod) Ping() (interfaces.IPingResponse, error) {
	ret := PingResponse{Name: spm.Name, Target: spm.Target}

	if spm.Target == "" {
		return ret, errors.New("Address to ping is not given.")
	}

	duration, _ := time.ParseDuration(fmt.Sprintf("%dms", spm.Timeout))

	if spm.Timeout > 0 {
		conn, err := net.DialTimeout("tcp", spm.Target, time.Duration(duration))
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
		conn, err := net.Dial("tcp", spm.Target)

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
