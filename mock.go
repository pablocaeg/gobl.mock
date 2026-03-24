// Package mock generates realistic, valid GOBL invoices for any supported tax regime.
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
func Envelope(opts ...Option) (*gobl.Envelope, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	r := newRand(o)
	country := o.Regime
	addon := resolveAddon(country, o.Addon)
	locale := getLocale(country)
	ac := resolveAddonConfig(addon)

	regime := tax.RegimeDefFor(l10n.Code(country))
	if regime == nil {
		return nil, fmt.Errorf("unsupported regime: %s", country)
	}

	inv := buildInvoice(r, o, country, addon, locale, ac)

	// Apply template overrides if provided.
	if o.Template != nil {
		applyTemplate(inv, o.Template)
	}

	// Apply invoice type (standard, credit-note, corrective, debit-note, proforma).
	if o.Type != "" {
		applyInvoiceType(r, inv, o.Type, addon, ac, regime)
	}

	env, err := gobl.Envelop(inv)
	if err != nil {
		return nil, fmt.Errorf("building envelope: %w", err)
	}
	return env, nil
}

func buildInvoice(r *rand.Rand, o *Options, country l10n.TaxCountryCode, addon cbc.Key, locale *localeData, ac *addonConfig) *bill.Invoice {
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
	if o.Simplified {
		inv.SetTags(tax.TagSimplified)
	}

	if ac != nil && ac.InvoiceTaxExt != nil {
		inv.Tax = &bill.Tax{Ext: ac.InvoiceTaxExt(r)}
	}
	if ac != nil && len(ac.Notes) > 0 {
		inv.Notes = ac.Notes
	}

	inv.Supplier = buildParty(r, country, locale, locale.SupplierNames, ac, true)
	if !o.Simplified {
		inv.Customer = buildParty(r, country, locale, locale.CustomerNames, ac, false)
	}
	inv.Lines = buildLines(r, country, locale, ac, o.Lines)
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

	return inv
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

// resolveAddonConfig returns the addon config, falling back to dynamic
// resolution for unknown addons.
func resolveAddonConfig(addon cbc.Key) *addonConfig {
	if addon == "" {
		return nil
	}
	if ac, ok := addons[addon]; ok {
		return ac
	}
	// Dynamic fallback: unknown addon, return empty config.
	// Scenarios will handle most extension auto-setting.
	return &addonConfig{}
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
			Country:  l10n.ISOCountryCode(country),
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
	var item itemData
	if r.IntN(2) == 0 && len(locale.Products) > 0 {
		item = pick(r, locale.Products)
	} else {
		item = pick(r, locale.Services)
	}

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

	combo := pickTaxCombo(r, country)
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
		digits := make([]byte, 20)
		for i := range digits {
			digits[i] = byte('0' + r.IntN(10))
		}
		instructions.CreditTransfer = []*pay.CreditTransfer{{
			IBAN: locale.IBANPrefix + string(digits),
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

func pickTaxCombo(_ *rand.Rand, country l10n.TaxCountryCode) *tax.Combo {
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
	return items[r.IntN(len(items))]
}

func newRand(o *Options) *rand.Rand {
	if o.HasSeed {
		return rand.New(rand.NewPCG(uint64(o.Seed), uint64(o.Seed)))
	}
	return rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()+1)))
}
