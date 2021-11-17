package simple

type IPingService interface {
	String() string
	configure()
	startMethods()
}
