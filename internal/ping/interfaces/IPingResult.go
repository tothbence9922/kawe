package simple

type IPingResult interface {
	String() string
	Json() string
	GetServiceName() string
	GetResponses() map[string](IPingResponse)
}
