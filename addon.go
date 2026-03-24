package mock

import (
	"fmt"
	"math/rand/v2"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

// addonConfig holds addon-specific generation requirements.
type addonConfig struct {
	SupplierExt        func(r *rand.Rand) tax.Extensions
	CustomerExt        func(r *rand.Rand) tax.Extensions
	ItemExt            func(r *rand.Rand) tax.Extensions
	InvoiceTaxExt      func(r *rand.Rand) tax.Extensions
	CustomerPostalCode func(r *rand.Rand) cbc.Code
	CorrectionExt      func(r *rand.Rand) tax.Extensions
	CorrectionStamps   func(r *rand.Rand) []stampData
	SupplierPeople     bool
	SupplierInboxes    bool
	SupplierPhones     bool
	CustomerInboxes    bool
	RequiresOrdering   bool
}

type stampData struct {
	Provider cbc.Key
	Value    string
}

var addons = map[cbc.Key]*addonConfig{
	"es-facturae-v3": {
		CorrectionExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"es-facturae-correction": "01"}
		},
	},
	"de-xrechnung-v3": {
		SupplierPeople:   true,
		SupplierInboxes:  true,
		SupplierPhones:   true,
		CustomerInboxes:  true,
		RequiresOrdering: true,
	},
	"mx-cfdi-v4": {
		SupplierExt: func(r *rand.Rand) tax.Extensions {
			return tax.Extensions{"mx-cfdi-fiscal-regime": cbc.Code(pick(r, mxFiscalRegimes))}
		},
		CustomerExt: func(r *rand.Rand) tax.Extensions {
			return tax.Extensions{
				"mx-cfdi-fiscal-regime": cbc.Code(pick(r, mxFiscalRegimes)),
				"mx-cfdi-use":           cbc.Code(pick(r, mxCFDIUses)),
			}
		},
		ItemExt: func(r *rand.Rand) tax.Extensions {
			return tax.Extensions{"mx-cfdi-prod-serv": cbc.Code(pick(r, mxProdServCodes))}
		},
		InvoiceTaxExt: func(r *rand.Rand) tax.Extensions {
			return tax.Extensions{"mx-cfdi-issue-place": cbc.Code(pick(r, mxPostalCodes))}
		},
		CustomerPostalCode: func(r *rand.Rand) cbc.Code {
			return cbc.Code(pick(r, mxPostalCodes))
		},
		CorrectionStamps: func(_ *rand.Rand) []stampData {
			return []stampData{{Provider: "sat-uuid", Value: fmt.Sprintf("a1b2c3d4-e5f6-7890-abcd-%012d", 1)}}
		},
	},
}

// CFDI-specific data tables.
var (
	mxFiscalRegimes = []string{"601", "603", "612", "620", "621", "626"}
	mxCFDIUses      = []string{"G01", "G02", "G03", "S01"}
	mxProdServCodes = []string{"10101500", "43211500", "80111600", "84111500", "81112100", "82101500"}
	mxPostalCodes   = []string{"06600", "01000", "44100", "64000", "72000", "76000"}
)

func defaultAddonForRegime(country l10n.TaxCountryCode) cbc.Key {
	defaults := map[l10n.TaxCountryCode]cbc.Key{
		"MX": "mx-cfdi-v4",
	}
	return defaults[country]
}
