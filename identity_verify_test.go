package mock

import (
	"math/rand/v2"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/require"

	_ "github.com/invopop/gobl/addons"
	_ "github.com/invopop/gobl/regimes"
)

// TestTaxIDValidation generates 50 tax IDs per regime and verifies each one
// passes GOBL's full validation pipeline by building a minimal invoice.
// This is the definitive proof that our generators produce correct check digits.
func TestTaxIDValidation(t *testing.T) {
	regimes := []l10n.TaxCountryCode{
		"AR", "AT", "BE", "BR", "CH", "CO", "DE", "DK",
		"EL", "ES", "FR", "GB", "GR", "IE", "IN", "IT",
		"MX", "NL", "PL", "PT", "SE", "SG",
	}

	for _, cc := range regimes {
		t.Run(string(cc), func(t *testing.T) {
			for seed := int64(0); seed < 50; seed++ {
				r := rand.New(rand.NewPCG(uint64(seed), uint64(seed)))
				code := generateTaxID(r, cc, true)
				if code == "" {
					continue // some regimes don't require codes
				}

				err := validateTaxIDViaInvoice(cc, code)
				require.NoError(t, err, "seed=%d code=%s", seed, code)
			}
		})
	}
}

// TestTaxIDValidation_Individuals verifies generated individual tax IDs
// for regimes that distinguish between org and person (ES, BR, MX, AR).
func TestTaxIDValidation_Individuals(t *testing.T) {
	regimes := []l10n.TaxCountryCode{"ES", "BR", "MX", "AR"}

	for _, cc := range regimes {
		t.Run(string(cc), func(t *testing.T) {
			for seed := int64(0); seed < 50; seed++ {
				r := rand.New(rand.NewPCG(uint64(seed), uint64(seed)))
				code := generateTaxID(r, cc, false) // individual
				if code == "" {
					continue
				}

				err := validateTaxIDViaInvoice(cc, code)
				require.NoError(t, err, "seed=%d code=%s", seed, code)
			}
		})
	}
}

// validateTaxIDViaInvoice builds a minimal invoice with the given tax ID
// and runs it through GOBL's full pipeline. If it succeeds, the tax ID
// passed all regime-specific validation including check digit verification.
func validateTaxIDViaInvoice(country l10n.TaxCountryCode, code cbc.Code) error {
	price := num.MakeAmount(100, 2)
	inv := &bill.Invoice{
		Series: "TEST",
		Code:   "001",
		Supplier: &org.Party{
			Name: "Test",
			TaxID: &tax.Identity{
				Country: country,
				Code:    code,
			},
		},
		Lines: []*bill.Line{{
			Quantity: num.MakeAmount(1, 0),
			Item:     &org.Item{Name: "Test", Price: &price},
			Taxes:    tax.Set{pickTaxCombo(nil, country)},
		}},
	}

	_, err := gobl.Envelop(inv)
	return err
}
