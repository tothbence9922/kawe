package simple

type IPingResponse interface {
	String() string
	Json() string
	GetTarget() string
}
