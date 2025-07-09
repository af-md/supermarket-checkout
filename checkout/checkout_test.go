package checkout

import "testing"

func TestScanItem(t *testing.T) {
	checkout := NewCheckout()

	item := item{SKU: "A"}

	err := checkout.Scan(item.SKU)

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


// Negatvie edge cases tests for Scan method
func TestScanEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		sku      string
		errorMsg string
	}{
		{
			name:     "empty SKU",
			sku:      "",
			errorMsg: "SKU cannot be empty",
		},
		{
			name:     "whitespace only SKU",
			sku:      " ",
			errorMsg: "SKU cannot have whitespace",
		},
		// further edge cases can be added here but in the interest of time I will keep it simple
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkout := NewCheckout()
			
			err := checkout.Scan(tt.sku)
			
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
			} else if err.Error() != tt.errorMsg {
				t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
			}
		})
	}
}
