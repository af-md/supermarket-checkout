package pricing


type IPricingService interface {
	GetPrice(Sku string) (int, error)
}