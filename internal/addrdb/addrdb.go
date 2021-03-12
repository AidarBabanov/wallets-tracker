package addrdb

type Address struct {
	Address string `csv:"address"`
}

type AddressDatabase interface {
	IsEmpty() bool
	Next() (Address, error)
	Index() int
	Len() int
	All() []Address
	Add(address Address) error
}
