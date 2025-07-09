package checkout

type ICheckout interface {
	Scan(SKU string) (err error)
	GetTotalPrice() (totalPrice int, err error)
}

type Checkout struct {
	items []string
}

func NewCheckout() *Checkout {
	return &Checkout{
		items: make([]string, 0),
	}
}

func (c *Checkout) Scan(SKU string) error {
	panic("not implemented")
}

func (c *Checkout) GetTotalPrice() (int, error) {
	panic("not implemented")
}
