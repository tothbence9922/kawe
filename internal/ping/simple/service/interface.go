package simple

type PingerService interface {
	String() string
	configure()
	startMethods()
}
