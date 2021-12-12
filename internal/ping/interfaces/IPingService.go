package ping

type IPingService interface {
	String() string
	Configure()
	StartMethods()
}
