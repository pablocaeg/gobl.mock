package mock_test

import (
	"testing"

	mock "github.com/pablocaeg/gobl.mock"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// All regimes that should produce valid invoices.
var allRegimes = []l10n.TaxCountryCode{
	"ES", "DE", "MX", "AR", "AT", "BE", "BR", "CA", "CH", "CO",
	"DK", "FR", "GB", "GR", "IE", "IN", "IT", "NL", "PL", "PT",
	"SE", "SG", "US",
}

func TestEnvelope_AllRegimes(t *testing.T) {
	for _, cc := range allRegimes {
		t.Run(string(cc), func(t *testing.T) {
			env, err := mock.Envelope(
				mock.WithRegime(cc),
				mock.WithSeed(42),
			)
			require.NoError(t, err)
			require.NoError(t, env.Validate())
		})
	}
}

func TestEnvelope_AllRegimes_MultiplSeeds(t *testing.T) {
	seeds := []int64{1, 42, 100, 999, 12345}
	for _, cc := range allRegimes {
		for _, seed := range seeds {
			t.Run(string(cc), func(t *testing.T) {
				env, err := mock.Envelope(mock.WithRegime(cc), mock.WithSeed(seed))
				require.NoError(t, err)
				require.NoError(t, env.Validate())
			})
		}
	}
}

func TestEnvelope_CreditNote_AllRegimes(t *testing.T) {
	for _, cc := range allRegimes {
		t.Run(string(cc), func(t *testing.T) {
			env, err := mock.Envelope(
				mock.WithRegime(cc),
				mock.WithCredit(),
				mock.WithSeed(42),
			)
			require.NoError(t, err)
			require.NoError(t, env.Validate())
		})
	}
}

func TestInvoice_ES(t *testing.T) {
	inv, err := mock.Invoice(mock.WithRegime(l10n.ES.Tax()), mock.WithSeed(42))
	require.NoError(t, err)
	assert.Equal(t, l10n.ES.Tax(), inv.GetRegime())
	assert.Equal(t, bill.InvoiceTypeStandard, inv.Type)
	assert.NotEmpty(t, inv.Supplier.Name)
	assert.NotNil(t, inv.Customer)
	assert.Len(t, inv.Lines, 3)
	assert.NotNil(t, inv.Totals)
}

func TestInvoice_Addons(t *testing.T) {
	tests := []struct {
		regime l10n.TaxCountryCode
		addon  cbc.Key
	}{
		{"ES", "es-facturae-v3"},
		{"DE", "de-xrechnung-v3"},
		{"MX", "mx-cfdi-v4"},
	}
	for _, tc := range tests {
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

func TestInvoice_CreditNote(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime(l10n.ES.Tax()),
		mock.WithCredit(),
		mock.WithSeed(42),
	)
	require.NoError(t, err)
	assert.Equal(t, bill.InvoiceTypeCreditNote, inv.Type)
	assert.NotEmpty(t, inv.Preceding)
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
