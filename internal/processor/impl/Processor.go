package processor

import (
	configurations "github.com/tothbence9922/kawe/internal/configuration"
	availability "github.com/tothbence9922/kawe/internal/processor/impl/availability"
	processorInterfaces "github.com/tothbence9922/kawe/internal/processor/interfaces"
)

func GetProcessor(pc configurations.ProcessorConfiguration) processorInterfaces.IProcessor {

	switch {
	case pc.Type == "AVAILABILITY":
		cfg := pc.Params.(map[string]interface{})
		ret := new(availability.AvailabilityProcessor)
		ret.Threshold = float32(cfg["Threshold"].(float64))
		return ret
	default:
		ret := new(availability.AvailabilityProcessor)
		ret.Threshold = 100
		return ret
	}
}
