package callback

type Type string

const (
	UnsupportedType  Type = ""
	UpdateType       Type = "0"
	ReverseRouteType Type = "1"
)

type Callback struct {
	Type Type
	Data Data
}

type Data struct {
	Origin, Destination string
}
