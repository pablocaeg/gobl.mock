package mock

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
)

// Tax ID check character tables for Spain.
const (
	esCheckLetters    = "TRWAGMYFPDXBNJZSQVHLCKE"
	esOrgTypes        = "ABCDEFGHJNPQRSUVW"
	esOrgCheckLetters = "JABCDEFGHI"
	esForeignerTypes  = "XYZ"
)

// generateTaxID generates a valid tax identity code for the given country.
func generateTaxID(r *rand.Rand, country l10n.TaxCountryCode, isOrg bool) cbc.Code {
	switch country {
	case l10n.ES.Tax():
		if isOrg {
			return generateESOrgTaxID(r)
		}
		return generateESNationalTaxID(r)
	case l10n.DE.Tax():
		return generateDETaxID(r)
	case l10n.MX.Tax():
		if isOrg {
			return generateMXCompanyTaxID(r)
		}
		return generateMXPersonTaxID(r)
	default:
		return ""
	}
}

// generateESNationalTaxID generates a valid Spanish national NIF (8 digits + check letter).
func generateESNationalTaxID(r *rand.Rand) cbc.Code {
	n := r.IntN(99999998) + 1 // 1 to 99999998
	check := esCheckLetters[n%23]
	return cbc.Code(fmt.Sprintf("%08d%c", n, check))
}

// generateESOrgTaxID generates a valid Spanish organization CIF (letter + 7 digits + check).
func generateESOrgTaxID(r *rand.Rand) cbc.Code {
	typeLetter := esOrgTypes[r.IntN(len(esOrgTypes))]

	digits := make([]int, 7)
	for i := range digits {
		digits[i] = r.IntN(10)
	}

	// Luhn-like check digit algorithm matching GOBL's verification.
	sumEven := 0
	sumOdd := 0
	for k, v := range digits {
		if k&1 == 0 {
			d := v * 2
			if d > 9 {
				d -= 9
			}
			sumOdd += d
		} else {
			sumEven += v
		}
	}
	cdc := (10 - (sumEven+sumOdd)%10) % 10
	check := esOrgCheckLetters[cdc]

	var sb strings.Builder
	sb.WriteByte(typeLetter)
	for _, d := range digits {
		sb.WriteByte(byte('0' + d))
	}
	sb.WriteByte(byte(check))
	return cbc.Code(sb.String())
}

// generateESForeignerTaxID generates a valid Spanish foreigner NIE (X/Y/Z + 7 digits + check).
func generateESForeignerTaxID(r *rand.Rand) cbc.Code {
	typeIdx := r.IntN(len(esForeignerTypes))
	typeLetter := esForeignerTypes[typeIdx]

	number := r.IntN(10000000) // 0 to 9999999
	composite := typeIdx*10000000 + number
	check := esCheckLetters[composite%23]

	return cbc.Code(fmt.Sprintf("%c%07d%c", typeLetter, number, check))
}

// generateDETaxID generates a valid German USt-IdNr (9 digits with check digit).
func generateDETaxID(r *rand.Rand) cbc.Code {
	digits := make([]int, 8)
	digits[0] = r.IntN(9) + 1 // first digit 1-9
	for i := 1; i < 8; i++ {
		digits[i] = r.IntN(10)
	}

	// Iterative modular check digit algorithm.
	p := 10
	for i := 0; i < 8; i++ {
		sum := (digits[i] + p) % 10
		if sum == 0 {
			sum = 10
		}
		p = (2 * sum) % 11
	}
	cd := 11 - p
	if cd == 10 {
		cd = 0
	}

	var sb strings.Builder
	for _, d := range digits {
		sb.WriteByte(byte('0' + d))
	}
	sb.WriteByte(byte('0' + cd))
	return cbc.Code(sb.String())
}

// generateMXCompanyTaxID generates a valid Mexican RFC for a company (3 alpha + 6 date + 3 alnum).
func generateMXCompanyTaxID(r *rand.Rand) cbc.Code {
	prefix := randomAlpha(r, 3)
	date := randomMXDate(r)
	suffix := randomAlphaNum(r, 3)
	return cbc.Code(prefix + date + suffix)
}

// generateMXPersonTaxID generates a valid Mexican RFC for a person (4 alpha + 6 date + 3 alnum).
func generateMXPersonTaxID(r *rand.Rand) cbc.Code {
	prefix := randomAlpha(r, 4)
	date := randomMXDate(r)
	suffix := randomAlphaNum(r, 3)
	return cbc.Code(prefix + date + suffix)
}

func randomMXDate(r *rand.Rand) string {
	year := r.IntN(40) + 60 // 60-99
	month := r.IntN(12) + 1
	day := r.IntN(28) + 1
	return fmt.Sprintf("%02d%02d%02d", year, month, day)
}

func randomAlpha(r *rand.Rand, n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.IntN(len(letters))]
	}
	return string(b)
}

func randomAlphaNum(r *rand.Rand, n int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[r.IntN(len(chars))]
	}
	return string(b)
}
