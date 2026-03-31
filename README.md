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

All 24 GOBL tax regimes and all 20 addons are supported. Domain-specific scenarios generate realistic invoices for hotels, freelancers, restaurants, e-commerce, and cross-border B2B — with correct tax rates, retained taxes, document charges, and service periods.

## Use Cases

- **Testing GOBL conversion tools** (gobl.cfdi, gobl.facturae, gobl.xinvoice, gobl.html) — generate valid input invoices for any country instead of maintaining static YAML fixtures.
- **Integration testing** — any system that processes GOBL invoices can use `mock.Invoice()` to generate test data programmatically.
- **Regression testing GOBL itself** — generate thousands of invoices across all regimes and verify the regime system works correctly after changes.
- **Exploring GOBL** — see what a valid invoice looks like for any country, addon, and business domain.
- **Load testing** — generate diverse invoice data at scale with different seeds for stress testing invoice processing pipelines.
- **Templated testing** — provide a partial invoice (specific customer, specific line items) and let gobl.mock fill in the rest with valid data.

## Installation

```bash
go install github.com/pablocaeg/gobl.mock/cmd/gobl.mock@latest
```

## Usage

### CLI

```bash
# Basic usage
gobl.mock --regime ES                          # Spanish invoice
gobl.mock --regime DE --addon de-xrechnung-v3  # German XRechnung
gobl.mock --regime MX --lines 15               # Mexican invoice, 15 lines
gobl.mock --regime FR --type credit-note       # French credit note
gobl.mock --regime IT --simplified             # Italian simplified
gobl.mock --regime ES --seed 42                # Reproducible output
gobl.mock --regime ES -o invoice.json          # Write to file

# Domain-specific scenarios
gobl.mock --regime IT --scenario hotel          # Hotel with 10% VAT + lodging tax
gobl.mock --regime ES --scenario freelance      # Freelancer with IRPF retention
gobl.mock --regime MX --scenario freelance      # Freelancer with ISR + retained IVA
gobl.mock --regime FR --scenario restaurant     # Restaurant with 10% food VAT
gobl.mock --regime DE --scenario ecommerce      # E-commerce with shipping charge
gobl.mock --regime ES --scenario reverse-charge # Cross-border reverse charge (ES->NL)

# Combine scenarios with addons and invoice types
gobl.mock --regime IT --addon it-sdi-v1 --scenario hotel --type credit-note
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

// Domain-specific scenarios.
env, err := mock.Envelope(
    mock.WithRegime("IT"),
    mock.WithAddon("it-sdi-v1"),
    mock.WithScenario(mock.ScenarioHotel),
    mock.WithSeed(42),
)
// -> Camera Matrimoniale (4 nights, Mar 25-30) @ 10% VAT
//    Tassa di Soggiorno (exempt)
//    SDI extensions auto-applied

inv, err := mock.Invoice(
    mock.WithRegime(l10n.ES.Tax()),
    mock.WithScenario(mock.ScenarioFreelance),
)
// -> Desarrollo de software, 40h (Feb 2026) @ 90.00/h, 10% discount
//    VAT 21% + IRPF 15% retained
//    Supplier: Maria Francisca Montero Esteban
//    Payment due in 30 days

// Use a template to override specific fields.
inv, err = mock.Invoice(
    mock.WithRegime(l10n.ES.Tax()),
    mock.WithTemplate(&bill.Invoice{
        Customer: mySpecificCustomer,
        Lines:    mySpecificLines,
    }),
)
```

## Scenarios

> **Work in progress.** The scenarios feature was proposed by [@juanmoliner](https://github.com/juanmoliner) to generate realistic, domain-specific invoices for testing GOBL conversion tools and exploring real-world invoice patterns. Feedback and additional scenario ideas are welcome.

Scenarios generate domain-specific invoices with realistic line items, tax treatments, service periods, and party structures. They compose with any regime, addon, invoice type, and seed.

| Key | Description | Features |
|-----|-------------|----------|
| `hotel` | Hotel accommodation | Reduced/intermediate VAT for rooms; exempt lodging taxes; stay period on room lines |
| `freelance` | Self-employed professional services | Retained taxes (IRPF, IRPEF, ISR, RR); hourly rates; billing period; line discounts; supplier as person; due-date payment |
| `restaurant` | Restaurant / catering | Intermediate VAT for food (FR, IT); localized menu items; 10% service charge |
| `ecommerce` | E-commerce / retail | Physical goods; shipping as delivery charge; localized product names |
| `reverse-charge` | Cross-border B2B | Reverse-charge VAT tag; foreign customer; legal note auto-applied by GOBL |

### Hotel Regime Overrides

Localized room names and the correct reduced/intermediate VAT rate per country. Countries not listed use the standard rate as fallback.

| Country | Room Rate | Source |
|---------|-----------|--------|
| IT | Intermediate (10%) | gobl.builder template; DPR 633/72, Table A, Part III |
| ES | Reduced (10%) | Ley 37/1992, art. 91 |
| FR | Intermediate (10%) | CGI art. 278 bis; GOBL regime: "prestations de logement" |
| DE | Reduced (7%) | GOBL regime description: "hotel accommodations" |
| CH | Intermediate (3.8%) | GOBL regime description: "accommodation services" |
| AT | Reduced (10%) | [Austrian USP](https://www.usp.gv.at/en/themen/steuern-finanzen/umsatzsteuer-ueberblick/steuersaetze-und-steuerbefreiungen-der-umsatzsteuer.html) |
| BE | Intermediate (12%) | [Royal Decree n. 20, Feb 2026](https://kpmg.com/be/en/home/insights/2026/02/itx-federal-government-approves-increase-in-the-vat-rate-on-hotels-campings-and-pesticides.html) |
| PL | Reduced (8%) | [Polish Ministry of Finance](https://podatki-arch.mf.gov.pl/en/value-added-tax/general-vat-rules-and-rates/list-of-vat-rates) |
| PT | Reduced (6%) | [PwC Portugal](https://taxsummaries.pwc.com/portugal/corporate/other-taxes) / CIVA List I |
| SE | Reduced (12%) | [Skatteverket](https://www.skatteverket.se/servicelankar/otherlanguages/englishengelska/businessesandemployers/startingandrunningaswedishbusiness/declaringtaxesbusinesses/vat/vatratesandvatexemption.4.676f4884175c97df419255d.html) |
| EL | Reduced (13%) | [PwC Greece](https://taxsummaries.pwc.com/greece/corporate/other-taxes) |
| IE | Reduced (13.5%) | [Irish Revenue](https://www.revenue.ie/en/vat/vat-on-services/accommodation/guest-and-holiday/index.aspx) |

### Freelance Regime Overrides

Localized service names plus retained taxes where the regime defines them.

| Country | Retained Taxes | Source |
|---------|---------------|--------|
| ES | IRPF 15% (professional rate) | [gobl.builder template](https://github.com/invopop/gobl.builder); Ley 35/2006 |
| IT | IRPEF 20% + SDI extension | [gobl.builder template](https://github.com/invopop/gobl.builder); DPR 600/73 |
| MX | ISR 10% + RVAT 10.67% | [LISR art. 106](https://www.diputados.gob.mx/LeyesBiblio/pdf/LISR.pdf); [LIVA art. 1-A](https://www.diputados.gob.mx/LeyesBiblio/pdf/LIVA.pdf) |
| CO | ReteRenta 11% | [Estatuto Tributario art. 392](https://estatuto.co/392) |
| DE, FR, PT, AT, NL, BE, PL, SE, CH | Standard VAT (localized names only) | GOBL regime definitions |

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
| `pl-favat-v3` | PL | Polish FA_VAT / KSeF |
| `pt-saft-v1` | PT | Portuguese SAF-T |

## How It Works

1. Reads the regime definition from GOBL (`tax.RegimeDefFor()`) to determine tax categories, rates, and currency.
2. If a scenario is set, applies domain-specific items, tax combos, charges, periods, and party structure. Regime overrides provide localized names and correct tax rates per country.
3. Generates tax IDs with correct check digits using per-country algorithms that mirror GOBL's own validators.
4. Applies addon-specific extensions, identities, and structural requirements.
5. Wraps in a GOBL envelope which triggers `Calculate()` — GOBL computes all tax totals.
6. GOBL runs `Validate()` — if the envelope builds, the invoice is guaranteed valid.

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

Generated test data and examples:

```bash
go generate ./...  # Regenerate test lists and example JSONs
```

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
