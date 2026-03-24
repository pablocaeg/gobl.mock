# gobl.mock

<img src="https://github.com/invopop/gobl/blob/main/gobl_logo_black_rgb.svg#gh-light-mode-only" width="80" height="97" alt="GOBL Logo">
<img src="https://github.com/invopop/gobl/blob/main/gobl_logo_white_rgb.svg#gh-dark-mode-only" width="80" height="97" alt="GOBL Logo">

Generate realistic, valid [GOBL](https://gobl.org) invoices for any supported tax regime.

Released under the Apache 2.0 [LICENSE](LICENSE).

[![Lint](https://github.com/pablocaeg/gobl.mock/actions/workflows/lint.yaml/badge.svg)](https://github.com/pablocaeg/gobl.mock/actions/workflows/lint.yaml)
[![Test](https://github.com/pablocaeg/gobl.mock/actions/workflows/test.yaml/badge.svg)](https://github.com/pablocaeg/gobl.mock/actions/workflows/test.yaml)
[![GoDoc](https://godoc.org/github.com/pablocaeg/gobl.mock?status.svg)](https://godoc.org/github.com/pablocaeg/gobl.mock)

## Introduction

gobl.mock is a Go library and CLI tool that generates valid GOBL invoices for testing and development. It reads tax regime definitions dynamically from GOBL, generates tax IDs with correct check digits, and produces invoices that pass the full GOBL validation pipeline.

All 24 GOBL tax regimes and all 20 addons are supported. Every generated invoice passes `gobl validate`.

## Use Cases

- **Testing GOBL conversion tools** (gobl.cfdi, gobl.facturae, gobl.xinvoice, gobl.html) — generate valid input invoices for any country instead of maintaining static YAML fixtures.
- **Integration testing** — any system that processes GOBL invoices can use `mock.Invoice()` to generate test data programmatically.
- **Regression testing GOBL itself** — generate thousands of invoices across all regimes and verify the regime system works correctly after changes.
- **Exploring GOBL** — see what a valid invoice looks like for any country and addon combination.
- **Load testing** — generate diverse invoice data at scale with different seeds for stress testing invoice processing pipelines.
- **Templated testing** — provide a partial invoice (specific customer, specific line items) and let gobl.mock fill in the rest with valid data.

## Installation

```bash
go install github.com/pablocaeg/gobl.mock/cmd/gobl.mock@latest
```

## Usage

### CLI

```bash
gobl.mock --regime ES                          # Spanish invoice
gobl.mock --regime DE --addon de-xrechnung-v3  # German XRechnung
gobl.mock --regime MX --lines 15               # Mexican invoice, 15 lines
gobl.mock --regime FR --type credit-note       # French credit note
gobl.mock --regime ES --type corrective        # Spanish corrective invoice
gobl.mock --regime IT --type proforma          # Italian proforma
gobl.mock --regime IT --simplified             # Italian simplified
gobl.mock --regime BR --addon br-nfe-v4        # Brazilian NF-e
gobl.mock --regime PT --addon pt-saft-v1       # Portuguese SAF-T
gobl.mock --regime ES --seed 42                # Reproducible output
gobl.mock --regime ES -o invoice.json          # Write to file
```

### Library

```go
// Generate a valid GOBL envelope for any regime.
env, err := mock.Envelope(
    mock.WithRegime(l10n.ES.Tax()),
    mock.WithAddon("es-facturae-v3"),
    mock.WithLines(5),
    mock.WithSeed(42),
)

// Extract the calculated invoice directly.
inv, err := mock.Invoice(
    mock.WithRegime(l10n.DE.Tax()),
    mock.WithType(bill.InvoiceTypeCreditNote),
)

// Use a template to test a specific scenario.
inv, err := mock.Invoice(
    mock.WithRegime(l10n.ES.Tax()),
    mock.WithTemplate(&bill.Invoice{
        Customer: mySpecificCustomer,
        Lines:    mySpecificLines,
    }),
)
```

## Supported Regimes

Tax categories, rates, and currencies are read dynamically from GOBL regime definitions. Tax IDs are generated with correct check digits.

| Code | Country | Check Digit Algorithm |
|------|---------|----------------------|
| AE | United Arab Emirates | Format validation |
| AR | Argentina | Mod-11 |
| AT | Austria | Luhn+4 |
| BE | Belgium | Mod-97 |
| BR | Brazil | Dual Mod-11 (CNPJ/CPF) |
| CA | Canada | No code required |
| CH | Switzerland | Mod-11 |
| CO | Colombia | Prime-weighted Mod-11 |
| DE | Germany | ISO 7064 |
| DK | Denmark | Mod-11 |
| ES | Spain | Mod-23 / Luhn (NIF/CIF/NIE) |
| FR | France | SIREN Luhn + Mod-97 |
| GB | United Kingdom | Weighted sum |
| GR | Greece | Powers-of-2 |
| IE | Ireland | Mod-23 letter check |
| IN | India | Mod-36 (GSTIN) |
| IT | Italy | Luhn (Partita IVA) |
| MX | Mexico | Format validation (RFC) |
| NL | Netherlands | Mod-11 |
| PL | Poland | Weighted Mod-11 |
| PT | Portugal | Weighted Mod-11 |
| SE | Sweden | Luhn |
| SG | Singapore | Format validation (UEN) |
| US | United States | No code required |

## Supported Addons

| Key | Country | Description |
|-----|---------|-------------|
| `ar-arca-v4` | AR | Argentina ARCA electronic invoicing |
| `br-nfe-v4` | BR | Brazilian NF-e (Nota Fiscal Eletronica) |
| `br-nfse-v1` | BR | Brazilian NFS-e (Nota Fiscal de Servicos) |
| `co-dian-v2` | CO | Colombian DIAN tax authority format |
| `de-xrechnung-v3` | DE | German XRechnung B2G e-invoicing |
| `de-zugferd-v2` | DE | German ZUGFeRD hybrid invoice |
| `es-facturae-v3` | ES | Spanish FacturaE B2G format |
| `es-sii-v1` | ES | Spanish SII real-time reporting |
| `es-tbai-v1` | ES | Basque Country TicketBAI |
| `es-verifactu-v1` | ES | Spanish VeriFactu anti-fraud |
| `eu-en16931-v2017` | EU | European EN16931 standard (use with any EU regime) |
| `fr-choruspro-v1` | FR | French Chorus Pro B2G |
| `fr-ctc-flow2-v1` | FR | French B2B CTC |
| `fr-facturx-v1` | FR | French Factur-X hybrid |
| `gr-mydata-v1` | GR | Greek MyData tax reporting |
| `it-sdi-v1` | IT | Italian SDI electronic invoicing |
| `it-ticket-v1` | IT | Italian fiscal receipts |
| `mx-cfdi-v4` | MX | Mexican CFDI (auto-applied for MX) |
| `pl-favat-v1` | PL | Polish FA_VAT / KSeF |
| `pt-saft-v1` | PT | Portuguese SAF-T |

## How It Works

1. Reads the regime definition from GOBL (`tax.RegimeDefFor()`) to determine tax categories, rates, and currency.
2. Generates tax IDs with correct check digits using per-country algorithms that mirror GOBL's own validators.
3. Applies addon-specific extensions, identities, and structural requirements.
4. Wraps in a GOBL envelope which triggers `Calculate()` — GOBL computes all tax totals.
5. GOBL runs `Validate()` — if the envelope builds, the invoice is guaranteed valid.

Tax ID correctness is verified by generating each ID and running it through GOBL's own validation pipeline. The generators do not self-validate — they are validated against the same code that validates real invoices.

## Development

This project uses [Mage](https://magefile.org/) for build automation:

```bash
mage lint     # Run golangci-lint
mage test     # Run tests
mage testrace # Run tests with race detector
mage build    # Build binary
mage install  # Install binary
mage check    # Full CI pipeline (lint + test)
```

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
