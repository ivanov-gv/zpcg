package callback

import (
	"time"
)

type Type string

const (
	UnsupportedType  Type = ""
	UpdateType       Type = "0"
	ReverseRouteType Type = "1"
)

type Callback struct {
	Type             Type
	UpdateData       UpdateData
	ReverseRouteData ReverseRouteData
}

type UpdateData struct {
	Origin, Destination string
	Date                time.Time
}

type ReverseRouteData struct {
	Origin, Destination string
}
