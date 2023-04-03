package processor

import (
	configTypes "github.com/tothbence9922/kawe/internal/configuration/types"
	availability "github.com/tothbence9922/kawe/internal/processor/impl/availability"
	processorInterfaces "github.com/tothbence9922/kawe/internal/processor/interfaces"
)

func GetProcessor(pc configTypes.ProcessorConfiguration) processorInterfaces.IProcessor {

	switch {
	case pc.Type == "PERCENTAGE":
		cfg := pc.Params.(map[string]interface{})
		ret := new(availability.PercentageProcessor)
		ret.Threshold = float32(cfg["Threshold"].(float64))
		return ret
	case pc.Type == "UNIT":
		cfg := pc.Params.(map[string]interface{})
		ret := new(availability.UnitProcessor)
		ret.Threshold = float32(cfg["Threshold"].(float64))
		return ret
	default:
		ret := new(availability.PercentageProcessor)
		ret.Threshold = 100
		return ret
	}
}
