package availability

import (
	"encoding/json"

	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type AvailabilityProcessedData struct {
	Result     pingInterfaces.IPingResult
	Available  bool
	Percentage float32
}

func (pd *AvailabilityProcessedData) GetServiceName() string {

	return pd.Result.GetServiceName()
}

func (pd *AvailabilityProcessedData) String() string {

	ret := ""

	ret += pd.Result.String() + "\n"

	if pd.Available {
		ret += "Service available"
	} else {
		ret += "Service unavailable"
	}

	return ret
}

func (pd *AvailabilityProcessedData) Json() string {

	jsonString, _ := json.Marshal(pd)
	return string(jsonString)
}
