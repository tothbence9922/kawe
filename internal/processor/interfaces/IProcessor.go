package processor

import (
	pingInterfaces "github.com/tothbence9922/kawe/internal/ping/interfaces"
)

type IProcessor interface {
	ProcessData(pingInterfaces.IPingResult)
}
