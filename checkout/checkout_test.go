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
