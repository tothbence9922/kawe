package processor

type IProcessedData interface {
	String() string
	Json() string
	GetServiceName() string
	GetAvailability() bool
	GetThreshold() float32
	GetValue() float32
}
