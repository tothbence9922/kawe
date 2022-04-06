package availability

import (
	"github.com/tothbence9922/kawe/internal/aggregator"
	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type UnitProcessor struct {
	Threshold float32
}

func (ap *UnitProcessor) ProcessData(result pingInterfaces.IPingResult) {

	processedData := new(UnitProcessedData)
	processedData.Result = result
	processedData.Threshold = float32(ap.Threshold)

	responses := result.GetResponses()
	totalCount := 0
	successCount := float32(0)

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
