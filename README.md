# gobl.mock

Generate realistic, valid [GOBL](https://gobl.org) invoices for any supported tax regime.

## Features

- Works with **all 23 GOBL tax regimes** out of the box
- Reads tax categories, rates, and currencies dynamically from GOBL definitions
- Valid tax IDs with correct check digits for 20+ countries
- Addon support: FacturaE, XRechnung, CFDI
- Credit notes with proper preceding references per regime
- Simplified invoices
- Deterministic output with seed support
- Every generated invoice passes `gobl validate`

## Installation

```bash
go install github.com/pablocaeg/gobl.mock/cmd/gobl.mock@latest
```

## CLI Usage

```bash
gobl.mock --regime ES                       # Spanish invoice
gobl.mock --regime DE --addon de-xrechnung-v3  # German XRechnung
gobl.mock --regime MX --lines 15            # Mexican invoice, 15 lines
gobl.mock --regime FR --credit              # French credit note
gobl.mock --regime IT --simplified          # Italian simplified
gobl.mock --regime BR --seed 42             # Brazilian, reproducible
gobl.mock --regime PT -o invoice.json       # Portuguese, write to file
```

## Library Usage

```go
env, err := mock.Envelope(
    mock.WithRegime(l10n.ES.Tax()),
    mock.WithLines(5),
    mock.WithSeed(42),
)

inv, err := mock.Invoice(
    mock.WithRegime(l10n.DE.Tax()),
    mock.WithAddon("de-xrechnung-v3"),
    mock.WithCredit(),
)
```

## Supported Regimes

All regimes supported by GOBL work dynamically. Tax rates and currencies are read from the regime definition at runtime.

| Code | Country | Tax ID Generation |
|------|---------|-------------------|
| AR | Argentina | CUIT/CUIL with mod-11 check |
| AT | Austria | USt-IdNr with Luhn+4 check |
| BE | Belgium | Enterprise with mod-97 check |
| BR | Brazil | CNPJ/CPF with dual mod-11 check |
| CA | Canada | No tax ID required |
| CH | Switzerland | UID with mod-11 check |
| CO | Colombia | NIT with prime-weighted mod-11 check |
| DE | Germany | USt-IdNr with ISO 7064 check |
| DK | Denmark | CVR with mod-11 check |
| ES | Spain | NIF/CIF/NIE with mod-23/Luhn check |
| FR | France | TVA with SIREN Luhn + mod-97 check |
| GB | United Kingdom | VAT with weighted sum check |
| GR | Greece | AFM with powers-of-2 check |
| IE | Ireland | VAT with mod-23 letter check |
| IN | India | GSTIN with mod-36 check |
| IT | Italy | Partita IVA with Luhn check |
| MX | Mexico | RFC format validation |
| NL | Netherlands | BTW-id with mod-11 check |
| PL | Poland | NIP with weighted mod-11 check |
| PT | Portugal | NIF with weighted mod-11 check |
| SE | Sweden | Org number with Luhn check |
| SG | Singapore | UEN/GST format validation |
| US | United States | No tax ID required |

## Supported Addons

| Addon | Key | Notes |
|-------|-----|-------|
| FacturaE v3 | `es-facturae-v3` | Spanish B2G format |
| XRechnung v3 | `de-xrechnung-v3` | German B2G (auto-adds contact/payment fields) |
| CFDI v4 | `mx-cfdi-v4` | Mexican e-invoicing (auto-applied for MX) |

## How It Works

1. Reads the regime definition from GOBL (`tax.RegimeDefFor()`) to determine tax categories, rates, and currency
2. Generates tax IDs with correct check digits using per-country algorithms
3. Builds a valid invoice structure with realistic parties, line items, and payment details
4. Wraps in a GOBL envelope which triggers `Calculate()` and `Validate()`
5. Returns a fully calculated, valid envelope

The only hardcoded data is locale content (company names, city names) for ES/DE/MX. All other regimes use a generic English fallback that still produces valid invoices.

## Development

```bash
mage lint     # Run linter
mage test     # Run tests
mage build    # Build binary
mage check    # Full CI pipeline
```

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
