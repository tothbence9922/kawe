package types

import "fmt"

type PodConfiguration struct {
	Name        string
	Address     string
	Port        string
	Enabled     bool
	Periodicity int
	Timeout     int
}

func (p PodConfiguration) String() string {
	enabledString := "DISABLED"
	if p.Enabled {
		enabledString = "ENABLED"
	}
	ret := p.Name + " " + p.Address + " " + enabledString + " " + fmt.Sprint(p.Periodicity) + " " + fmt.Sprint(p.Timeout)
	return ret
}
