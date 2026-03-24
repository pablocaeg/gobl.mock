package mock_test

import (
	"testing"

	mock "github.com/pablocaeg/gobl.mock"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoice_ES(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime(l10n.ES.Tax()),
		mock.WithSeed(42),
	)
	require.NoError(t, err)
	assert.Equal(t, l10n.ES.Tax(), inv.GetRegime())
	assert.Equal(t, bill.InvoiceTypeStandard, inv.Type)
	assert.NotEmpty(t, inv.Supplier.Name)
	assert.NotNil(t, inv.Customer)
	assert.Len(t, inv.Lines, 3)
	assert.NotNil(t, inv.Totals)
}

func TestInvoice_DE(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime(l10n.DE.Tax()),
		mock.WithSeed(42),
	)
	require.NoError(t, err)
	assert.Equal(t, l10n.DE.Tax(), inv.GetRegime())
	assert.NotEmpty(t, inv.Supplier.Name)
}

func TestInvoice_MX(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime(l10n.MX.Tax()),
		mock.WithSeed(42),
	)
	require.NoError(t, err)
	assert.Equal(t, l10n.MX.Tax(), inv.GetRegime())
	assert.NotNil(t, inv.Tax)
	assert.NotEmpty(t, inv.Tax.Ext)
}

func TestInvoice_ESFacturae(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime(l10n.ES.Tax()),
		mock.WithAddon("es-facturae-v3"),
		mock.WithSeed(42),
	)
	require.NoError(t, err)
	assert.NotNil(t, inv.Customer)
	assert.NotNil(t, inv.Customer.TaxID)
}

func TestInvoice_DEXRechnung(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime(l10n.DE.Tax()),
		mock.WithAddon("de-xrechnung-v3"),
		mock.WithSeed(42),
	)
	require.NoError(t, err)
	assert.NotEmpty(t, inv.Supplier.People)
	assert.NotEmpty(t, inv.Supplier.Inboxes)
	assert.NotEmpty(t, inv.Supplier.Telephones)
	assert.NotNil(t, inv.Ordering)
}

func TestInvoice_CreditNote(t *testing.T) {
	for _, tc := range []struct {
		name   string
		regime l10n.TaxCountryCode
		addon  cbc.Key
	}{
		{"ES", l10n.ES.Tax(), ""},
		{"DE", l10n.DE.Tax(), ""},
		{"MX", l10n.MX.Tax(), ""},
		{"ES-FacturaE", l10n.ES.Tax(), "es-facturae-v3"},
		{"MX-CFDI", l10n.MX.Tax(), "mx-cfdi-v4"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := []mock.Option{
				mock.WithRegime(tc.regime),
				mock.WithCredit(),
				mock.WithSeed(42),
			}
			if tc.addon != "" {
				opts = append(opts, mock.WithAddon(tc.addon))
			}
			inv, err := mock.Invoice(opts...)
			require.NoError(t, err)
			assert.Equal(t, bill.InvoiceTypeCreditNote, inv.Type)
			assert.NotEmpty(t, inv.Preceding)
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
}

func TestInvoice_LineCount(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime(l10n.ES.Tax()),
		mock.WithLines(15),
		mock.WithSeed(42),
	)
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
	assert.Equal(t, inv1.Lines[0].Item.Name, inv2.Lines[0].Item.Name)
}

func TestEnvelope(t *testing.T) {
	for _, regime := range []l10n.TaxCountryCode{
		l10n.ES.Tax(),
		l10n.DE.Tax(),
		l10n.MX.Tax(),
	} {
		t.Run(string(regime), func(t *testing.T) {
			env, err := mock.Envelope(
				mock.WithRegime(regime),
				mock.WithSeed(42),
			)
			require.NoError(t, err)
			assert.NotNil(t, env.Head)
			assert.NotNil(t, env.Document)

			// Envelope must be valid.
			require.NoError(t, env.Validate())
		})
	}
}

func TestInvoice_UnsupportedRegime(t *testing.T) {
	_, err := mock.Invoice(mock.WithRegime(l10n.TaxCountryCode("ZZ")))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported regime")
}
