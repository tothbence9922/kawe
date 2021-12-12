package processor

type IProcessedData interface {
	String() string
	Json() string
	GetServiceName() string
	GetAvailability() bool
}
