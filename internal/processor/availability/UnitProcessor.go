package availability

import (
	"github.com/tothbence9922/kawe/internal/aggregator"
	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type UnitProcessor struct {
	Threshold int
}

func (ap *UnitProcessor) ProcessData(result pingInterfaces.IPingResult) {

	processedData := new(UnitProcessedData)
	processedData.Result = result
	processedData.Threshold = ap.Threshold
	processedData.ProcessorType = "UNIT"
	responses := result.GetResponses()
	totalCount := 0
	successCount := 0

	for _, value := range responses {
		totalCount++
		if value.GetSuccess() {
			successCount++
		}
	}

	processedData.AvailableUnits = successCount

	processedData.Available = false
	if processedData.Threshold <= processedData.AvailableUnits {
		processedData.Available = true
	}

	commonChannel := aggregator.GetInstance().Channel

	commonChannel <- processedData

}
