package availability

import (
	"encoding/json"

	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type PercentageProcessedData struct {
	Result        pingInterfaces.IPingResult
	Available     bool
	Threshold     int
	Percentage    int
	ProcessorType string
}

func (pd *PercentageProcessedData) GetValue() int {

	return pd.Percentage
}

func (pd *PercentageProcessedData) GetThreshold() int {

	return pd.Threshold
}

func (pd *PercentageProcessedData) GetAvailability() bool {

	return pd.Available
}

func (pd *PercentageProcessedData) GetServiceName() string {

	return pd.Result.GetServiceName()
}
func (pd *PercentageProcessedData) GetProcessorType() string {

	return pd.ProcessorType
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
