package checkout

import (
	"errors"
	"fmt"
	"supermarket-checkout/pricing"
)

const (
	// SKU related error messages
	EmptySKUError      = "SKU cannot be empty"
	WhitespaceSKUError = "SKU cannot have whitespace"

	// Generate price related error messages
	NoItemsError      = "no items scanned, checkout empty"
	PriceSchemeError  = "error occurred when getting pricing scheme. Error found %s"
	ItemNotFoundError = "item %s not found in pricing scheme"
)

type ICheckout interface {
	Scan(SKU string) (err error)
	GetTotalPrice() (totalPrice int, err error)
}

type Checkout struct {
	items          map[string]int
	pricingService pricing.IPricingService
}

func NewCheckout(pricingService pricing.IPricingService) *Checkout {
	return &Checkout{
		items:          make(map[string]int),
		pricingService: pricingService,
	}
}

func (c *Checkout) Scan(SKU string) error {
	// could have trimmed the whitespace here, but I wanted to show how I would handle multiple edge cases
	if SKU == "" {
		return errors.New(EmptySKUError)
	}

	if SKU == " " {
		return errors.New(WhitespaceSKUError)
	}

	c.items[SKU] = c.items[SKU] + 1
	return nil
}

func (c *Checkout) GetTotalPrice() (int, error) {
	if len(c.items) == 0 {
		return 0, errors.New(NoItemsError)
	}

	pricingScheme, err := c.pricingService.GetPricingScheme()
	if err != nil {
		return 0, fmt.Errorf(PriceSchemeError, err.Error())
	}

	return c.CalculateTotalPrice(pricingScheme)
}

func (c *Checkout) CalculateTotalPrice(pricingScheme pricing.PricingScheme) (int, error) {
	totalPrice := 0

	for sku, quantity := range c.items {
		pricingRules, ok := pricingScheme.Items[sku]
		if !ok {
			fmt.Printf(ItemNotFoundError, sku)
			continue
		}

		if quantity >= pricingRules.DiscountThreshold && pricingRules.DiscountEnabled {
			discountGroups := quantity / pricingRules.DiscountThreshold
			itemsAfterDiscount := quantity % pricingRules.DiscountThreshold

			totalPrice += discountGroups * pricingRules.DiscountPrice
			totalPrice += itemsAfterDiscount * pricingRules.Price
		} else {
			totalPrice += quantity * pricingRules.Price
		}
	}

	return totalPrice, nil
}
