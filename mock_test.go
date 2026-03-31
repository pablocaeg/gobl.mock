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

// TestEnvelope_Scenarios validates each scenario against key regimes.
func TestEnvelope_Scenarios(t *testing.T) {
	tests := []struct {
		scenario cbc.Key
		regime   l10n.TaxCountryCode
	}{
		{mock.ScenarioHotel, "ES"},
		{mock.ScenarioHotel, "IT"},
		{mock.ScenarioHotel, "FR"},
		{mock.ScenarioHotel, "DE"},
		{mock.ScenarioHotel, "GB"},
		{mock.ScenarioFreelance, "ES"},
		{mock.ScenarioFreelance, "IT"},
		{mock.ScenarioFreelance, "FR"},
		{mock.ScenarioFreelance, "DE"},
		{mock.ScenarioReverseCharge, "ES"},
		{mock.ScenarioReverseCharge, "IT"},
		{mock.ScenarioReverseCharge, "DE"},
		{mock.ScenarioReverseCharge, "FR"},
		{mock.ScenarioRestaurant, "ES"},
		{mock.ScenarioRestaurant, "IT"},
		{mock.ScenarioRestaurant, "FR"},
		{mock.ScenarioRestaurant, "DE"},
		{mock.ScenarioEcommerce, "ES"},
		{mock.ScenarioEcommerce, "DE"},
		{mock.ScenarioEcommerce, "FR"},
		{mock.ScenarioEcommerce, "GB"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%s", tt.regime, tt.scenario), func(t *testing.T) {
			env, err := mock.Envelope(
				mock.WithRegime(tt.regime),
				mock.WithScenario(tt.scenario),
				mock.WithSeed(42),
			)
			require.NoError(t, err)
			require.NoError(t, env.Validate())
		})
	}
}

// TestEnvelope_Scenario_AllRegimes verifies that every scenario works with
// every regime (using fallback items for regimes without overrides).
func TestEnvelope_Scenario_AllRegimes(t *testing.T) {
	scenarios := []cbc.Key{mock.ScenarioHotel, mock.ScenarioFreelance, mock.ScenarioReverseCharge, mock.ScenarioRestaurant, mock.ScenarioEcommerce}
	for _, sc := range scenarios {
		for _, cc := range generatedRegimes {
			t.Run(fmt.Sprintf("%s/%s", cc, sc), func(t *testing.T) {
				env, err := mock.Envelope(
					mock.WithRegime(cc),
					mock.WithScenario(sc),
					mock.WithSeed(42),
				)
				require.NoError(t, err)
				require.NoError(t, env.Validate())
			})
		}
	}
}

// TestEnvelope_ScenarioWithAddon validates scenarios compose with addons.
func TestEnvelope_ScenarioWithAddon(t *testing.T) {
	tests := []struct {
		scenario cbc.Key
		regime   l10n.TaxCountryCode
		addon    cbc.Key
	}{
		{mock.ScenarioHotel, "IT", "it-sdi-v1"},
		{mock.ScenarioHotel, "ES", "es-facturae-v3"},
		{mock.ScenarioHotel, "DE", "de-xrechnung-v3"},
		{mock.ScenarioHotel, "AT", "eu-en16931-v2017"},
		{mock.ScenarioFreelance, "IT", "it-sdi-v1"},
		{mock.ScenarioFreelance, "ES", "es-facturae-v3"},
		{mock.ScenarioFreelance, "MX", "mx-cfdi-v4"},
		{mock.ScenarioFreelance, "CO", "co-dian-v2"},
		{mock.ScenarioReverseCharge, "ES", "es-facturae-v3"},
		{mock.ScenarioReverseCharge, "DE", "de-xrechnung-v3"},
		{mock.ScenarioReverseCharge, "IT", "it-sdi-v1"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%s+%s", tt.regime, tt.scenario, tt.addon), func(t *testing.T) {
			env, err := mock.Envelope(
				mock.WithRegime(tt.regime),
				mock.WithAddon(tt.addon),
				mock.WithScenario(tt.scenario),
				mock.WithSeed(42),
			)
			require.NoError(t, err)
			require.NoError(t, env.Validate())
		})
	}
}

// TestEnvelope_ScenarioWithCreditNote validates scenarios compose with invoice types.
func TestEnvelope_ScenarioWithCreditNote(t *testing.T) {
	scenarios := []cbc.Key{mock.ScenarioHotel, mock.ScenarioFreelance}
	for _, sc := range scenarios {
		t.Run(string(sc), func(t *testing.T) {
			env, err := mock.Envelope(
				mock.WithRegime(l10n.ES.Tax()),
				mock.WithScenario(sc),
				mock.WithCredit(),
				mock.WithSeed(42),
			)
			require.NoError(t, err)
			require.NoError(t, env.Validate())
		})
	}
}

// TestInvoice_ScenarioHotel verifies hotel scenario produces domain-specific content.
func TestInvoice_ScenarioHotel(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime("IT"),
		mock.WithScenario(mock.ScenarioHotel),
		mock.WithSeed(42),
	)
	require.NoError(t, err)

	// Should have hotel-specific items.
	hotelItems := map[string]bool{"Camera Matrimoniale": true, "Tassa di Soggiorno": true}
	for _, line := range inv.Lines {
		assert.True(t, hotelItems[line.Item.Name], "unexpected item: %s", line.Item.Name)
	}
}

// TestInvoice_ScenarioFreelance verifies freelance scenario has people and retained taxes.
func TestInvoice_ScenarioFreelance(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime(l10n.ES.Tax()),
		mock.WithScenario(mock.ScenarioFreelance),
		mock.WithSeed(42),
	)
	require.NoError(t, err)

	// Supplier should have a person.
	require.NotEmpty(t, inv.Supplier.People)
	assert.NotEmpty(t, inv.Supplier.People[0].Name.Given)

	// Lines should include IRPF retained tax.
	hasIRPF := false
	for _, line := range inv.Lines {
		for _, tc := range line.Taxes {
			if tc.Category == "IRPF" {
				hasIRPF = true
			}
		}
	}
	assert.True(t, hasIRPF, "freelance ES should have IRPF retained tax")

	// Payment should have due date terms.
	require.NotNil(t, inv.Payment)
	require.NotNil(t, inv.Payment.Terms)
	assert.Equal(t, cbc.Key("due-date"), inv.Payment.Terms.Key)
}

// TestInvoice_ScenarioReverseCharge verifies reverse-charge tag and cross-border customer.
func TestInvoice_ScenarioReverseCharge(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime(l10n.ES.Tax()),
		mock.WithScenario(mock.ScenarioReverseCharge),
		mock.WithSeed(42),
	)
	require.NoError(t, err)

	assert.True(t, inv.HasTags(tax.TagReverseCharge))
	// Customer should be from a different country.
	assert.NotEqual(t, l10n.ES.Tax(), inv.Customer.TaxID.Country)
}

// TestInvoice_ScenarioFreelanceMX verifies Mexican freelance has ISR and retained VAT.
func TestInvoice_ScenarioFreelanceMX(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime("MX"),
		mock.WithScenario(mock.ScenarioFreelance),
		mock.WithSeed(42),
	)
	require.NoError(t, err)

	hasISR := false
	hasRVAT := false
	for _, line := range inv.Lines {
		for _, tc := range line.Taxes {
			if tc.Category == "ISR" {
				hasISR = true
			}
			if tc.Category == "RVAT" {
				hasRVAT = true
			}
		}
	}
	assert.True(t, hasISR, "MX freelance should have ISR retained tax")
	assert.True(t, hasRVAT, "MX freelance should have RVAT retained tax")
}

// TestInvoice_ScenarioFreelanceCO verifies Colombian freelance has ReteRenta.
func TestInvoice_ScenarioFreelanceCO(t *testing.T) {
	inv, err := mock.Invoice(
		mock.WithRegime("CO"),
		mock.WithScenario(mock.ScenarioFreelance),
		mock.WithSeed(42),
	)
	require.NoError(t, err)

	hasRR := false
	for _, line := range inv.Lines {
		for _, tc := range line.Taxes {
			if tc.Category == "RR" {
				hasRR = true
			}
		}
	}
	assert.True(t, hasRR, "CO freelance should have ReteRenta retained tax")
}

func TestInvoice_UnknownScenario(t *testing.T) {
	_, err := mock.Invoice(
		mock.WithRegime(l10n.ES.Tax()),
		mock.WithScenario("nonexistent"),
		mock.WithSeed(42),
	)
	assert.ErrorContains(t, err, "unknown scenario")
}
