package simple

type PingResponse interface {
	String() string
	Json() string
	GetTarget() string
}
