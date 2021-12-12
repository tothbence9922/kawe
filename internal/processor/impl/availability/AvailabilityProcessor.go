package availability

import (
	"github.com/tothbence9922/kawe/internal/aggregator"
	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type AvailabilityProcessor struct {
	Percentage float32
}

func (ap *AvailabilityProcessor) ProcessData(result pingInterfaces.IPingResult) {

	processedData := new(AvailabilityProcessedData)
	processedData.Result = result

	responses := result.GetResponses()
	totalCount := 0
	successCount := 0

	for _, value := range responses {
		totalCount++
		if value.GetSuccess() {
			successCount++
		}
	}

	percentage := float32(successCount) / float32(totalCount) * float32(100)
	processedData.Percentage = percentage

	processedData.Available = false
	if ap.Percentage <= percentage {
		processedData.Available = true
	}

	commonChannel := aggregator.GetInstance().Channel

	commonChannel <- processedData

}
