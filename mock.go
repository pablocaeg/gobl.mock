// Package mock generates realistic, valid GOBL invoices for any supported tax regime.
//
//go:generate go run ./internal/generate/main.go
package mock

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"

	_ "github.com/invopop/gobl/addons"  // register all addons
	_ "github.com/invopop/gobl/regimes" // register all regimes
)

// Invoice generates a valid GOBL invoice for the given options.
// It is safe for concurrent use from multiple goroutines.
func Invoice(opts ...Option) (*bill.Invoice, error) {
	env, err := Envelope(opts...)
	if err != nil {
		return nil, err
	}
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("unexpected document type")
	}
	return inv, nil
}

// Envelope generates a valid GOBL envelope containing an invoice.
// It is safe for concurrent use from multiple goroutines.
func Envelope(opts ...Option) (*gobl.Envelope, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.lines < 1 {
		o.lines = 1
	}

	r := newRand(o)
	country := o.regime
	addon := resolveAddon(country, o.addon)
	locale := getLocale(country)
	ac := resolveAddonConfig(addon)

	regime := tax.RegimeDefFor(l10n.Code(country))
	if regime == nil {
		return nil, fmt.Errorf("unsupported regime: %s", country)
	}

	sc, err := resolveScenario(o.scenario)
	if err != nil {
		return nil, err
	}

	inv := buildInvoice(r, o, country, addon, locale, ac, sc)

	// Apply template overrides if provided.
	if o.template != nil {
		applyTemplate(inv, o.template)
	}

	// Apply invoice type (standard, credit-note, corrective, debit-note, proforma).
	if o.invType != "" {
		applyInvoiceType(r, inv, o.invType, addon, ac, regime)
	}

	env, err := gobl.Envelop(inv)
	if err != nil {
		return nil, fmt.Errorf("building envelope: %w", err)
	}
	return env, nil
}

func buildInvoice(r *rand.Rand, o *options, country l10n.TaxCountryCode, addon cbc.Key, locale *localeData, ac *addonConfig, sc *scenarioConfig) *bill.Invoice {
	series := cbc.Code("MOCK")
	if ac != nil && ac.Series != "" {
		series = cbc.Code(ac.Series)
	} else if ac != nil && ac.NumericSeries {
		series = cbc.Code(fmt.Sprintf("%04d", r.IntN(9998)+1))
	}

	inv := &bill.Invoice{
		Series:    series,
		Code:      cbc.Code(fmt.Sprintf("%05d", r.IntN(99999)+1)),
		IssueDate: cal.Today(),
		Currency:  regimeCurrency(country),
	}
	inv.SetRegime(country)

	if addon != "" {
		inv.SetAddons(addon)
	}
	if o.simplified {
		inv.SetTags(tax.TagSimplified)
	}

	if ac != nil && ac.InvoiceTaxExt != nil {
		inv.Tax = &bill.Tax{Ext: ac.InvoiceTaxExt(r)}
	}
	if ac != nil && len(ac.Notes) > 0 {
		inv.Notes = ac.Notes
	}

	inv.Supplier = buildParty(r, country, locale, locale.SupplierNames, ac, true)
	if !o.simplified {
		inv.Customer = buildParty(r, country, locale, locale.CustomerNames, ac, false)
	}
	inv.Lines = buildLines(r, country, locale, ac, o.lines)
	inv.Payment = buildPayment(r, locale, ac)

	if ac != nil && ac.RequiresOrdering {
		inv.Ordering = &bill.Ordering{
			Code: cbc.Code(fmt.Sprintf("PO-%05d", r.IntN(99999)+1)),
		}
		if ac.NumericSeries {
			start := cal.Today().Add(0, -1, 0)
			end := cal.Today()
			inv.Ordering.Period = &cal.Period{Start: start, End: end}
		}
	}

	if sc != nil {
		applyScenario(r, inv, o, country, addon, ac, sc)
	}

	return inv
}

// applyScenario applies domain-specific overrides to the invoice.
func applyScenario(r *rand.Rand, inv *bill.Invoice, o *options, country l10n.TaxCountryCode, addon cbc.Key, ac *addonConfig, sc *scenarioConfig) {
	// Tags.
	if len(sc.Tags) > 0 {
		tags := sc.Tags
		if o.simplified {
			tags = append(tags, tax.TagSimplified)
		}
		inv.SetTags(tags...)
	}

	// Supplier: add people if needed.
	applyScenarioToSupplier(r, inv.Supplier, sc, country)

	// Customer: scenario may override country for cross-border.
	if !o.simplified && sc.CustomerCountry != nil {
		custCountry := sc.CustomerCountry(r, country)
		custLocale := getLocale(custCountry)
		inv.Customer = buildParty(r, custCountry, custLocale, custLocale.CustomerNames, ac, false)
	}

	// Lines: scenario items replace locale items when present.
	if len(sc.resolveItems(country)) > 0 {
		inv.Lines = buildScenarioLines(r, country, addon, ac, sc, o.lines)
	}

	// Reverse-charge without custom items: override tax combos.
	applyReverseChargeTaxes(inv, sc, country)

	// Document-level charges (skip for addons that reject them).
	if sc.Charges != nil && !addonBlocksCharges(addon) {
		inv.Charges = sc.Charges(r, country)
	}

	// Payment terms override.
	if sc.PaymentTerms != nil {
		inv.Payment.Terms = sc.PaymentTerms(r)
	}
}

// applyReverseChargeTaxes overrides line tax combos for reverse-charge
// scenarios that use locale items (no custom scenario items).
// Only applies to VAT regimes — reverse charge is a VAT-specific concept.
func applyReverseChargeTaxes(inv *bill.Invoice, sc *scenarioConfig, country l10n.TaxCountryCode) {
	if len(sc.resolveItems(country)) > 0 {
		return
	}
	for _, tag := range sc.Tags {
		if tag == tax.TagReverseCharge {
			combo := pickTaxCombo(country)
			if combo.Category == tax.CategoryVAT {
				combo.Key = tax.KeyReverseCharge
				combo.Rate = ""
				for _, line := range inv.Lines {
					line.Taxes = tax.Set{combo}
				}
			}
			return
		}
	}
}

// applyInvoiceType sets the invoice type and adds preceding references
// for correction types (credit-note, corrective, debit-note).
func applyInvoiceType(r *rand.Rand, inv *bill.Invoice, invType cbc.Key, addon cbc.Key, ac *addonConfig, regime *tax.RegimeDef) {
	inv.Type = invType

	switch invType {
	case bill.InvoiceTypeProforma:
		return // proforma needs no preceding reference
	case bill.InvoiceTypeStandard:
		return
	}

	// Correction types need a preceding reference.
	// Check if the regime supports this type.
	if !regimeSupportsType(regime, addon, invType) {
		// Fall back to credit-note which is universally supported.
		inv.Type = bill.InvoiceTypeCreditNote
	}

	yesterday := cal.Today().Add(0, 0, -1)
	ref := &org.DocumentRef{
		Identify:  uuid.Identify{UUID: uuid.V7()},
		Type:      bill.InvoiceTypeStandard,
		Series:    inv.Series,
		Code:      cbc.Code(fmt.Sprintf("%05d", r.IntN(99999)+1)),
		IssueDate: &yesterday,
	}

	if ac != nil {
		if ac.CorrectionExt != nil {
			ref.Ext = ac.CorrectionExt(r)
		}
		if ac.CorrectionStamps != nil {
			for _, s := range ac.CorrectionStamps(r) {
				ref.Stamps = append(ref.Stamps, &head.Stamp{
					Provider: s.Provider,
					Value:    s.Value,
				})
			}
		}
	}

	inv.Preceding = []*org.DocumentRef{ref}
}

// regimeSupportsType checks if the regime (with addon) supports a given correction type.
func regimeSupportsType(regime *tax.RegimeDef, addon cbc.Key, invType cbc.Key) bool {
	sets := []tax.CorrectionSet{regime.Corrections}
	if addon != "" {
		if ad := tax.AddonForKey(addon); ad != nil {
			sets = append(sets, ad.Corrections)
		}
	}
	for _, cs := range sets {
		for _, cd := range cs {
			for _, t := range cd.Types {
				if t == invType {
					return true
				}
			}
		}
	}
	return false
}

// applyTemplate merges non-zero fields from the template into the invoice.
func applyTemplate(inv *bill.Invoice, tmpl *bill.Invoice) {
	if tmpl.Series != "" {
		inv.Series = tmpl.Series
	}
	if tmpl.Code != "" {
		inv.Code = tmpl.Code
	}
	if !tmpl.IssueDate.IsZero() {
		inv.IssueDate = tmpl.IssueDate
	}
	if tmpl.Currency != "" {
		inv.Currency = tmpl.Currency
	}
	if tmpl.Supplier != nil {
		inv.Supplier = tmpl.Supplier
	}
	if tmpl.Customer != nil {
		inv.Customer = tmpl.Customer
	}
	if tmpl.Lines != nil {
		inv.Lines = tmpl.Lines
	}
	if tmpl.Payment != nil {
		inv.Payment = tmpl.Payment
	}
	if tmpl.Ordering != nil {
		inv.Ordering = tmpl.Ordering
	}
	if tmpl.Tax != nil {
		inv.Tax = tmpl.Tax
	}
	if tmpl.Notes != nil {
		inv.Notes = tmpl.Notes
	}
	if tmpl.Preceding != nil {
		inv.Preceding = tmpl.Preceding
	}
}

// resolveAddonConfig returns the addon config for a known addon, or builds
// a dynamic config for unknown addons by reading extension definitions from
// GOBL. Scenarios handle most extension auto-setting; the dynamic fallback
// provides first-valid-value defaults for extensions not covered by scenarios.
func resolveAddonConfig(addon cbc.Key) *addonConfig {
	if addon == "" {
		return nil
	}
	if ac, ok := addons[addon]; ok {
		return ac
	}
	return buildDynamicAddonConfig(addon)
}

func buildDynamicAddonConfig(addon cbc.Key) *addonConfig {
	ad := tax.AddonForKey(addon)
	if ad == nil {
		return &addonConfig{}
	}

	// Build a lookup of extension key → first valid value.
	firstValue := make(map[cbc.Key]cbc.Code)
	for _, def := range ad.Extensions {
		if len(def.Values) > 0 {
			firstValue[def.Key] = def.Values[0].Code
		}
	}

	if len(firstValue) == 0 {
		return &addonConfig{}
	}

	// Probe the addon's validator to discover which extensions are required
	// on which object types, then assign first valid values accordingly.
	mapping := probeAddonExtensions(ad)

	ac := &addonConfig{}

	if keys, ok := mapping["invoice.tax"]; ok {
		ac.InvoiceTaxExt = func(_ *rand.Rand) tax.Extensions {
			return pickExtValues(keys, firstValue)
		}
	}
	if keys, ok := mapping["party"]; ok {
		fn := func(_ *rand.Rand) tax.Extensions {
			return pickExtValues(keys, firstValue)
		}
		ac.SupplierExt = fn
		ac.CustomerExt = fn
	}
	if keys, ok := mapping["combo"]; ok {
		ac.ComboExt = func(_ *rand.Rand) tax.Extensions {
			return pickExtValues(keys, firstValue)
		}
	}
	if keys, ok := mapping["item"]; ok {
		ac.ItemExt = func(_ *rand.Rand) tax.Extensions {
			return pickExtValues(keys, firstValue)
		}
	}
	if keys, ok := mapping["pay"]; ok {
		ac.PaymentExt = func(_ *rand.Rand) tax.Extensions {
			return pickExtValues(keys, firstValue)
		}
	}

	return ac
}

// pickExtValues creates an Extensions map using the first valid value
// for each requested key.
func pickExtValues(keys []cbc.Key, firstValue map[cbc.Key]cbc.Code) tax.Extensions {
	ext := make(tax.Extensions, len(keys))
	for _, k := range keys {
		if v, ok := firstValue[k]; ok {
			ext[k] = v
		}
	}
	return ext
}

func buildParty(r *rand.Rand, country l10n.TaxCountryCode, locale *localeData, names []string, ac *addonConfig, isSupplier bool) *org.Party {
	city := pick(r, locale.Cities)
	postalCode := city.Code
	if fn, ok := postalCodeFormats[country]; ok {
		postalCode = fn(r)
	}

	p := &org.Party{
		Name: pick(r, names),
		TaxID: &tax.Identity{
			Country: country,
			Code:    generateTaxID(r, country, true),
		},
		Addresses: []*org.Address{{
			Street:   pick(r, locale.Streets),
			Number:   fmt.Sprintf("%d", r.IntN(200)+1),
			Locality: city.Name,
			Region:   city.Region,
			Code:     postalCode,
			Country:  isoCountry(country),
		}},
		Emails: []*org.Email{{Address: "billing@example.com"}},
	}

	if ac == nil {
		return p
	}

	if isSupplier {
		if ac.SupplierExt != nil {
			p.Ext = ac.SupplierExt(r)
		}
		if ac.SupplierState != "" {
			p.Addresses[0].State = ac.SupplierState
		}
		if ac.SupplierIdentities != nil {
			p.Identities = ac.SupplierIdentities(r)
		}
		if ac.SupplierPeople {
			p.People = []*org.Person{{
				Name:   &org.Name{Given: "Test", Surname: "Person"},
				Emails: []*org.Email{{Address: "contact@example.com"}},
			}}
		}
		if ac.SupplierInboxes {
			p.Inboxes = []*org.Inbox{{Email: "invoices@example.com"}}
		}
		if ac.SupplierPhones {
			p.Telephones = []*org.Telephone{{Number: "+1234567890"}}
		}
	} else {
		if ac.CustomerExt != nil {
			p.Ext = ac.CustomerExt(r)
		}
		if ac.CustomerState != "" {
			p.Addresses[0].State = ac.CustomerState
		}
		if ac.CustomerPostalCode != nil {
			p.Addresses[0].Code = ac.CustomerPostalCode(r)
		}
		if ac.CustomerIdentities != nil {
			p.Identities = ac.CustomerIdentities(r)
		}
		if ac.CustomerInboxes {
			p.Inboxes = []*org.Inbox{{Email: "billing@customer.com"}}
		}
	}
	return p
}

func buildLines(r *rand.Rand, country l10n.TaxCountryCode, locale *localeData, ac *addonConfig, count int) []*bill.Line {
	lines := make([]*bill.Line, count)
	for i := range lines {
		lines[i] = buildLine(r, country, locale, ac)
	}
	return lines
}

func buildLine(r *rand.Rand, country l10n.TaxCountryCode, locale *localeData, ac *addonConfig) *bill.Line {
	// Build a pool of available items from products and services.
	items := make([]itemData, 0, len(locale.Products)+len(locale.Services))
	items = append(items, locale.Products...)
	items = append(items, locale.Services...)
	item := pick(r, items)

	price, _ := num.AmountFromString(item.Price)
	lineItem := &org.Item{
		Name:  item.Name,
		Price: &price,
	}
	if item.Unit != "" {
		lineItem.Unit = org.Unit(item.Unit)
	}
	if ac != nil && ac.ItemExt != nil {
		lineItem.Ext = ac.ItemExt(r)
	}
	if fn, ok := itemIdentities[country]; ok {
		lineItem.Identities = fn(r)
	}

	combo := pickTaxCombo(country)
	if ac != nil && ac.ComboExt != nil {
		combo.Ext = ac.ComboExt(r)
	}

	taxes := tax.Set{combo}
	if ac != nil && ac.ExtraCombos != nil {
		taxes = append(taxes, ac.ExtraCombos(r)...)
	}

	return &bill.Line{
		Quantity: num.MakeAmount(int64(r.IntN(20)+1), 0),
		Item:     lineItem,
		Taxes:    taxes,
	}
}

func buildPayment(r *rand.Rand, locale *localeData, ac *addonConfig) *bill.PaymentDetails {
	key := locale.PaymentKey
	if key == "" {
		key = pay.MeansKeyCreditTransfer
	}

	instructions := &pay.Instructions{Key: key}
	if locale.IBANPrefix != "" {
		instructions.CreditTransfer = []*pay.CreditTransfer{{
			IBAN: generateIBAN(r, locale.IBANPrefix),
			Name: "Bank Account",
		}}
	}
	if ac != nil && ac.PaymentExt != nil {
		instructions.Ext = ac.PaymentExt(r)
	}

	terms := &pay.Terms{Key: pay.TermKeyInstant, Notes: "Payment due upon receipt."}
	if ac != nil && ac.RequiresDueDates {
		due := cal.Today().Add(0, 1, 0)
		pct, _ := num.PercentageFromString("100%")
		terms = &pay.Terms{
			Key:      pay.TermKeyDueDate,
			DueDates: []*pay.DueDate{{Date: &due, Percent: &pct}},
		}
	}

	return &bill.PaymentDetails{Instructions: instructions, Terms: terms}
}

func resolveAddon(country l10n.TaxCountryCode, explicit cbc.Key) cbc.Key {
	if explicit != "" {
		return explicit
	}
	return defaultAddonForRegime(country)
}

func regimeCurrency(country l10n.TaxCountryCode) currency.Code {
	r := tax.RegimeDefFor(l10n.Code(country))
	if r != nil {
		return r.GetCurrency()
	}
	return currency.USD
}

func pickTaxCombo(country l10n.TaxCountryCode) *tax.Combo {
	regime := tax.RegimeDefFor(l10n.Code(country))
	if regime == nil {
		return &tax.Combo{Category: tax.CategoryVAT, Rate: tax.KeyStandard}
	}

	var primary *tax.CategoryDef
	for _, cat := range regime.Categories {
		if cat.Retained || cat.Informative {
			continue
		}
		if cat.Code == tax.CategoryVAT || cat.Code == tax.CategoryGST || cat.Code == tax.CategoryST {
			primary = cat
			break
		}
		if primary == nil {
			primary = cat
		}
	}
	if primary == nil {
		return &tax.Combo{Category: tax.CategoryVAT, Rate: tax.KeyStandard}
	}

	if primary.Code == tax.CategoryVAT {
		return &tax.Combo{Category: primary.Code, Rate: tax.KeyStandard}
	}

	for _, rate := range primary.Rates {
		if len(rate.Values) > 0 && !rate.Values[0].Disabled {
			return &tax.Combo{Category: primary.Code, Rate: rate.Rate}
		}
	}

	pct := num.MakePercentage(10, 2)
	return &tax.Combo{Category: primary.Code, Percent: &pct}
}

// isoCountry converts a tax country code to an ISO country code.
// Most are identical, but Greece uses "EL" for tax and "GR" for ISO.
func isoCountry(tc l10n.TaxCountryCode) l10n.ISOCountryCode {
	if tc == "EL" {
		return "GR"
	}
	return l10n.ISOCountryCode(tc)
}

// BBAN lengths per country (IBAN total length - 4).
var bbanLengths = map[string]int{
	"AE": 19, "AT": 16, "BE": 12, "BR": 25, "CH": 17,
	"DE": 18, "DK": 14, "ES": 20, "FR": 23, "GB": 18,
	"GR": 23, "IE": 18, "IT": 23, "NL": 14, "PL": 24,
	"PT": 21, "SE": 20,
}

// generateIBAN creates an IBAN with valid mod-97 check digits.
func generateIBAN(r *rand.Rand, country string) string {
	bbanLen := bbanLengths[country]
	if bbanLen == 0 {
		bbanLen = 18 // safe default
	}

	// Generate random BBAN (digits only for simplicity).
	bban := make([]byte, bbanLen)
	bban[0] = byte('1' + r.IntN(9)) // avoid leading zero
	for i := 1; i < bbanLen; i++ {
		bban[i] = byte('0' + r.IntN(10))
	}

	// Compute check digits: rearrange to BBAN + country + "00", convert letters, mod 97.
	check := ibanCheckDigits(country, string(bban))
	return fmt.Sprintf("%s%s%s", country, check, string(bban))
}

// ibanCheckDigits computes the 2-digit IBAN check using ISO 7064 mod 97-10.
func ibanCheckDigits(country, bban string) string {
	// Rearrange: BBAN + CountryLetters + "00"
	raw := bban + country + "00"

	// Convert letters to numbers: A=10, B=11, ..., Z=35
	var numeric string
	for _, ch := range raw {
		if ch >= 'A' && ch <= 'Z' {
			numeric += fmt.Sprintf("%d", ch-'A'+10)
		} else {
			numeric += string(ch)
		}
	}

	// Compute mod 97 using iterative chunk approach (avoids big.Int).
	remainder := 0
	for _, ch := range numeric {
		remainder = (remainder*10 + int(ch-'0')) % 97
	}

	return fmt.Sprintf("%02d", 98-remainder)
}

var postalCodeFormats = map[l10n.TaxCountryCode]func(r *rand.Rand) cbc.Code{
	"BR": func(r *rand.Rand) cbc.Code {
		return cbc.Code(fmt.Sprintf("%05d-%03d", r.IntN(99999), r.IntN(999)))
	},
}

var itemIdentities = map[l10n.TaxCountryCode]func(r *rand.Rand) []*org.Identity{
	"IN": func(r *rand.Rand) []*org.Identity {
		return []*org.Identity{{
			Type: "HSN",
			Code: cbc.Code(fmt.Sprintf("%04d%04d", r.IntN(9999)+1, r.IntN(9999))),
		}}
	},
}

func pick[T any](r *rand.Rand, items []T) T {
	if len(items) == 0 {
		var zero T
		return zero
	}
	return items[r.IntN(len(items))]
}

// addonBlocksCharges returns true if the addon rejects document-level charges.
// Currently only mx-cfdi-v4 does this.
func addonBlocksCharges(addon cbc.Key) bool {
	return addon == "mx-cfdi-v4"
}

func newRand(o *options) *rand.Rand {
	if o.hasSeed {
		return rand.New(rand.NewPCG(uint64(o.seed), uint64(o.seed)))
	}
	return rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()+1)))
}
