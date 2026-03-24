package mock_test

import (
	"fmt"
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

// TestEnvelope_AllRegimes validates every GOBL regime using the
// auto-generated regime list from go generate.
func TestEnvelope_AllRegimes(t *testing.T) {
	for _, cc := range generatedRegimes {
		t.Run(string(cc), func(t *testing.T) {
			env, err := mock.Envelope(mock.WithRegime(cc), mock.WithSeed(42))
			require.NoError(t, err)
			require.NoError(t, env.Validate())
		})
	}
}

// TestEnvelope_AllRegimes_MultipleSeeds verifies that different random
// seeds all produce valid invoices for every regime.
func TestEnvelope_AllRegimes_MultipleSeeds(t *testing.T) {
	for _, cc := range generatedRegimes {
		for _, seed := range []int64{1, 42, 100, 999, 12345} {
			t.Run(fmt.Sprintf("%s/seed=%d", cc, seed), func(t *testing.T) {
				env, err := mock.Envelope(mock.WithRegime(cc), mock.WithSeed(seed))
				require.NoError(t, err)
				require.NoError(t, env.Validate())
			})
		}
	}
}

// TestEnvelope_AllAddons validates every GOBL addon using the
// auto-generated addon list from go generate.
func TestEnvelope_AllAddons(t *testing.T) {
	for _, tc := range generatedAddons {
		t.Run(string(tc.Regime)+"+"+string(tc.Key), func(t *testing.T) {
			env, err := mock.Envelope(
				mock.WithRegime(tc.Regime),
				mock.WithAddon(tc.Key),
				mock.WithSeed(42),
			)
			require.NoError(t, err)
			require.NoError(t, env.Validate())
		})
	}
}

// TestEnvelope_CreditNote_AllRegimes verifies credit notes for every regime.
func TestEnvelope_CreditNote_AllRegimes(t *testing.T) {
	for _, cc := range generatedRegimes {
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

func TestInvoice_WithCredit(t *testing.T) {
	inv, err := mock.Invoice(mock.WithRegime(l10n.ES.Tax()), mock.WithCredit(), mock.WithSeed(42))
	require.NoError(t, err)
	assert.Equal(t, bill.InvoiceTypeCreditNote, inv.Type)
	assert.NotEmpty(t, inv.Preceding)
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
	assert.ErrorContains(t, err, "unsupported regime")
}
