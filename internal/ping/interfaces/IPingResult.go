package ping

type IPingResult interface {
	String() string
	Json() string
	GetServiceName() string
	GetServiceLabel() string
	GetAnnotations() map[string]string
	GetResponses() map[string](IPingResponse)
	AddResponse(IPingResponse)
}
