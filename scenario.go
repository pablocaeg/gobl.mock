package mock

import (
	"fmt"
	"math/rand/v2"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

// Scenario keys for domain-specific invoice generation.
const (
	ScenarioHotel         cbc.Key = "hotel"
	ScenarioFreelance     cbc.Key = "freelance"
	ScenarioReverseCharge cbc.Key = "reverse-charge"
	ScenarioRestaurant    cbc.Key = "restaurant"
	ScenarioEcommerce     cbc.Key = "ecommerce"
)

// scenarioConfig defines a domain-specific invoice generation profile.
type scenarioConfig struct {
	// Tags to apply to the invoice (e.g. reverse-charge).
	Tags []cbc.Key

	// Items to use instead of locale products/services.
	Items []scenarioItem

	// SupplierPeople adds a Person entry to the supplier party.
	SupplierPeople bool

	// Charges generates document-level charges for the invoice.
	Charges func(r *rand.Rand, country l10n.TaxCountryCode) []*bill.Charge

	// CustomerCountry overrides the customer's country for cross-border scenarios.
	CustomerCountry func(r *rand.Rand, country l10n.TaxCountryCode) l10n.TaxCountryCode

	// PaymentTerms overrides the default payment terms.
	PaymentTerms func(r *rand.Rand) *pay.Terms

	// RegimeOverrides provides per-regime customizations.
	RegimeOverrides map[l10n.TaxCountryCode]*scenarioOverride
}

// scenarioItem defines a line item template with its associated tax set.
type scenarioItem struct {
	Name     string
	Price    string
	Unit     string          // optional (e.g. "h" for hours)
	Discount *num.Percentage // optional line-level discount
	// Quantity returns a realistic quantity for this item. When nil,
	// defaults to random 1-20.
	Quantity func(r *rand.Rand) num.Amount
	// Period returns a service period for this line. When nil, no period is set.
	Period func(r *rand.Rand) *cal.Period
	// TaxCombo returns the tax.Set for this item in the given country.
	// When nil, falls back to the default pickTaxCombo behavior.
	TaxCombo func(country l10n.TaxCountryCode, ac *addonConfig, r *rand.Rand) tax.Set
}

// scenarioOverride customizes a scenario for a specific regime.
type scenarioOverride struct {
	// Items replaces the base items for this regime.
	Items []scenarioItem
}

// resolveScenario looks up the scenario config for the given key.
// Returns nil if the key is empty or unknown.
func resolveScenario(key cbc.Key) (*scenarioConfig, error) {
	if key == "" {
		return nil, nil
	}
	sc, ok := scenarios[key]
	if !ok {
		return nil, fmt.Errorf("unknown scenario: %s", key)
	}
	return sc, nil
}

// resolveItems returns the regime-specific items if available, otherwise the base items.
func (sc *scenarioConfig) resolveItems(country l10n.TaxCountryCode) []scenarioItem {
	if ov, ok := sc.RegimeOverrides[country]; ok && len(ov.Items) > 0 {
		return ov.Items
	}
	return sc.Items
}

// buildScenarioLines builds invoice lines from scenario items.
func buildScenarioLines(r *rand.Rand, country l10n.TaxCountryCode, addon cbc.Key, ac *addonConfig, sc *scenarioConfig, count int) []*bill.Line {
	items := sc.resolveItems(country)
	lines := make([]*bill.Line, count)
	for i := range lines {
		si := items[r.IntN(len(items))]
		lines[i] = buildScenarioLine(r, country, addon, ac, si)
	}
	return lines
}

// buildScenarioLine builds a single invoice line from a scenario item.
func buildScenarioLine(r *rand.Rand, country l10n.TaxCountryCode, addon cbc.Key, ac *addonConfig, si scenarioItem) *bill.Line {
	price, _ := num.AmountFromString(si.Price)
	lineItem := &org.Item{
		Name:  si.Name,
		Price: &price,
	}
	if si.Unit != "" {
		lineItem.Unit = org.Unit(si.Unit)
	}
	if ac != nil && ac.ItemExt != nil {
		lineItem.Ext = ac.ItemExt(r)
	}
	if fn, ok := itemIdentities[country]; ok {
		lineItem.Identities = fn(r)
	}

	// Determine tax set.
	var taxes tax.Set
	if si.TaxCombo != nil {
		taxes = si.TaxCombo(country, ac, r)
	} else {
		combo := pickTaxCombo(country)
		if ac != nil && ac.ComboExt != nil {
			combo.Ext = ac.ComboExt(r)
		}
		taxes = tax.Set{combo}
	}

	// Apply addon combo extensions to any combo that doesn't already have them.
	if ac != nil && ac.ComboExt != nil {
		for _, combo := range taxes {
			if len(combo.Ext) == 0 {
				combo.Ext = ac.ComboExt(r)
			}
		}
	}

	// EN16931-family addons require cef-vatex on exempt VAT combos (BR-E-10).
	// This is a conditional requirement that the probe cannot discover.
	if en16931Addons[addon] {
		for _, combo := range taxes {
			if combo.Category == tax.CategoryVAT && combo.Key == tax.KeyExempt {
				if _, ok := combo.Ext["cef-vatex"]; !ok {
					combo.Ext = combo.Ext.Merge(tax.Extensions{
						"cef-vatex": "VATEX-EU-132",
					})
				}
			}
		}
	}

	if ac != nil && ac.ExtraCombos != nil {
		taxes = append(taxes, ac.ExtraCombos(r)...)
	}

	// Quantity: use scenario-specific function or default random 1-20.
	qty := num.MakeAmount(int64(r.IntN(20)+1), 0)
	if si.Quantity != nil {
		qty = si.Quantity(r)
	}

	line := &bill.Line{
		Quantity: qty,
		Item:     lineItem,
		Taxes:    taxes,
	}

	if si.Discount != nil {
		line.Discounts = []*bill.LineDiscount{{
			Percent: si.Discount,
			Reason:  "Discount",
		}}
	}

	if si.Period != nil {
		line.Period = si.Period(r)
	}

	return line
}

// EN16931-family addons that require cef-vatex on exempt VAT combos.
var en16931Addons = map[cbc.Key]bool{
	"eu-en16931-v2017": true,
	"de-xrechnung-v3":  true,
	"de-zugferd-v2":    true,
	"fr-facturx-v1":    true,
}

// Per-locale freelancer names for the SupplierPeople feature.
var freelancerNames = map[l10n.TaxCountryCode]*org.Name{
	"ES": {Given: "Maria Francisca", Surname: "Montero", Surname2: "Esteban"},
	"MX": {Given: "Ana Lucia", Surname: "Hernandez", Surname2: "Torres"},
	"AR": {Given: "Valentina", Surname: "Lopez", Surname2: "Fernandez"},
	"CO": {Given: "Camila Andrea", Surname: "Restrepo", Surname2: "Mejia"},
	"IT": {Given: "Giulia", Surname: "Rossi"},
	"FR": {Given: "Marie", Surname: "Dubois"},
	"DE": {Given: "Anna", Surname: "Schmidt"},
	"AT": {Given: "Sophie", Surname: "Huber"},
	"CH": {Given: "Lea", Surname: "Mueller"},
	"PT": {Given: "Ana", Surname: "Silva"},
	"NL": {Given: "Emma", Surname: "de Vries"},
	"BE": {Given: "Louise", Surname: "Peeters"},
	"PL": {Given: "Anna", Surname: "Kowalska"},
	"SE": {Given: "Emma", Surname: "Andersson"},
	"GB": {Given: "Sarah", Surname: "Johnson"},
}

var defaultFreelancerName = &org.Name{Given: "Maria", Surname: "Garcia"}

// applyScenarioToSupplier modifies the supplier party based on scenario config.
func applyScenarioToSupplier(_ *rand.Rand, p *org.Party, sc *scenarioConfig, country l10n.TaxCountryCode) {
	if sc.SupplierPeople {
		name := defaultFreelancerName
		if n, ok := freelancerNames[country]; ok {
			name = n
		}
		p.People = []*org.Person{{
			Name:   name,
			Emails: []*org.Email{{Address: "contact@example.com"}},
		}}
	}
}

// Quantity helpers for realistic amounts.
func qtyHours(r *rand.Rand) num.Amount {
	// Realistic work hours: 8, 16, 24, 32, 40, 60, 80, 120, 160
	hours := []int64{8, 16, 24, 32, 40, 60, 80, 120, 160}
	return num.MakeAmount(hours[r.IntN(len(hours))], 0)
}

func qtyNights(r *rand.Rand) num.Amount {
	return num.MakeAmount(int64(r.IntN(5)+1), 0) // 1-5 nights
}

func qtyOne(_ *rand.Rand) num.Amount {
	return num.MakeAmount(1, 0)
}

func qtySmall(r *rand.Rand) num.Amount {
	return num.MakeAmount(int64(r.IntN(4)+1), 0) // 1-4
}

func qtyProduct(r *rand.Rand) num.Amount {
	return num.MakeAmount(int64(r.IntN(5)+1), 0) // 1-5
}

// Period helpers for realistic date ranges.

// periodStay generates a hotel stay period ending yesterday (recent checkout).
func periodStay(r *rand.Rand) *cal.Period {
	nights := r.IntN(5) + 1
	end := cal.Today().Add(0, 0, -1)
	start := end.Add(0, 0, -nights)
	return &cal.Period{
		Start: start,
		End:   end,
	}
}

// periodBillingMonth generates last month as a billing period.
func periodBillingMonth(_ *rand.Rand) *cal.Period {
	today := cal.Today()
	// First day of current month, then back one month for start, back one day for end.
	start := today.Add(0, -1, -(today.Day - 1))
	end := today.Add(0, 0, -today.Day)
	return &cal.Period{
		Label: "Billing period",
		Start: start,
		End:   end,
	}
}
