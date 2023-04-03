package availability

import (
	"encoding/json"

	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type PercentageProcessedData struct {
	Result     pingInterfaces.IPingResult
	Available  bool
	Threshold  float32
	Percentage float32
}

func (pd *PercentageProcessedData) GetValue() float32 {

	return pd.Percentage
}

func (pd *PercentageProcessedData) GetThreshold() float32 {

	return pd.Threshold
}

func (pd *PercentageProcessedData) GetAvailability() bool {

	return pd.Available
}

func (pd *PercentageProcessedData) GetServiceName() string {

	return pd.Result.GetServiceName()
}

func (pd *PercentageProcessedData) String() string {

	ret := ""

	ret += pd.Result.String() + "\n"

	if pd.Available {
		ret += "Service available"
	} else {
		ret += "Service unavailable"
	}

	return ret
}

func (pd *PercentageProcessedData) Json() string {

	jsonString, _ := json.Marshal(pd)
	return string(jsonString)
}
