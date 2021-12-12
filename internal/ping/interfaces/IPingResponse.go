package ping

import "time"

type IPingResponse interface {
	String() string
	Json() string
	GetTarget() string
	GetSuccess() bool
	GetError() string
	GetTimestamp() time.Time
}
