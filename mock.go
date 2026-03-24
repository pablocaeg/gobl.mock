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

	_ "github.com/invopop/gobl/addons"
	_ "github.com/invopop/gobl/regimes"
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
	ac := addons[addon]

	// Verify the regime exists.
	if tax.RegimeDefFor(l10n.Code(country)) == nil {
		return nil, fmt.Errorf("unsupported regime: %s", country)
	}

	inv := &bill.Invoice{
		Series:    "MOCK",
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

	// Invoice-level extensions from addon.
	if ac != nil && ac.InvoiceTaxExt != nil {
		inv.Tax = &bill.Tax{Ext: ac.InvoiceTaxExt(r)}
	}

	inv.Supplier = buildParty(r, country, locale, locale.SupplierNames, ac, true)
	if !o.Simplified {
		inv.Customer = buildParty(r, country, locale, locale.CustomerNames, ac, false)
	}
	inv.Lines = buildLines(r, country, locale, ac, o.Lines)
	inv.Payment = buildPayment(r, country, locale)

	if ac != nil && ac.RequiresOrdering {
		inv.Ordering = &bill.Ordering{
			Code: cbc.Code(fmt.Sprintf("PO-%05d", r.IntN(99999)+1)),
		}
	}

	if o.Credit {
		applyCreditNote(r, inv, addon, ac)
	}

	env, err := gobl.Envelop(inv)
	if err != nil {
		return nil, fmt.Errorf("building envelope: %w", err)
	}
	return env, nil
}

func buildParty(r *rand.Rand, country l10n.TaxCountryCode, locale *localeData, names []string, ac *addonConfig, isSupplier bool) *org.Party {
	city := pick(r, locale.Cities)
	postalCode := city.Code
	// Regime-specific postal code formats.
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
		if ac.CustomerPostalCode != nil {
			p.Addresses[0].Code = ac.CustomerPostalCode(r)
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
	line := &bill.Line{
		Quantity: num.MakeAmount(int64(r.IntN(20)+1), 0),
		Item:     lineItem,
		Taxes:    tax.Set{combo},
	}

	if r.IntN(5) == 0 {
		pct, _ := num.PercentageFromString("10%")
		line.Discounts = []*bill.LineDiscount{{Percent: &pct, Reason: "Discount"}}
	}
	return line
}

func buildPayment(r *rand.Rand, country l10n.TaxCountryCode, locale *localeData) *bill.PaymentDetails {
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

	return &bill.PaymentDetails{
		Instructions: instructions,
		Terms:        &pay.Terms{Key: pay.TermKeyInstant, Notes: "Payment due upon receipt."},
	}
}

func applyCreditNote(r *rand.Rand, inv *bill.Invoice, addon cbc.Key, ac *addonConfig) {
	inv.Type = bill.InvoiceTypeCreditNote
	yesterday := cal.Today().Add(0, 0, -1)

	ref := &org.DocumentRef{
		Identify:  uuid.Identify{UUID: uuid.V7()},
		Type:      bill.InvoiceTypeStandard,
		Series:    "MOCK",
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

// pickTaxCombo selects a valid tax combo for the regime by reading the
// regime definition dynamically. Returns either a rate-based combo (for
// regimes with defined rates) or a percent-based combo (for regimes
// where rates are not enumerated).
func pickTaxCombo(_ *rand.Rand, country l10n.TaxCountryCode) *tax.Combo {
	regime := tax.RegimeDefFor(l10n.Code(country))
	if regime == nil {
		return &tax.Combo{Category: tax.CategoryVAT, Rate: tax.KeyStandard}
	}

	// Find the primary tax category: prefer VAT/GST/ST, skip retained/informative.
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

	// For VAT, the "standard" rate key is specially normalized to "general" by GOBL.
	if primary.Code == tax.CategoryVAT {
		return &tax.Combo{Category: primary.Code, Rate: tax.KeyStandard}
	}

	// For other categories (GST, ST, etc.), find the first rate with values and use it directly.
	for _, rate := range primary.Rates {
		if len(rate.Values) > 0 && !rate.Values[0].Disabled {
			return &tax.Combo{Category: primary.Code, Rate: rate.Rate}
		}
	}

	// No usable rates — provide an explicit percent.
	pct := num.MakePercentage(10, 2)
	return &tax.Combo{Category: primary.Code, Percent: &pct}
}

// Regime-specific postal code generators for countries that validate the format.
var postalCodeFormats = map[l10n.TaxCountryCode]func(r *rand.Rand) cbc.Code{
	"BR": func(r *rand.Rand) cbc.Code {
		return cbc.Code(fmt.Sprintf("%05d-%03d", r.IntN(99999), r.IntN(999)))
	},
}

// Regime-specific item identity requirements.
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
