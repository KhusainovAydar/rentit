package flat

const (
	MOSCOW uint8 = 1
	SPB
)

type Flat struct {
	URL       string
	Rooms     uint8
	Address   string
	Lat       float64
	Lon       float64
	Price     uint64
	Area      uint64
	FromOwner bool
}
