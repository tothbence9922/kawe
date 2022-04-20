package availability

import (
	"encoding/json"

	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type UnitProcessedData struct {
	Result         pingInterfaces.IPingResult
	Available      bool
	Threshold      float32
	AvailableUnits float32
}

func (pd *UnitProcessedData) GetValue() float32 {

	return pd.AvailableUnits
}

func (pd *UnitProcessedData) GetThreshold() float32 {

	return pd.Threshold
}

func (pd *UnitProcessedData) GetAvailability() bool {

	return pd.Available
}

func (pd *UnitProcessedData) GetServiceName() string {

	return pd.Result.GetServiceName()
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
