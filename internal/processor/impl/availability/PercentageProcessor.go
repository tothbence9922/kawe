package availability

import (
	"github.com/tothbence9922/kawe/internal/aggregator"
	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type PercentageProcessor struct {
	Threshold float32
}

func (ap *PercentageProcessor) ProcessData(result pingInterfaces.IPingResult) {

	processedData := new(PercentageProcessedData)
	processedData.Result = result
	processedData.Threshold = float32(ap.Threshold)

	responses := result.GetResponses()
	totalCount := 0
	successCount := 0

	for _, value := range responses {
		totalCount++
		if value.GetSuccess() {
			successCount++
		}
	}

	percentage := ((float32(successCount) / float32(totalCount)) * float32(100))
	processedData.Percentage = percentage

	processedData.Available = false
	if processedData.Threshold <= processedData.Percentage {
		processedData.Available = true
	}

	commonChannel := aggregator.GetInstance().Channel

	commonChannel <- processedData

}
