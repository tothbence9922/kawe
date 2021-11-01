package ping

type PingResponse interface {
	String() string
}

type PingerMethod interface {
	String() string
	ping() (PingResponse, error)
	getPeriodicity() int
}

type PingerService interface {
	String() string
	configure()
	startMethods()
}
