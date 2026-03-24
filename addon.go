package mock

import (
	"fmt"
	"math/rand/v2"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// addonConfig holds addon-specific generation requirements that cannot be
// derived dynamically from the GOBL addon definition.
type addonConfig struct {
	// Extension generators per object.
	SupplierExt   func(r *rand.Rand) tax.Extensions
	CustomerExt   func(r *rand.Rand) tax.Extensions
	ItemExt       func(r *rand.Rand) tax.Extensions
	InvoiceTaxExt func(r *rand.Rand) tax.Extensions
	PaymentExt    func(r *rand.Rand) tax.Extensions
	ComboExt      func(r *rand.Rand) tax.Extensions

	// Address overrides.
	CustomerPostalCode func(r *rand.Rand) cbc.Code

	// Correction handling.
	CorrectionExt    func(r *rand.Rand) tax.Extensions
	CorrectionStamps func(r *rand.Rand) []stampData

	// Extra tax combos added to every line (e.g. BR NF-e requires ICMS+PIS+COFINS).
	ExtraCombos func(r *rand.Rand) tax.Set

	// Item identities (e.g. BR NFS-e requires service code).
	ItemIdentities func(r *rand.Rand) []*org.Identity

	// Series override (some addons require specific series formats).
	Series string

	// Notes to add (some addons require specific note keys).
	Notes []*org.Note

	// Party identities (e.g. FR Chorus Pro requires SIREN/SIRET).
	SupplierIdentities func(r *rand.Rand) []*org.Identity
	CustomerIdentities func(r *rand.Rand) []*org.Identity

	// Address state override (e.g. BR requires state code).
	SupplierState cbc.Code
	CustomerState cbc.Code

	// Structural requirements.
	SupplierPeople   bool
	SupplierInboxes  bool
	SupplierPhones   bool
	CustomerInboxes  bool
	RequiresOrdering bool
	NumericSeries    bool // series must be a number
	RequiresDueDates bool // payment terms need due_dates
}

type stampData struct {
	Provider cbc.Key
	Value    string
}

var addons = map[cbc.Key]*addonConfig{
	// --- Spain ---
	"es-facturae-v3": {
		CorrectionExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"es-facturae-correction": "01"}
		},
	},
	"es-sii-v1":       {},
	"es-verifactu-v1": {},
	"es-tbai-v1": {
		InvoiceTaxExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"es-tbai-region": "BI"} // Bizkaia
		},
		Notes: []*org.Note{{Key: cbc.Key("general"), Text: "Test invoice"}},
	},

	// --- Germany ---
	"de-xrechnung-v3": {
		SupplierPeople:   true,
		SupplierInboxes:  true,
		SupplierPhones:   true,
		CustomerInboxes:  true,
		RequiresOrdering: true,
	},
	"de-zugferd-v2": {},

	// --- France ---
	"fr-facturx-v1": {},
	"fr-choruspro-v1": {
		SupplierExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"fr-choruspro-scheme": "1"}
		},
		CustomerExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"fr-choruspro-scheme": "1"}
		},
		CustomerIdentities: func(_ *rand.Rand) []*org.Identity {
			return []*org.Identity{{Type: "SIRET", Code: "73282932000074"}}
		},
		RequiresOrdering: true,
	},

	// --- Italy ---
	"it-sdi-v1": {},

	// --- Mexico ---
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

	// --- Argentina ---
	"ar-arca-v4": {
		NumericSeries:    true,
		RequiresOrdering: true, // ordering with period required
		RequiresDueDates: true,
		CustomerExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"ar-arca-vat-status": "1"}
		},
	},

	// --- Brazil ---
	"br-nfe-v4": {
		Series: "1",
		Notes:  []*org.Note{{Key: "reason", Text: "Venda de mercadoria"}},
		InvoiceTaxExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"br-nfe-presence": "1"}
		},
		SupplierExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"br-ibge-municipality": "3550308"}
		},
		SupplierIdentities: func(_ *rand.Rand) []*org.Identity {
			return []*org.Identity{{Key: "br-nfe-state-reg", Code: "110042490114"}}
		},
		CustomerExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"br-ibge-municipality": "3550308"}
		},
		ItemExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"br-nfe-cfop": "5102"}
		},
		PaymentExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"br-nfe-payment-means": "03"}
		},
		SupplierState: "SP",
		CustomerState: "SP",
		ExtraCombos: func(_ *rand.Rand) tax.Set {
			icms := num.MakePercentage(18, 2)
			pis := num.MakePercentage(165, 4)
			cofins := num.MakePercentage(760, 4)
			return tax.Set{
				{Category: "ICMS", Percent: &icms},
				{Category: "PIS", Percent: &pis},
				{Category: "COFINS", Percent: &cofins},
			}
		},
	},
	"br-nfse-v1": {
		SupplierExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{
				"br-nfse-simples":          "2",
				"br-nfse-fiscal-incentive": "2",
				"br-ibge-municipality":     "3550308",
			}
		},
		SupplierIdentities: func(_ *rand.Rand) []*org.Identity {
			return []*org.Identity{{Key: "br-nfse-municipal-reg", Code: "12345678"}}
		},
		SupplierState: "SP",
		ItemExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"br-nfse-service": "01.07"}
		},
		ExtraCombos: func(_ *rand.Rand) tax.Set {
			iss := num.MakePercentage(5, 2)
			return tax.Set{
				{Category: "ISS", Percent: &iss, Ext: tax.Extensions{"br-nfse-iss-liability": "1"}},
			}
		},
	},

	// --- Colombia ---
	"co-dian-v2": {
		SupplierExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"co-dian-municipality": "11001"} // Bogota
		},
		CustomerExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"co-dian-municipality": "11001"}
		},
	},

	// --- EU ---
	"eu-en16931-v2017": {},

	// --- Greece ---
	"gr-mydata-v1": {
		ComboExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{
				"gr-mydata-vat-rate":    "1", // Standard 24%
				"gr-mydata-income-cat":  "category1_1",
				"gr-mydata-income-type": "E3_561_001",
			}
		},
		PaymentExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"gr-mydata-payment-means": "3"} // Bank transfer
		},
	},

	// --- Poland ---
	"pl-favat-v1": {},

	// --- Portugal ---
	"pt-saft-v1": {
		Series: "FT MOCK",
		ItemExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"pt-saft-product-type": "S"} // Service
		},
		ComboExt: func(_ *rand.Rand) tax.Extensions {
			return tax.Extensions{"pt-saft-tax-rate": "NOR"} // Normal rate
		},
	},
}

// Data tables for specific addons.
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
