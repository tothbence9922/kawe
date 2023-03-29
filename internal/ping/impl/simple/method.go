package simple

import (
	"errors"
	"fmt"
	"net"
	"time"

	interfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type PingMethod struct {
	Target      string
	Labels      map[string]string
	Annotations map[string]string
	Method      string
	Timeout     int
	Periodicity int
}

func (spm PingMethod) GetTarget() string {

	return spm.Target
}
func (spm PingMethod) GetLabels() map[string]string {

	return spm.Labels
}

func (spm PingMethod) String() string {

	return fmt.Sprintf("Target\t\tMethod\t\tTimeout\t\tPeriodicity\n%s\t\t%s\t\t%d\t\t%d\n", spm.Target, spm.Method, spm.Timeout, spm.Periodicity)
}

func (spm PingMethod) GetPeriodicity() (Periodicity int) {

	return spm.Periodicity
}

func (spm PingMethod) Ping() (interfaces.IPingResponse, error) {

	ret := PingResponse{Target: spm.Target, Labels: spm.Labels, Annotations: spm.Annotations}

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
