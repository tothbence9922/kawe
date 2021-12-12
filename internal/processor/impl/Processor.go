package processor

import (
	configurations "github.com/tothbence9922/kawe/internal/configuration"
	availability "github.com/tothbence9922/kawe/internal/processor/impl/availability"
	processorInterfaces "github.com/tothbence9922/kawe/internal/processor/interfaces"
)

func GetProcessor(pc configurations.ProcessorConfiguration) processorInterfaces.IProcessor {

	switch pc.Type {
	case "AVAILABILITY":
		cfg := pc.Params.(AvailabiltiyParams)
		ret := new(availability.AvailabilityProcessor)
		ret.Percentage = cfg.Percentage
		return ret
	default:
		ret := new(availability.AvailabilityProcessor)
		ret.Percentage = 100
		return ret
	}
}
