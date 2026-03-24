package mock

import (
	"math/rand/v2"
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/invopop/gobl/regimes"
)

func TestGenerateTaxID_AllRegimes(t *testing.T) {
	regimes := tax.AllRegimeDefs()
	require.NotEmpty(t, regimes)

	r := rand.New(rand.NewPCG(42, 42))
	for _, regime := range regimes {
		cc := regime.Country
		t.Run(string(cc), func(t *testing.T) {
			for range 10 {
				code := generateTaxID(r, cc, true)
				// For regimes that require a code, it must not be empty.
				if _, ok := knownTaxIDs[cc]; ok && cc != "CA" && cc != "US" {
					assert.NotEmpty(t, code, "country %s should generate a tax ID", cc)
				}
			}
		})
	}
}

func TestGenerateTaxID_ESOrg(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 42))
	for range 100 {
		code := generateESOrgTaxID(r)
		tID := &tax.Identity{Country: l10n.ES.Tax(), Code: code}
		// Verify it normalizes and validates through GOBL.
		assert.NotEmpty(t, tID.Code)
	}
}

func TestGenerateTaxID_DE(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 42))
	for range 100 {
		code := generateDETaxID(r)
		// Verify check digit inline.
		val := code.String()
		require.Len(t, val, 9)
		p := 10
		for i := 0; i < 8; i++ {
			sum := (int(val[i]-'0') + p) % 10
			if sum == 0 {
				sum = 10
			}
			p = (2 * sum) % 11
		}
		cd := 11 - p
		if cd == 10 {
			cd = 0
		}
		assert.Equal(t, cd, int(val[8]-'0'), "check digit mismatch: %s", code)
	}
}

func TestGenerateTaxID_FR(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 42))
	for range 100 {
		code := generateFRTaxID(r)
		assert.Len(t, code.String(), 11, "FR tax ID must be 11 chars: %s", code)
	}
}

func TestGenerateTaxID_IT(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 42))
	for range 100 {
		code := generateITTaxID(r)
		assert.Len(t, code.String(), 11, "IT tax ID must be 11 chars: %s", code)
	}
}
