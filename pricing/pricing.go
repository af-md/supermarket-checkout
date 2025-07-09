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
