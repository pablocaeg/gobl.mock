package mock

import (
	"math/rand/v2"
	"regexp"
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateESTaxID_National(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 42))
	for range 100 {
		code := generateESNationalTaxID(r)
		tID := &tax.Identity{Country: l10n.ES.Tax(), Code: code}
		es.TaxIdentityKey(tID)
		assert.Equal(t, es.TaxIdentityNational, es.TaxIdentityKey(tID), "code: %s", code)
	}
}

func TestGenerateESTaxID_Org(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 42))
	for range 100 {
		code := generateESOrgTaxID(r)
		tID := &tax.Identity{Country: l10n.ES.Tax(), Code: code}
		assert.Equal(t, es.TaxIdentityOrg, es.TaxIdentityKey(tID), "code: %s", code)
	}
}

func TestGenerateESTaxID_Foreigner(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 42))
	for range 100 {
		code := generateESForeignerTaxID(r)
		tID := &tax.Identity{Country: l10n.ES.Tax(), Code: code}
		assert.Equal(t, es.TaxIdentityForeigner, es.TaxIdentityKey(tID), "code: %s", code)
	}
}

func TestGenerateDETaxID(t *testing.T) {
	deFormat := regexp.MustCompile(`^[1-9]\d{8}$`)
	r := rand.New(rand.NewPCG(42, 42))
	for range 100 {
		code := generateDETaxID(r)
		require.True(t, deFormat.MatchString(code.String()), "format mismatch: %s", code)
		// Verify the check digit algorithm.
		val := code.String()
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

func TestGenerateMXTaxID_Company(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 42))
	for range 100 {
		code := generateMXCompanyTaxID(r)
		assert.Equal(t, mx.TaxIdentityTypeCompany, mx.DetermineTaxCodeType(code), "code: %s", code)
	}
}

func TestGenerateMXTaxID_Person(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 42))
	for range 100 {
		code := generateMXPersonTaxID(r)
		assert.Equal(t, mx.TaxIdentityTypePerson, mx.DetermineTaxCodeType(code), "code: %s", code)
	}
}
