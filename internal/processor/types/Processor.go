package processor

import (
	configTypes "github.com/tothbence9922/kawe/internal/configuration/types"
	availability "github.com/tothbence9922/kawe/internal/processor/availability"
	processorInterfaces "github.com/tothbence9922/kawe/internal/processor/interfaces"
)

func GetProcessor(pc configTypes.ProcessorConfiguration) processorInterfaces.IProcessor {

	switch {
	case pc.Type == "PERCENTAGE":
		ret := new(availability.PercentageProcessor)
		ret.Threshold = pc.Threshold
		return ret
	case pc.Type == "UNIT":
		ret := new(availability.UnitProcessor)
		ret.Threshold = pc.Threshold
		return ret
	default:
		ret := new(availability.PercentageProcessor)
		ret.Threshold = 100
		return ret
	}
}
