package ping

type IPingMethod interface {
	String() string
	Ping() (IPingResponse, error)
	GetPeriodicity() int
	GetTarget() string
}
