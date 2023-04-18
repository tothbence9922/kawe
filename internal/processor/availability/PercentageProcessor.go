package availability

import (
	"math"

	"github.com/tothbence9922/kawe/internal/aggregator"
	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type PercentageProcessor struct {
	Threshold int
}

func (ap *PercentageProcessor) ProcessData(result pingInterfaces.IPingResult) {

	processedData := new(PercentageProcessedData)
	processedData.Result = result
	processedData.Threshold = ap.Threshold
	processedData.ProcessorType = "PERCENTAGE"

	responses := result.GetResponses()
	totalCount := 0
	successCount := 0

	for _, value := range responses {
		totalCount++
		if value.GetSuccess() {
			successCount++
		}
	}

	percentage := float64(((float64(successCount) / float64(totalCount)) * float64(100)))
	processedData.Percentage = int(int64(math.Floor(percentage)))

	processedData.Available = false
	if processedData.Threshold <= processedData.Percentage {
		processedData.Available = true
	}

	commonChannel := aggregator.GetInstance().Channel

	commonChannel <- processedData

}
