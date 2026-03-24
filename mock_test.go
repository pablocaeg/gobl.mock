package mock_test

import (
	"testing"

	mock "github.com/pablocaeg/gobl.mock"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var allRegimes = []l10n.TaxCountryCode{
	"AE", "AR", "AT", "BE", "BR", "CA", "CH", "CO", "DE", "DK",
	"ES", "FR", "GB", "GR", "IE", "IN", "IT", "MX", "NL", "PL",
	"PT", "SE", "SG", "US",
}

var allAddons = []struct {
	regime l10n.TaxCountryCode
	addon  cbc.Key
}{
	{"AR", "ar-arca-v4"}, {"BR", "br-nfe-v4"}, {"BR", "br-nfse-v1"},
	{"CO", "co-dian-v2"}, {"DE", "de-xrechnung-v3"}, {"DE", "de-zugferd-v2"},
	{"ES", "es-facturae-v3"}, {"ES", "es-sii-v1"}, {"ES", "es-tbai-v1"}, {"ES", "es-verifactu-v1"},
	{"FR", "fr-choruspro-v1"}, {"FR", "fr-ctc-flow2-v1"}, {"FR", "fr-facturx-v1"},
	{"DE", "eu-en16931-v2017"},
	{"GR", "gr-mydata-v1"}, {"IT", "it-sdi-v1"}, {"IT", "it-ticket-v1"},
	{"MX", "mx-cfdi-v4"}, {"PL", "pl-favat-v1"}, {"PT", "pt-saft-v1"},
}

func TestEnvelope_AllRegimes(t *testing.T) {
	for _, cc := range allRegimes {
		t.Run(string(cc), func(t *testing.T) {
			env, err := mock.Envelope(mock.WithRegime(cc), mock.WithSeed(42))
			require.NoError(t, err)
			require.NoError(t, env.Validate())
		})
	}
}

func TestEnvelope_AllRegimes_MultipleSeeds(t *testing.T) {
	for _, cc := range allRegimes {
		for _, seed := range []int64{1, 42, 100, 999, 12345} {
			t.Run(string(cc), func(t *testing.T) {
				env, err := mock.Envelope(mock.WithRegime(cc), mock.WithSeed(seed))
				require.NoError(t, err)
				require.NoError(t, env.Validate())
			})
		}
	}
}

func TestEnvelope_AllAddons(t *testing.T) {
	for _, tc := range allAddons {
		t.Run(string(tc.regime)+"+"+string(tc.addon), func(t *testing.T) {
			env, err := mock.Envelope(
				mock.WithRegime(tc.regime),
				mock.WithAddon(tc.addon),
				mock.WithSeed(42),
			)
			require.NoError(t, err)
			require.NoError(t, env.Validate())
		})
	}
}

func TestEnvelope_CreditNote_AllRegimes(t *testing.T) {
	for _, cc := range allRegimes {
		t.Run(string(cc), func(t *testing.T) {
			env, err := mock.Envelope(
				mock.WithRegime(cc),
				mock.WithType(bill.InvoiceTypeCreditNote),
				mock.WithSeed(42),
			)
			require.NoError(t, err)
			require.NoError(t, env.Validate())
		})
	}
}

func TestInvoice_Types(t *testing.T) {
	types := []cbc.Key{
		bill.InvoiceTypeCreditNote,
		bill.InvoiceTypeCorrective,
		bill.InvoiceTypeDebitNote,
		bill.InvoiceTypeProforma,
	}
	for _, invType := range types {
		t.Run(string(invType), func(t *testing.T) {
			inv, err := mock.Invoice(
				mock.WithRegime(l10n.ES.Tax()),
				mock.WithType(invType),
				mock.WithSeed(42),
			)
			require.NoError(t, err)
			assert.Equal(t, invType, inv.Type)
			if invType != bill.InvoiceTypeProforma {
				assert.NotEmpty(t, inv.Preceding)
			}
		})
	}
}

func TestInvoice_Template(t *testing.T) {
	price := num.MakeAmount(999, 2)
	tmpl := &bill.Invoice{
		Customer: &org.Party{
			Name: "Custom Customer S.L.",
			TaxID: &tax.Identity{
				Country: l10n.ES.Tax(),
				Code:    "54387763P",
			},
		},
		Lines: []*bill.Line{{
			Quantity: num.MakeAmount(1, 0),
			Item:     &org.Item{Name: "Custom item", Price: &price},
			Taxes:    tax.Set{{Category: tax.CategoryVAT, Rate: tax.KeyStandard}},
		}},
	}

	inv, err := mock.Invoice(
		mock.WithRegime(l10n.ES.Tax()),
		mock.WithTemplate(tmpl),
		mock.WithSeed(42),
	)
	require.NoError(t, err)
	assert.Equal(t, "Custom Customer S.L.", inv.Customer.Name)
	assert.Len(t, inv.Lines, 1)
	assert.Equal(t, "Custom item", inv.Lines[0].Item.Name)
	// Supplier should still be generated.
	assert.NotEmpty(t, inv.Supplier.Name)
}

func TestInvoice_Simplified(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime(l10n.ES.Tax()),
		mock.WithSimplified(),
		mock.WithSeed(42),
	)
	require.NoError(t, err)
	assert.Nil(t, inv.Customer)
	assert.True(t, inv.HasTags(tax.TagSimplified))
}

func TestInvoice_Lines(t *testing.T) {
	inv, err := mock.Invoice(mock.WithRegime(l10n.ES.Tax()), mock.WithLines(15), mock.WithSeed(42))
	require.NoError(t, err)
	assert.Len(t, inv.Lines, 15)
}

func TestInvoice_Seed(t *testing.T) {
	inv1, err := mock.Invoice(mock.WithRegime(l10n.ES.Tax()), mock.WithSeed(42))
	require.NoError(t, err)
	inv2, err := mock.Invoice(mock.WithRegime(l10n.ES.Tax()), mock.WithSeed(42))
	require.NoError(t, err)
	assert.Equal(t, inv1.Supplier.Name, inv2.Supplier.Name)
	assert.Equal(t, inv1.Supplier.TaxID.Code, inv2.Supplier.TaxID.Code)
}

func TestInvoice_UnsupportedRegime(t *testing.T) {
	_, err := mock.Invoice(mock.WithRegime("ZZ"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported regime")
}
