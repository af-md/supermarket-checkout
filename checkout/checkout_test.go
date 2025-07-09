package checkout

import (
	"testing"
)

type MockPricingService struct{}

func (m *MockPricingService) GetPrice(sku string) (int, error) {
	if sku == "A" {
		return 50, nil
	}
	return 0, nil
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

	if checkout.items[0] != "A" {
		t.Errorf("Expected item 'A', got %s", checkout.items[0])
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
