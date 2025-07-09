# Supermarket Checkout System

A Go implementation of a supermarket checkout system with support for bulk discounts.

See [Problem Specification](ProblemSpecification.md) for detailed requirements.

## Quick Start

```bash
# Run all tests
go test ./...

# Run the demo
go run main.go
```

## Demo

The demo shows the key scenario from the problem specification:
- Scanning items B, A, B in any order
- Applying the bulk discount (2 Bs for 45 instead of 60)
- Final total: 95 (45 + 50)