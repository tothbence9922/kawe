package availability

import (
	"encoding/json"

	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type UnitProcessedData struct {
	Result         pingInterfaces.IPingResult
	Available      bool
	Threshold      int
	AvailableUnits int
	ProcessorType  string
}

func (pd *UnitProcessedData) GetValue() int {

	return pd.AvailableUnits
}

func (pd *UnitProcessedData) GetThreshold() int {

	return pd.Threshold
}

func (pd *UnitProcessedData) GetAvailability() bool {

	return pd.Available
}

func (pd *UnitProcessedData) GetProcessorType() string {

	return pd.ProcessorType
}

func (pd *UnitProcessedData) GetServiceName() string {

	return pd.Result.GetServiceName()
}

func (pd *UnitProcessedData) GetServiceLabel() string {

	return pd.Result.GetServiceLabel()
}

func (pd *UnitProcessedData) String() string {

	ret := ""

	ret += pd.Result.String() + "\n"

	if pd.Available {
		ret += "Service available"
	} else {
		ret += "Service unavailable"
	}

	return ret
}

func (pd *UnitProcessedData) Json() string {

	jsonString, _ := json.Marshal(pd)
	return string(jsonString)
}
