# gobl.mock

Generate realistic, valid [GOBL](https://gobl.org) invoices for any supported tax regime.

Released under the Apache 2.0 [LICENSE](https://github.com/pablocaeg/gobl.mock/blob/main/LICENSE), Copyright 2025 Pablo Carrasco.

[![Lint](https://github.com/pablocaeg/gobl.mock/actions/workflows/lint.yaml/badge.svg)](https://github.com/pablocaeg/gobl.mock/actions/workflows/lint.yaml)
[![Test](https://github.com/pablocaeg/gobl.mock/actions/workflows/test.yaml/badge.svg)](https://github.com/pablocaeg/gobl.mock/actions/workflows/test.yaml)
[![GoDoc](https://godoc.org/github.com/pablocaeg/gobl.mock?status.svg)](https://godoc.org/github.com/pablocaeg/gobl.mock)

## Introduction

gobl.mock generates valid GOBL invoices for testing and development. It supports all 23 GOBL tax regimes and all 17 addons out of the box. Every generated invoice passes `gobl validate`.

Key features:

- Read tax categories, rates, and currencies dynamically from GOBL regime definitions.
- Generate valid tax IDs with correct check digits for 20+ countries.
- Support all GOBL addons with the required extensions, identities, and structural fields.
- Produce credit notes, simplified invoices, and configurable line counts.
- Provide deterministic output with seed support for reproducible testing.

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
gobl.mock --regime FR --credit                 # French credit note
gobl.mock --regime IT --simplified             # Italian simplified
gobl.mock --regime BR --addon br-nfe-v4        # Brazilian NF-e
gobl.mock --regime PT --addon pt-saft-v1       # Portuguese SAF-T
gobl.mock --regime ES --seed 42                # Reproducible output
gobl.mock --regime ES -o invoice.json          # Write to file
```

### Library

```go
env, err := mock.Envelope(
    mock.WithRegime(l10n.ES.Tax()),
    mock.WithAddon("es-facturae-v3"),
    mock.WithLines(5),
    mock.WithSeed(42),
)

inv, err := mock.Invoice(
    mock.WithRegime(l10n.DE.Tax()),
    mock.WithCredit(),
)
```

## Supported Regimes

All regimes generate valid tax IDs with correct check digits.

| Code | Country | Check Digit Algorithm |
|------|---------|----------------------|
| AR | Argentina | Mod-11 |
| AT | Austria | Luhn+4 |
| BE | Belgium | Mod-97 |
| BR | Brazil | Dual Mod-11 |
| CA | Canada | - |
| CH | Switzerland | Mod-11 |
| CO | Colombia | Prime-weighted Mod-11 |
| DE | Germany | ISO 7064 |
| DK | Denmark | Mod-11 |
| ES | Spain | Mod-23 / Luhn |
| FR | France | SIREN Luhn + Mod-97 |
| GB | United Kingdom | Weighted sum |
| GR | Greece | Powers-of-2 |
| IE | Ireland | Mod-23 letter |
| IN | India | Mod-36 |
| IT | Italy | Luhn |
| MX | Mexico | Format validation |
| NL | Netherlands | Mod-11 |
| PL | Poland | Weighted Mod-11 |
| PT | Portugal | Weighted Mod-11 |
| SE | Sweden | Luhn |
| SG | Singapore | Format validation |
| US | United States | - |

## Supported Addons

| Key | Country | Description |
|-----|---------|-------------|
| `ar-arca-v4` | AR | Argentina ARCA |
| `br-nfe-v4` | BR | Brazilian NF-e |
| `br-nfse-v1` | BR | Brazilian NFS-e |
| `co-dian-v2` | CO | Colombian DIAN |
| `de-xrechnung-v3` | DE | German XRechnung |
| `de-zugferd-v2` | DE | German ZUGFeRD |
| `es-facturae-v3` | ES | Spanish FacturaE |
| `es-sii-v1` | ES | Spanish SII |
| `es-tbai-v1` | ES | Basque TicketBAI |
| `es-verifactu-v1` | ES | Spanish VeriFactu |
| `eu-en16931-v2017` | EU | European EN16931 |
| `fr-choruspro-v1` | FR | French Chorus Pro |
| `fr-facturx-v1` | FR | French Factur-X |
| `gr-mydata-v1` | GR | Greek MyData |
| `it-sdi-v1` | IT | Italian SDI |
| `mx-cfdi-v4` | MX | Mexican CFDI |
| `pl-favat-v1` | PL | Polish FA_VAT |
| `pt-saft-v1` | PT | Portuguese SAF-T |

## How It Works

1. Reads the regime definition from GOBL to determine tax categories, rates, and currency.
2. Generates tax IDs with correct check digits using per-country algorithms.
3. Applies addon-specific extensions, identities, and structural requirements.
4. Wraps in a GOBL envelope which triggers `Calculate()` and `Validate()`.
5. Returns a fully calculated, valid envelope.

## Development

This project uses [Mage](https://magefile.org/) for build automation:

```bash
mage lint     # Run linter
mage test     # Run tests
mage build    # Build binary
mage check    # Full CI pipeline
```

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
