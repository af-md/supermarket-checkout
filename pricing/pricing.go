package pricing

type IPricingService interface {
	GetPricingScheme() (PricingScheme, error)
}

type PricedItem struct {
	Price             int
	DiscountThreshold int
	DiscountPrice     int
	DiscountEnabled   bool
}

type PricingScheme struct {
	Items map[string]PricedItem
}

type PricingService struct{}

func NewPricingService() *PricingService {
	return &PricingService{}
}

func (p *PricingService) GetPricingScheme() (PricingScheme, error) {
	return PricingScheme{
		Items: map[string]PricedItem{
			"A": {Price: 50, DiscountThreshold: 3, DiscountPrice: 130, DiscountEnabled: true},
			"B": {Price: 30, DiscountThreshold: 2, DiscountPrice: 45, DiscountEnabled: true},
			"C": {Price: 20},
			"D": {Price: 15},
		},
	}, nil
}
