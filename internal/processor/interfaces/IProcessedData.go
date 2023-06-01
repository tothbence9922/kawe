package processor

type IProcessedData interface {
	String() string
	Json() string
	GetServiceName() string
	GetServiceLabel() string
	GetAvailability() bool
	GetProcessorType() string
	GetThreshold() int
	GetValue() int
}
