package checkout

type ICheckout interface {
	Scan(SKU string) (err error)
	GetTotalPrice() (totalPrice int, err error)
}

type checkout struct {
	items []string
}

type item struct {
	SKU string
}

func NewCheckout() *checkout {
	return &checkout{
		items: make([]string, 0),
	}
}

func (c *checkout) Scan(SKU string) error {
	c.items = append(c.items, SKU)
	return nil
}

func (c *checkout) GetTotalPrice() (int, error) {
	panic("not implemented")
}
