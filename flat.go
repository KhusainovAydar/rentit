package rentit

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

type FlatsRequest struct {
	City       uint8
	Rooms      []uint8
	MinPrice   uint64
	MaxPrice   uint64
	LastUpdate uint64
	FromOwner  bool
}
