# gobl.mock

Generate realistic, valid [GOBL](https://gobl.org) invoices for testing and development.

## Features

- Generates invoices that pass `gobl build` validation
- Supports multiple tax regimes: Spain (ES), Germany (DE), Mexico (MX)
- Addon support: FacturaE, XRechnung, CFDI
- Valid tax IDs with correct check digits
- Locale-appropriate company names, addresses, and products
- Credit notes with proper preceding references
- Simplified invoices
- Deterministic output with seed support

## Installation

```bash
go install github.com/pablocaeg/gobl.mock/cmd/gobl.mock@latest
```

## CLI Usage

```bash
# Spanish invoice
gobl.mock --regime ES

# German XRechnung invoice
gobl.mock --regime DE --addon de-xrechnung-v3

# Mexican invoice with 15 line items
gobl.mock --regime MX --lines 15

# Spanish credit note
gobl.mock --regime ES --credit

# Simplified invoice
gobl.mock --regime ES --simplified

# Reproducible output
gobl.mock --regime ES --seed 42

# Write to file
gobl.mock --regime ES --output invoice.json
```

## Library Usage

```go
package main

import (
	"fmt"

	mock "github.com/pablocaeg/gobl.mock"
	"github.com/invopop/gobl/l10n"
)

func main() {
	env, err := mock.Envelope(
		mock.WithRegime(l10n.ES.Tax()),
		mock.WithLines(5),
		mock.WithSeed(42),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated invoice in envelope %s\n", env.Head.UUID)
}
```

## Supported Regimes

| Country | Code | Default Addon | Credit Notes | Simplified |
|---------|------|---------------|-------------|------------|
| Spain   | ES   | -             | ✓           | ✓          |
| Germany | DE   | -             | ✓           | ✓          |
| Mexico  | MX   | mx-cfdi-v4    | ✓           | ✓          |

## Supported Addons

| Addon | Key | Country |
|-------|-----|---------|
| FacturaE v3 | `es-facturae-v3` | ES |
| XRechnung v3 | `de-xrechnung-v3` | DE |
| CFDI v4 | `mx-cfdi-v4` | MX |

## Development

This project uses [Mage](https://magefile.org/) for build automation:

```bash
mage lint     # Run linter
mage test     # Run tests
mage testrace # Run tests with race detector
mage build    # Build binary
mage install  # Install binary
mage check    # Full CI pipeline
```

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
