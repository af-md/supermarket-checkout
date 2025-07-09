package main

import (
	"fmt"
	"supermarket-checkout/checkout"
	"supermarket-checkout/pricing"
)

func main() {
	fmt.Println("=== Supermarket Checkout Demo ===")
	fmt.Println("Testing mixed discounted items scenario:")
	fmt.Println()

	// Create  services
	pricingService := pricing.NewPricingService()
	checkout := checkout.NewCheckout(pricingService)

	// Scan items: B, A, B (as mentioned in problem specification)
	fmt.Println("Scanning items: B, A, B")
	
	if err := checkout.Scan("B"); err != nil {
		fmt.Printf("Error scanning B: %v\n", err)
		return
	}
	
	if err := checkout.Scan("A"); err != nil {
		fmt.Printf("Error scanning A: %v\n", err)
		return
	}
	
	if err := checkout.Scan("B"); err != nil {
		fmt.Printf("Error scanning B: %v\n", err)
		return
	}

	total, err := checkout.GetTotalPrice()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Total: %d\n", total)
	fmt.Printf("Expected: 95 (2 Bs for 45 + 1 A for 50)\n")
	
	if total == 95 {
		fmt.Println("✅ PASS - Discount applied correctly!")
	} else {
		fmt.Println("❌ FAIL")
	}
}