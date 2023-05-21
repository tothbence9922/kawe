package types

import "fmt"

type ProcessorConfiguration struct {
	Type      string
	Threshold int
}

func (pc ProcessorConfiguration) String() string {

	return fmt.Sprintf("\tType\n%s\n", pc.Type)
}
