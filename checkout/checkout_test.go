package checkout

import (
	"errors"
	"supermarket-checkout/pricing"
	"testing"
)

type MockPricingService struct{}

func (m *MockPricingService) GetPricingScheme() (pricing.PricingScheme, error) {

	return pricing.PricingScheme{
		Items: map[string]pricing.PricedItem{
			"A": {Price: 50, DiscountThreshold: 3, DiscountPrice: 130},
			"B": {Price: 30, DiscountThreshold: 2, DiscountPrice: 45},
		},
	}, nil
}

// pricing service that simulates an error
type MockPricingServiceError struct{}

func (m *MockPricingServiceError) GetPricingScheme() (pricing.PricingScheme, error) {
	return pricing.PricingScheme{}, errors.New("pricing service error")
}

func TestScanItem(t *testing.T) {
	checkout := NewCheckout(&MockPricingService{})

	err := checkout.Scan("A")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(checkout.items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(checkout.items))
	}

	if checkout.items["A"] != 1 {
		t.Errorf("Expected item 'A' value to be 1, got %d", checkout.items["A"])
	}

}

// Negative edge cases tests for Scan method
func TestScan_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		sku      string
		errorMsg string
	}{
		{
			name:     "empty SKU",
			sku:      "",
			errorMsg: EmptySKUError,
		},
		{
			name:     "whitespace only SKU",
			sku:      " ",
			errorMsg: WhitespaceSKUError,
		},
		// further edge cases can be added here but in the interest of time I will keep it simple
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkout := NewCheckout(&MockPricingService{})

			err := checkout.Scan(tt.sku)

			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
			} else if err.Error() != tt.errorMsg {
				t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
			}
		})
	}
}

func TestGetTotalPrice_SingleItem(t *testing.T) {
	checkout := NewCheckout(&MockPricingService{})

	err := checkout.Scan("A")
	if err != nil {
		t.Fatalf("Failed to scan item: %v", err)
	}

	total, err := checkout.GetTotalPrice()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := 50
	if total != expected {
		t.Errorf("Expected total %d, got %d", expected, total)
	}
}

func TestGetTotalPrice_EmptyCheckout(t *testing.T) {
	checkout := NewCheckout(&MockPricingService{})

	total, err := checkout.GetTotalPrice()
	if err == nil {
		t.Error("Expected error for empty checkout, got nil")
	}

	if total != 0 {
		t.Errorf("Expected total 0 for empty checkout, got %d", total)
	}
}

// no discounts
func TestGetTotalPrice_MultipleItems(t *testing.T) {
	tests := []struct {
		name     string
		items    []string
		expected int
	}{
		{
			name:     "multiple item A",
			items:    []string{"A", "A", "A"},
			expected: 150,
		},
		{
			name:     "multiple random items",
			items:    []string{"A", "B", "C", "A", "B", "A", "D"},
			expected: 245,
		},
		{
			name:     "multiple item D",
			items:    []string{"D", "D", "D", "D"},
			expected: 60,
		},
		{
			name:     "multiple item C",
			items:    []string{"C", "C", "C"},
			expected: 60,
		},
		{
			name:     "multiple B items",
			items:    []string{"B", "B", "B", "B"},
			expected: 120,
		},
		{
			name:     "many mixed items",
			items:    []string{"A", "B", "C", "D", "A", "B", "C", "D", "A", "D", "C", "B"},
			expected: 345,
		},
		{
			name:     "unknown item",
			items:    []string{"X"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkout := NewCheckout(&MockPricingService{})

			for _, item := range tt.items {
				err := checkout.Scan(item)
				if err != nil {
					t.Fatalf("Failed to scan item %s: %v", item, err)
				}
			}

			total, err := checkout.GetTotalPrice()
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if total != tt.expected {
				t.Errorf("Expected total %d, got %d", tt.expected, total)
			}
		})
	}
}

func TestGetTotalPrice_ErrorInPricingService(t *testing.T) {
	// shoould be table test
	tests := []struct {
		name               string
		items              []string
		pricingServiceMock pricing.IPricingService
	}{
		{
			name:               "error in pricing service",
			items:              []string{"A", "A", "A"},
			pricingServiceMock: &MockPricingServiceError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkout := NewCheckout(tt.pricingServiceMock)

			for _, item := range tt.items {
				err := checkout.Scan(item)
				if err != nil {
					t.Fatalf("Failed to scan item %s: %v", item, err)
				}
			}

			total, err := checkout.GetTotalPrice()
			if err == nil {
				t.Error("Expected error from pricing service, got nil")
			}

			if total != 0 {
				t.Errorf("Expected total 0 due to pricing service error, got %d", total)
			}
		})
	}
}

func TestGetTotalPrice_MultipleItemsDiscountApplied(t *testing.T) {

	tests := []struct {
		name               string
		items              []string
		pricingServiceMock pricing.IPricingService
		expected           int
	}{
		{
			name:               "2 Bs for 45",
			items:              []string{"B", "B"},
			pricingServiceMock: &MockPricingService{},
			expected:           45,
		},
		{
			name:               "2 Bs for 45 plus 1 B for 30",
			items:              []string{"B", "B", "B"},
			pricingServiceMock: &MockPricingService{},
			expected:           75,
		},
		{
			name:               "11 Bs for 120",
			items:              []string{"B", "B", "B", "B", "B", "B", "B", "B", "B", "B", "B"},
			pricingServiceMock: &MockPricingService{},
			expected:           255,
		},
		{
			name:               "3 As for 130",
			items:              []string{"A", "A", "A"},
			pricingServiceMock: &MockPricingService{},
			expected:           130,
		},
		{
			name:               "3 As for 130 + 1 A for 50",
			items:              []string{"A", "A", "A", "A"},
			pricingServiceMock: &MockPricingService{},
			expected:           180,
		},
		{
			name:               "5 As, 5Bs",
			items:              []string{"A", "A", "A", "A", "A", "B", "B", "B", "B", "B"},
			pricingServiceMock: &MockPricingService{},
			expected:           350,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkout := NewCheckout(tt.pricingServiceMock)

			for _, item := range tt.items {
				err := checkout.Scan(item)
				if err != nil {
					t.Fatalf("Failed to scan item %s: %v", item, err)
				}
			}

			total, err := checkout.GetTotalPrice()
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if total != tt.expected {
				t.Errorf("Expected total %d, got %d", tt.expected, total)
			}
		})
	}

}
