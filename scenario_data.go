package mock

import (
	"math/rand/v2"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

var scenarios = map[cbc.Key]*scenarioConfig{
	ScenarioHotel:         hotelScenario(),
	ScenarioFreelance:     freelanceScenario(),
	ScenarioReverseCharge: reverseChargeScenario(),
	ScenarioRestaurant:    restaurantScenario(),
	ScenarioEcommerce:     ecommerceScenario(),
}

// hotelScenario generates hotel/accommodation invoices with room stays
// and mixed VAT rates. Each regime override is sourced and cited below.
//
// Countries where hotels use the standard rate (no override needed — the
// fallback pickTaxCombo already returns the correct rate):
//   - DK: 25% single rate
//   - AE: 5% single rate
//   - SG: 9% single rate
//   - NL: 21% standard since Jan 2026 (https://business.gov.nl/amendments/vat-overnight-accommodation-goes-up/)
//   - GB: 20% standard (https://www.gov.uk/guidance/hotels-holiday-accommodation-and-vat-notice-7093)
//   - AR, CA, CO, MX, US: standard rate for accommodation
func hotelScenario() *scenarioConfig {
	return &scenarioConfig{
		Items: []scenarioItem{
			{
				Name:     "Double room",
				Price:    "150.00",
				Quantity: qtyNights,
				Period:   periodStay,
				TaxCombo: func(country l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
					return tax.Set{pickTaxCombo(country)}
				},
			},
			{
				Name:     "Room service",
				Price:    "25.00",
				Quantity: qtyOne,
				TaxCombo: func(country l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
					return tax.Set{pickTaxCombo(country)}
				},
			},
		},
		RegimeOverrides: map[l10n.TaxCountryCode]*scenarioOverride{
			// Source: gobl.builder/static/templates/it/hotel.json
			// DPR 633/72, Table A, Part III — accommodation at intermediate (10%).
			// Tassa di Soggiorno exempt per art. 15, DPR 633/72.
			// GOBL regime: regimes/it/tax_categories.go (RateIntermediate = 10%)
			"IT": {
				Items: hotelItems("Camera Matrimoniale", "125.00", tax.RateIntermediate,
					"Tassa di Soggiorno", "1.00"),
			},
			// ES Ley 37/1992, art. 91 — accommodation at reduced (10%).
			// GOBL regime: regimes/es/tax_categories.go (RateReduced = 10%)
			"ES": {
				Items: hotelItems("Habitacion doble", "120.00", tax.RateReduced,
					"Tasa turistica", "2.50"),
			},
			// CGI art. 278 bis — accommodation at intermediate (10%).
			// GOBL regime description: "certaines prestations de logement".
			// GOBL regime: regimes/fr/tax_categories.go (RateIntermediate = 10%)
			"FR": {
				Items: hotelItems("Chambre double", "140.00", tax.RateIntermediate,
					"Taxe de sejour", "2.00"),
			},
			// GOBL regime description: "hotel accommodations" / "Hotelunterkünfte".
			// GOBL regime: regimes/de/tax_categories.go (RateReduced = 7%)
			"DE": {
				Items: hotelItems("Doppelzimmer", "130.00", tax.RateReduced,
					"Kurtaxe", "2.50"),
			},
			// GOBL regime description: "Applies to accommodation services."
			// GOBL regime: regimes/ch/tax_categories.go (RateIntermediate = 3.8%)
			"CH": {
				Items: hotelItems("Doppelzimmer", "180.00", tax.RateIntermediate,
					"Kurtaxe", "3.00"),
			},
			// Austrian Umsatzsteuergesetz (UStG) §10(2)(4) — accommodation at reduced (10%).
			// Source: https://www.usp.gv.at/en/themen/steuern-finanzen/umsatzsteuer-ueberblick/steuersaetze-und-steuerbefreiungen-der-umsatzsteuer.html
			// GOBL regime: regimes/at/tax_categories.go (RateReduced = 10%)
			"AT": {
				Items: hotelItems("Doppelzimmer", "135.00", tax.RateReduced,
					"Ortstaxe", "2.00"),
			},
			// Royal Decree n. 20 (Feb 2026) — accommodation moved to intermediate (12%).
			// Source: https://kpmg.com/be/en/home/insights/2026/02/itx-federal-government-approves-increase-in-the-vat-rate-on-hotels-campings-and-pesticides.html
			// GOBL regime: regimes/be/tax_categories.go (RateIntermediate = 12%)
			"BE": {
				Items: hotelItems("Chambre double", "145.00", tax.RateIntermediate,
					"Taxe de sejour", "2.00"),
			},
			// Polish VAT Act, Annex 3 — accommodation at reduced (8%).
			// Source: https://podatki-arch.mf.gov.pl/en/value-added-tax/general-vat-rules-and-rates/list-of-vat-rates
			// GOBL regime: regimes/pl/tax_categories.go (RateReduced = 8%)
			"PL": {
				Items: hotelItems("Pokoj dwuosobowy", "90.00", tax.RateReduced,
					"Oplata klimatyczna", "2.00"),
			},
			// Portuguese CIVA, List I, item 2.17 — accommodation at reduced (6%).
			// Source: https://taxsummaries.pwc.com/portugal/corporate/other-taxes
			// GOBL regime: regimes/pt/tax_categories.go (RateReduced = 6%)
			"PT": {
				Items: hotelItems("Quarto duplo", "95.00", tax.RateReduced,
					"Taxa turistica", "2.00"),
			},
			// Swedish Mervärdesskattelag — accommodation at reduced (12%).
			// Source: https://www.skatteverket.se/servicelankar/otherlanguages/englishengelska/businessesandemployers/startingandrunningaswedishbusiness/declaringtaxesbusinesses/vat/vatratesandvatexemption.4.676f4884175c97df419255d.html
			// GOBL regime: regimes/se/tax_categories.go (RateReduced = 12%)
			"SE": {
				Items: hotelItems("Dubbelrum", "150.00", tax.RateReduced,
					"", ""),
			},
			// Greek VAT Code (Law 2859/2000) — accommodation at reduced (13%).
			// Source: https://taxsummaries.pwc.com/greece/corporate/other-taxes
			// GOBL regime: regimes/gr/tax_categories.go (RateReduced = 13%)
			"EL": {
				Items: hotelItems("Double room", "110.00", tax.RateReduced,
					"", ""),
			},
			// Irish VAT Consolidation Act 2010, Schedule 3 — accommodation at reduced (13.5%).
			// Source: https://www.revenue.ie/en/vat/vat-on-services/accommodation/guest-and-holiday/index.aspx
			// GOBL regime: regimes/ie/tax_categories.go (RateReduced = 13.5%)
			"IE": {
				Items: hotelItems("Double room", "140.00", tax.RateReduced,
					"", ""),
			},
		},
	}
}

// hotelItems builds a scenarioItem slice for a hotel override. When taxName
// is non-empty, a VAT-exempt lodging tax line is included.
func hotelItems(roomName, roomPrice string, roomRate cbc.Key, taxName, taxPrice string) []scenarioItem {
	items := []scenarioItem{
		{
			Name:     roomName,
			Price:    roomPrice,
			Quantity: qtyNights,
			Period:   periodStay,
			TaxCombo: func(_ l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
				return tax.Set{
					{Category: tax.CategoryVAT, Rate: roomRate},
				}
			},
		},
	}
	if taxName != "" {
		items = append(items, scenarioItem{
			Name:     taxName,
			Price:    taxPrice,
			Quantity: qtyOne,
			TaxCombo: func(_ l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
				return tax.Set{
					{Category: tax.CategoryVAT, Key: tax.KeyExempt},
				}
			},
		})
	}
	return items
}

// freelanceScenario generates freelancer/self-employed invoices with hourly
// rates, retained taxes where applicable, line-level discounts, and due dates.
//
// Regime overrides provide two things:
//  1. Localized service names in the country's language.
//  2. Retained taxes where the regime defines them in GOBL.
//
// Countries with retained taxes in GOBL but WITHOUT overrides (rates are
// variable and context-dependent — too complex for a single default):
//   - BR: IRRF (1.5-15%), INSS, PISRet, COFINSRet, CSLL — rates depend on
//     service type, contract terms, and municipal rules.
//
// Sources for retained tax overrides:
//   - ES: gobl.builder/static/templates/es/invoice-freelance.json
//     VAT standard (21%) + IRPF professional rate (15%) per Ley 35/2006.
//     GOBL regime: regimes/es/tax_categories.go (TaxCategoryIRPF, TaxRatePro = 15%).
//   - IT: gobl.builder/static/templates/it/freelance.json
//     VAT general (22%) + IRPEF at 20% retained per DPR 600/73, art. 25.
//     SDI addon requires "it-sdi-retained": "A" (persone fisiche).
//     GOBL regime: regimes/it/tax_categories.go (TaxCategoryIRPEF).
//     GOBL addon: addons/it/sdi/extensions.go (ExtKeyRetained).
//   - MX: ISR 10% per LISR art. 106 + Retained IVA (2/3 of 16%) per LIVA art. 1-A.
//     Source: https://www.diputados.gob.mx/LeyesBiblio/pdf/LISR.pdf
//     Source: https://www.diputados.gob.mx/LeyesBiblio/pdf/LIVA.pdf
//     GOBL regime: regimes/mx/tax_categories.go (ISR, RVAT — no predefined rates).
//   - CO: ReteRenta 11% for professional services (honorarios, declarant)
//     per Estatuto Tributario art. 392.
//     Source: https://estatuto.co/392
//     Source: https://www.siigo.com/blog/tabla-de-retencion-en-la-fuente/
//     GOBL regime: regimes/co/tax_categories.go (RR — no predefined rates).
func freelanceDiscount() *num.Percentage {
	d := num.MakePercentage(10, 2)
	return &d
}

func freelanceScenario() *scenarioConfig {
	return &scenarioConfig{
		SupplierPeople: true,
		Items: []scenarioItem{
			{
				Name:     "Professional services",
				Price:    "100.00",
				Unit:     "h",
				Quantity: qtyHours,
				Period:   periodBillingMonth,
				Discount: freelanceDiscount(),
				TaxCombo: func(country l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
					return tax.Set{pickTaxCombo(country)}
				},
			},
		},
		PaymentTerms: func(_ *rand.Rand) *pay.Terms {
			due := cal.Today().Add(0, 1, 0)
			pct, _ := num.PercentageFromString("100%")
			return &pay.Terms{
				Key:      pay.TermKeyDueDate,
				DueDates: []*pay.DueDate{{Date: &due, Percent: &pct}},
			}
		},
		RegimeOverrides: map[l10n.TaxCountryCode]*scenarioOverride{
			// Source: gobl.builder/static/templates/es/invoice-freelance.json
			// GOBL regime: regimes/es/tax_categories.go, regimes/es/es.go
			"ES": {
				Items: freelanceItems("Desarrollo de software", "90.00",
					"Consultoria empresarial", "75.00",
					func(_ l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
						return tax.Set{
							{Category: tax.CategoryVAT, Rate: tax.KeyStandard},
							{Category: "IRPF", Rate: "pro"},
						}
					}),
			},
			// Source: gobl.builder/static/templates/it/freelance.json
			// GOBL addon: addons/it/sdi/extensions.go (it-sdi-retained: "A")
			"IT": {
				Items: freelanceItems("Servizi di sviluppo", "90.00",
					"Consulenza aziendale", "80.00",
					func(_ l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
						pct := num.MakePercentage(200, 3)
						return tax.Set{
							{Category: tax.CategoryVAT, Rate: tax.RateGeneral},
							{Category: "IRPEF", Percent: &pct, Ext: tax.Extensions{"it-sdi-retained": "A"}},
						}
					}),
			},
			// LISR art. 106: ISR 10% on professional services.
			// LIVA art. 1-A: Retained IVA = 2/3 of 16% = 10.6667%.
			// Source: https://www.diputados.gob.mx/LeyesBiblio/pdf/LISR.pdf
			// Source: https://www.diputados.gob.mx/LeyesBiblio/pdf/LIVA.pdf
			// GOBL regime: regimes/mx/tax_categories.go (ISR, RVAT)
			"MX": {
				Items: freelanceItems("Desarrollo de software", "1500.00",
					"Consultoria empresarial", "1200.00",
					func(_ l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
						isr := num.MakePercentage(100, 3)
						rvat := num.MakePercentage(106667, 6)
						return tax.Set{
							{Category: tax.CategoryVAT, Rate: tax.KeyStandard},
							{Category: "ISR", Percent: &isr},
							{Category: "RVAT", Percent: &rvat},
						}
					}),
			},
			// Estatuto Tributario art. 392: ReteRenta 11% for honorarios (declarant).
			// Source: https://estatuto.co/392
			// Source: https://www.siigo.com/blog/tabla-de-retencion-en-la-fuente/
			// GOBL regime: regimes/co/tax_categories.go (RR)
			"CO": {
				Items: freelanceItems("Desarrollo de software", "500000.00",
					"Consultoria empresarial", "400000.00",
					func(_ l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
						rr := num.MakePercentage(110, 3)
						return tax.Set{
							{Category: tax.CategoryVAT, Rate: tax.KeyStandard},
							{Category: "RR", Percent: &rr},
						}
					}),
			},
			// Localized items — no retained taxes in these regimes.
			// GOBL regime: regimes/de/tax_categories.go
			"DE": {
				Items: freelanceItemsSimple("Softwareentwicklung", "95.00",
					"Unternehmensberatung", "85.00"),
			},
			// GOBL regime: regimes/fr/tax_categories.go
			"FR": {
				Items: freelanceItemsSimple("Developpement logiciel", "90.00",
					"Conseil en entreprise", "80.00"),
			},
			// GOBL regime: regimes/pt/tax_categories.go
			"PT": {
				Items: freelanceItemsSimple("Desenvolvimento de software", "60.00",
					"Consultoria empresarial", "50.00"),
			},
			// GOBL regime: regimes/at/tax_categories.go
			"AT": {
				Items: freelanceItemsSimple("Softwareentwicklung", "90.00",
					"Unternehmensberatung", "80.00"),
			},
			// GOBL regime: regimes/nl/tax_categories.go
			"NL": {
				Items: freelanceItemsSimple("Softwareontwikkeling", "95.00",
					"Bedrijfsadvies", "85.00"),
			},
			// GOBL regime: regimes/be/tax_categories.go
			"BE": {
				Items: freelanceItemsSimple("Developpement logiciel", "90.00",
					"Conseil en entreprise", "80.00"),
			},
			// GOBL regime: regimes/pl/tax_categories.go
			"PL": {
				Items: freelanceItemsSimple("Tworzenie oprogramowania", "60.00",
					"Doradztwo biznesowe", "50.00"),
			},
			// GOBL regime: regimes/se/tax_categories.go
			"SE": {
				Items: freelanceItemsSimple("Programvaruutveckling", "95.00",
					"Foretagsradgivning", "85.00"),
			},
			// GOBL regime: regimes/ch/tax_categories.go
			"CH": {
				Items: freelanceItemsSimple("Softwareentwicklung", "120.00",
					"Unternehmensberatung", "110.00"),
			},
		},
	}
}

// freelanceItems builds a scenarioItem slice with a shared tax combo function
// for regimes that have retained taxes. The first item gets a 10% discount.
func freelanceItems(name1, price1, name2, price2 string,
	combo func(l10n.TaxCountryCode, *addonConfig, *rand.Rand) tax.Set) []scenarioItem {
	return []scenarioItem{
		{Name: name1, Price: price1, Unit: "h", Quantity: qtyHours, Period: periodBillingMonth, Discount: freelanceDiscount(), TaxCombo: combo},
		{Name: name2, Price: price2, Unit: "h", Quantity: qtyHours, Period: periodBillingMonth, TaxCombo: combo},
	}
}

// freelanceItemsSimple builds localized freelance items that use the regime's
// standard tax combo (no retained taxes). The first item gets a 10% discount.
func freelanceItemsSimple(name1, price1, name2, price2 string) []scenarioItem {
	combo := func(country l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
		return tax.Set{pickTaxCombo(country)}
	}
	return []scenarioItem{
		{Name: name1, Price: price1, Unit: "h", Quantity: qtyHours, Period: periodBillingMonth, Discount: freelanceDiscount(), TaxCombo: combo},
		{Name: name2, Price: price2, Unit: "h", Quantity: qtyHours, Period: periodBillingMonth, TaxCombo: combo},
	}
}

// restaurantScenario generates restaurant/catering invoices with food items,
// a service charge, and reduced VAT rates where applicable.
//
// Sources for regime overrides:
//   - FR: GOBL regime description: "restauration" under intermediate rate (10%).
//     CGI art. 278 bis. GOBL regime: regimes/fr/tax_categories.go.
//   - ES: Food served in restaurants uses standard IVA (21%) in Spain.
//     Reduced rate (10%) only applies to unprocessed food products, not restaurant service.
//     Source: Ley 37/1992, art. 91.
//   - DE: Restaurant food uses standard VAT (19%). Reduced (7%) is for takeaway only.
//     Source: https://www.bundesfinanzministerium.de/Content/DE/FAQ/2020-02-18-steuererleichterungen-fuer-die-gastronomie.html
//   - IT: Restaurant food uses standard IVA (22%) in most cases.
//     Reduced (10%) applies to some prepared food. We use intermediate as reasonable default.
func restaurantScenario() *scenarioConfig {
	return &scenarioConfig{
		Charges: func(_ *rand.Rand, _ l10n.TaxCountryCode) []*bill.Charge {
			pct := num.MakePercentage(100, 3)
			return []*bill.Charge{{
				Key:     bill.ChargeKeyHandling,
				Reason:  "Service charge",
				Percent: &pct,
			}}
		},
		Items: []scenarioItem{
			{Name: "Main course", Price: "18.00", Quantity: qtySmall},
			{Name: "Appetizer", Price: "12.00", Quantity: qtySmall},
			{Name: "Dessert", Price: "8.00", Quantity: qtySmall},
			{Name: "Wine bottle", Price: "25.00", Quantity: qtyOne},
		},
		RegimeOverrides: map[l10n.TaxCountryCode]*scenarioOverride{
			// CGI art. 278 bis — restaurant food at intermediate (10%).
			// GOBL regime description: "restauration".
			"FR": {
				Items: []scenarioItem{
					{Name: "Plat principal", Price: "22.00", Quantity: qtySmall,
						TaxCombo: vatIntermediate()},
					{Name: "Entree", Price: "14.00", Quantity: qtySmall,
						TaxCombo: vatIntermediate()},
					{Name: "Dessert", Price: "10.00", Quantity: qtySmall,
						TaxCombo: vatIntermediate()},
					{Name: "Bouteille de vin", Price: "30.00", Quantity: qtyOne},
				},
			},
			"IT": {
				Items: []scenarioItem{
					{Name: "Primo piatto", Price: "16.00", Quantity: qtySmall,
						TaxCombo: vatIntermediate()},
					{Name: "Antipasto", Price: "12.00", Quantity: qtySmall,
						TaxCombo: vatIntermediate()},
					{Name: "Dolce", Price: "8.00", Quantity: qtySmall,
						TaxCombo: vatIntermediate()},
					{Name: "Vino della casa", Price: "20.00", Quantity: qtyOne},
				},
			},
			"ES": {
				Items: []scenarioItem{
					{Name: "Plato principal", Price: "18.00", Quantity: qtySmall},
					{Name: "Entrante", Price: "12.00", Quantity: qtySmall},
					{Name: "Postre", Price: "8.00", Quantity: qtySmall},
					{Name: "Botella de vino", Price: "22.00", Quantity: qtyOne},
				},
			},
			"DE": {
				Items: []scenarioItem{
					{Name: "Hauptgericht", Price: "20.00", Quantity: qtySmall},
					{Name: "Vorspeise", Price: "14.00", Quantity: qtySmall},
					{Name: "Nachspeise", Price: "9.00", Quantity: qtySmall},
					{Name: "Flasche Wein", Price: "28.00", Quantity: qtyOne},
				},
			},
			"PT": {
				Items: []scenarioItem{
					{Name: "Prato principal", Price: "14.00", Quantity: qtySmall},
					{Name: "Entrada", Price: "8.00", Quantity: qtySmall},
					{Name: "Sobremesa", Price: "6.00", Quantity: qtySmall},
					{Name: "Garrafa de vinho", Price: "18.00", Quantity: qtyOne},
				},
			},
		},
	}
}

// ecommerceScenario generates e-commerce/retail invoices with physical goods,
// shipping as a document-level charge, and standard VAT rates.
func ecommerceScenario() *scenarioConfig {
	return &scenarioConfig{
		Charges: func(r *rand.Rand, _ l10n.TaxCountryCode) []*bill.Charge {
			amt := num.MakeAmount(int64(r.IntN(10)+5), 0)
			return []*bill.Charge{{
				Key:    bill.ChargeKeyDelivery,
				Reason: "Shipping",
				Amount: amt,
			}}
		},
		Items: []scenarioItem{
			{Name: "Wireless headphones", Price: "79.99", Quantity: qtyProduct},
			{Name: "USB-C charging cable", Price: "12.99", Quantity: qtyProduct},
			{Name: "Phone case", Price: "24.99", Quantity: qtyProduct},
			{Name: "Screen protector", Price: "9.99", Quantity: qtyProduct},
		},
		RegimeOverrides: map[l10n.TaxCountryCode]*scenarioOverride{
			"ES": {
				Items: []scenarioItem{
					{Name: "Auriculares inalambricos", Price: "79.99", Quantity: qtyProduct},
					{Name: "Cable de carga USB-C", Price: "12.99", Quantity: qtyProduct},
					{Name: "Funda para movil", Price: "24.99", Quantity: qtyProduct},
					{Name: "Protector de pantalla", Price: "9.99", Quantity: qtyProduct},
				},
			},
			"FR": {
				Items: []scenarioItem{
					{Name: "Ecouteurs sans fil", Price: "79.99", Quantity: qtyProduct},
					{Name: "Cable de charge USB-C", Price: "12.99", Quantity: qtyProduct},
					{Name: "Coque de telephone", Price: "24.99", Quantity: qtyProduct},
					{Name: "Protection d'ecran", Price: "9.99", Quantity: qtyProduct},
				},
			},
			"DE": {
				Items: []scenarioItem{
					{Name: "Kabellose Kopfhoerer", Price: "79.99", Quantity: qtyProduct},
					{Name: "USB-C Ladekabel", Price: "12.99", Quantity: qtyProduct},
					{Name: "Handyhuelle", Price: "24.99", Quantity: qtyProduct},
					{Name: "Displayschutzfolie", Price: "9.99", Quantity: qtyProduct},
				},
			},
			"IT": {
				Items: []scenarioItem{
					{Name: "Cuffie wireless", Price: "79.99", Quantity: qtyProduct},
					{Name: "Cavo di ricarica USB-C", Price: "12.99", Quantity: qtyProduct},
					{Name: "Cover per telefono", Price: "24.99", Quantity: qtyProduct},
					{Name: "Pellicola protettiva", Price: "9.99", Quantity: qtyProduct},
				},
			},
			"PT": {
				Items: []scenarioItem{
					{Name: "Auscultadores sem fios", Price: "79.99", Quantity: qtyProduct},
					{Name: "Cabo de carga USB-C", Price: "12.99", Quantity: qtyProduct},
					{Name: "Capa para telemovel", Price: "24.99", Quantity: qtyProduct},
					{Name: "Pelicula protetora", Price: "9.99", Quantity: qtyProduct},
				},
			},
		},
	}
}

// vatIntermediate returns a TaxCombo function for VAT intermediate rate.
func vatIntermediate() func(l10n.TaxCountryCode, *addonConfig, *rand.Rand) tax.Set {
	return func(_ l10n.TaxCountryCode, _ *addonConfig, _ *rand.Rand) tax.Set {
		return tax.Set{{Category: tax.CategoryVAT, Rate: tax.RateIntermediate}}
	}
}

// reverseChargeScenario generates invoices with reverse-charge VAT treatment
// for cross-border B2B transactions.
//
// The reverse-charge tag triggers GOBL's built-in scenario system
// (bill/invoice_scenarios.go) which auto-applies the legal note:
// "Reverse charge: Customer to account for VAT to the relevant tax authority."
//
// For IT+SDI addon, the normalizer at addons/it/sdi/tax_combo.go automatically
// sets "it-sdi-exempt" to "N6.9" when Key=reverse-charge.
//
// Source: EU VAT Directive 2006/112/EC, Articles 194-199.
// Cross-border B2B services within the EU use the reverse-charge mechanism.
func reverseChargeScenario() *scenarioConfig {
	// Cross-border customer countries by regime — common EU trade pairings.
	crossBorder := map[l10n.TaxCountryCode]l10n.TaxCountryCode{
		"ES": "NL",
		"IT": "DE",
		"FR": "BE",
		"DE": "NL",
		"NL": "DE",
		"BE": "FR",
		"AT": "DE",
		"PT": "ES",
		"SE": "DK",
		"DK": "SE",
		"PL": "DE",
		"IE": "FR",
		"GR": "IT",
		"EL": "IT",
	}

	return &scenarioConfig{
		Tags: []cbc.Key{tax.TagReverseCharge},
		// No custom items — uses locale items with reverse-charge tax treatment.
		Items: nil,
		CustomerCountry: func(_ *rand.Rand, country l10n.TaxCountryCode) l10n.TaxCountryCode {
			if cc, ok := crossBorder[country]; ok {
				return cc
			}
			return "NL"
		},
	}
}
