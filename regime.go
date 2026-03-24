package mock

import (
	"math/rand/v2"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

type cityData struct {
	Name   string
	Region string
	Code   cbc.Code // postal code
}

type itemData struct {
	Name string
	// Price as a string amount (e.g. "100.00").
	Price string
	// Unit is optional (e.g. "h", "service").
	Unit string
}

type taxRate struct {
	Category cbc.Code
	Rate     cbc.Key
}

// regimeConfig holds all the locale-specific data needed to generate
// realistic invoices for a given tax regime.
type regimeConfig struct {
	Country  l10n.TaxCountryCode
	Currency currency.Code

	SupplierNames []string
	CustomerNames []string
	Cities        []cityData
	Streets       []string

	Products []itemData
	Services []itemData

	// TaxRates for line items, first is default.
	TaxRates []taxRate

	// PaymentKey is the default payment means key.
	PaymentKey cbc.Key

	// IBANPrefix for generating plausible IBAN-like references.
	IBANPrefix string

	// SupplierExt returns addon-specific extensions for the supplier.
	SupplierExt func(r *rand.Rand) tax.Extensions
	// CustomerExt returns addon-specific extensions for the customer.
	CustomerExt func(r *rand.Rand) tax.Extensions
	// ItemExt returns addon-specific extensions for a line item.
	ItemExt func(r *rand.Rand) tax.Extensions
	// InvoiceTaxExt returns addon-specific extensions for the invoice tax object.
	InvoiceTaxExt func(r *rand.Rand) tax.Extensions
	// PaymentExt returns addon-specific extensions for payment instructions.
	PaymentExt func(r *rand.Rand) tax.Extensions

	// CustomerPostalCode returns a postal code for the customer address
	// when the addon requires it (e.g. CFDI).
	CustomerPostalCode func(r *rand.Rand) cbc.Code

	// DefaultAddon is the addon automatically applied when none is specified.
	DefaultAddon cbc.Key

	// SupplierPeople indicates the addon requires people on the supplier.
	SupplierPeople bool
	// SupplierInboxes indicates the addon requires inboxes on the supplier.
	SupplierInboxes bool
	// SupplierPhones indicates the addon requires telephones on the supplier.
	SupplierPhones bool
	// CustomerInboxes indicates the addon requires inboxes on the customer.
	CustomerInboxes bool
	// RequiresOrdering indicates the addon requires ordering with a code.
	RequiresOrdering bool
}

var regimeConfigs = map[l10n.TaxCountryCode]*regimeConfig{}

func getRegimeConfig(country l10n.TaxCountryCode, addon cbc.Key) *regimeConfig {
	cfg, ok := regimeConfigs[country]
	if !ok {
		return nil
	}
	// Apply addon overrides.
	return applyAddonConfig(cfg, addon)
}

// applyAddonConfig returns a copy of the config with addon-specific overrides.
func applyAddonConfig(base *regimeConfig, addon cbc.Key) *regimeConfig {
	if addon == "" {
		return base
	}
	fn, ok := addonConfigs[addon]
	if !ok {
		return base
	}
	return fn(base)
}

type addonConfigFunc func(base *regimeConfig) *regimeConfig

var addonConfigs = map[cbc.Key]addonConfigFunc{}

// pick returns a random element from a slice.
func pick[T any](r *rand.Rand, items []T) T {
	return items[r.IntN(len(items))]
}

func init() {
	regimeConfigs[l10n.ES.Tax()] = esConfig()
	regimeConfigs[l10n.DE.Tax()] = deConfig()
	regimeConfigs[l10n.MX.Tax()] = mxConfig()

	addonConfigs["es-facturae-v3"] = esFacturaeConfig
	addonConfigs["de-xrechnung-v3"] = deXRechnungConfig
	addonConfigs["mx-cfdi-v4"] = mxCFDIConfig
}

func (cfg *regimeConfig) paymentKey() cbc.Key {
	if cfg.PaymentKey != "" {
		return cfg.PaymentKey
	}
	return pay.MeansKeyCreditTransfer
}
