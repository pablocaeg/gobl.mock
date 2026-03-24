package mock

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
)

// Tax ID character tables.
const (
	esCheckLetters    = "TRWAGMYFPDXBNJZSQVHLCKE"
	esOrgTypes        = "ABCDEFGHJNPQRSUVW"
	esOrgCheckLetters = "JABCDEFGHI"

	ieCheckChars = "WABCDEFGHIJKLMNOPQRSTUV"
)

// knownTaxIDs provides fallback valid IDs from GOBL test files.
var knownTaxIDs = map[l10n.TaxCountryCode][]cbc.Code{
	"AE": {"123456789012345", "187654321098765"},
	"AR": {"30500010912", "20172543597"},
	"AT": {"U38516405", "U64727905"},
	"BE": {"0413172884", "0414445663"},
	"BR": {"05104582000170", "35549549506"},
	"CA": {"123456789"},
	"CH": {"E284156502", "E432825998"},
	"CO": {"412615332", "8110079918"},
	"DE": {"111111125", "282741168"},
	"DK": {"13585628", "88146328"},
	"EL": {"925667500", "064677095"},
	"ES": {"B98602642", "54387763P"},
	"FR": {"44732829320", "44391838042"},
	"GB": {"000472631", "350983637"},
	"IE": {"5343381W", "3628739L"},
	"IN": {"27AAPFU0939F1ZV", "29AAGCB7383J1Z4"},
	"IT": {"12345678903", "13029381004"},
	"MX": {"EKU9003173C9", "URE180429TM6"},
	"NL": {"000099995B57", "808661863B01"},
	"PL": {"9551893317", "1132191233"},
	"PT": {"999999990", "545259045"},
	"SE": {"556036079301", "556703748501"},
	"SG": {"M91234567X", "199912345A"},
}

// generateTaxID generates a valid tax identity code for the given country.
func generateTaxID(r *rand.Rand, country l10n.TaxCountryCode, isOrg bool) cbc.Code {
	switch country {
	case "AE":
		return generateAETaxID(r)
	case "AR":
		return generateARTaxID(r, isOrg)
	case "AT":
		return generateATTaxID(r)
	case "BE":
		return generateBETaxID(r)
	case "BR":
		return generateBRTaxID(r, isOrg)
	case "CA", "US":
		return "" // no code required
	case "CH":
		return generateCHTaxID(r)
	case "CO":
		return generateCOTaxID(r)
	case "DE":
		return generateDETaxID(r)
	case "DK":
		return generateDKTaxID(r)
	case "EL", "GR":
		return generateGRTaxID(r)
	case "ES":
		if isOrg {
			return generateESOrgTaxID(r)
		}
		return generateESNationalTaxID(r)
	case "FR":
		return generateFRTaxID(r)
	case "GB":
		return generateGBTaxID(r)
	case "IE":
		return generateIETaxID(r)
	case "IN":
		return generateINTaxID(r)
	case "IT":
		return generateITTaxID(r)
	case "MX":
		if isOrg {
			return generateMXCompanyTaxID(r)
		}
		return generateMXPersonTaxID(r)
	case "NL":
		return generateNLTaxID(r)
	case "PL":
		return generatePLTaxID(r)
	case "PT":
		return generatePTTaxID(r)
	case "SE":
		return generateSETaxID(r)
	case "SG":
		return generateSGTaxID(r)
	default:
		if ids, ok := knownTaxIDs[country]; ok && len(ids) > 0 {
			return pick(r, ids)
		}
		return ""
	}
}

// --- Helpers ---

func luhnCheckDigit(digits []int) int {
	sum := 0
	for i := len(digits) - 1; i >= 0; i-- {
		d := digits[i]
		if (len(digits)-1-i)%2 == 0 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return (10 - (sum % 10)) % 10
}

func digitsToString(digits []int) string {
	var sb strings.Builder
	for _, d := range digits {
		sb.WriteByte(byte('0' + d))
	}
	return sb.String()
}

func randomDigits(r *rand.Rand, n int) []int {
	d := make([]int, n)
	for i := range d {
		d[i] = r.IntN(10)
	}
	return d
}

// --- AE: 15 random digits ---

func generateAETaxID(r *rand.Rand) cbc.Code {
	d := randomDigits(r, 15)
	d[0] = r.IntN(9) + 1 // avoid leading zero
	return cbc.Code(digitsToString(d))
}

// --- AR: mod-11, prefixes [20,23,27,30,33] ---

func generateARTaxID(r *rand.Rand, isOrg bool) cbc.Code {
	prefixes := []string{"30", "33"}
	if !isOrg {
		prefixes = []string{"20", "23", "27"}
	}
	prefix := pick(r, prefixes)
	multipliers := []int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2}

	body := randomDigits(r, 8)
	digits := make([]int, 10)
	digits[0] = int(prefix[0] - '0')
	digits[1] = int(prefix[1] - '0')
	copy(digits[2:], body)

	sum := 0
	for i := 0; i < 10; i++ {
		sum += digits[i] * multipliers[i]
	}
	check := 11 - (sum % 11)
	switch check {
	case 11:
		check = 0
	case 10:
		check = 9
	}
	return cbc.Code(prefix + digitsToString(body) + strconv.Itoa(check))
}

// --- AT: Luhn+4 offset, format U + 8 digits ---

func generateATTaxID(r *rand.Rand) cbc.Code {
	mults := []int{1, 2, 1, 2, 1, 2, 1}
	digits := randomDigits(r, 7)

	sum := 0
	for i, m := range mults {
		x := digits[i] * m
		if x > 9 {
			sum += x/10 + x%10
		} else {
			sum += x
		}
	}
	check := (10 - ((sum + 4) % 10)) % 10
	return cbc.Code("U" + digitsToString(digits) + strconv.Itoa(check))
}

// --- BE: mod-97, 10-digit enterprise ---

func generateBETaxID(r *rand.Rand) cbc.Code {
	first := r.IntN(2) // 0 or 1
	body := randomDigits(r, 7)
	num, _ := strconv.Atoi(fmt.Sprintf("%d%s", first, digitsToString(body)))
	check := 97 - (num % 97)
	return cbc.Code(fmt.Sprintf("%d%s%02d", first, digitsToString(body), check))
}

// --- BR: CNPJ (14 digits) or CPF (11 digits) ---

func generateBRTaxID(r *rand.Rand, isOrg bool) cbc.Code {
	if isOrg {
		return generateBRCNPJ(r)
	}
	return generateBRCPF(r)
}

func generateBRCNPJ(r *rand.Rand) cbc.Code {
	digits := randomDigits(r, 12)
	// First check digit
	w1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += digits[i] * w1[i]
	}
	rem := sum % 11
	c1 := 0
	if rem >= 2 {
		c1 = 11 - rem
	}
	digits = append(digits, c1)

	// Second check digit
	w2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum = 0
	for i := 0; i < 13; i++ {
		sum += digits[i] * w2[i]
	}
	rem = sum % 11
	c2 := 0
	if rem >= 2 {
		c2 = 11 - rem
	}
	digits = append(digits, c2)
	return cbc.Code(digitsToString(digits))
}

func generateBRCPF(r *rand.Rand) cbc.Code {
	digits := randomDigits(r, 9)
	// First check digit
	w1 := []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * w1[i]
	}
	rem := sum % 11
	c1 := 0
	if rem >= 2 {
		c1 = 11 - rem
	}
	digits = append(digits, c1)

	// Second check digit
	w2 := []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
	sum = 0
	for i := 0; i < 10; i++ {
		sum += digits[i] * w2[i]
	}
	rem = sum % 11
	c2 := 0
	if rem >= 2 {
		c2 = 11 - rem
	}
	digits = append(digits, c2)
	return cbc.Code(digitsToString(digits))
}

// --- CH: mod-11, format E + 9 digits ---

func generateCHTaxID(r *rand.Rand) cbc.Code {
	mults := []int{5, 4, 3, 2, 7, 6, 5, 4}
	for {
		digits := randomDigits(r, 8)
		sum := 0
		for i, m := range mults {
			sum += digits[i] * m
		}
		check := 11 - (sum % 11)
		if check == 10 {
			continue // invalid combination, retry
		}
		if check == 11 {
			check = 0
		}
		return cbc.Code("E" + digitsToString(digits) + strconv.Itoa(check))
	}
}

// --- CO: mod-11 with prime multipliers ---

func generateCOTaxID(r *rand.Rand) cbc.Code {
	primes := []int{3, 7, 13, 17, 19, 23, 29, 37, 41, 43, 47, 53, 59, 67, 71}
	bodyLen := 8 + r.IntN(2) // 8 or 9 digits
	digits := randomDigits(r, bodyLen)
	digits[0] = r.IntN(9) + 1 // avoid leading zero

	sum := 0
	for i, v := range digits {
		sum += v * primes[bodyLen-i-1]
	}
	check := sum % 11
	if check >= 2 {
		check = 11 - check
	}
	return cbc.Code(digitsToString(digits) + strconv.Itoa(check))
}

// --- DE: ISO 7064, 9 digits ---

func generateDETaxID(r *rand.Rand) cbc.Code {
	digits := make([]int, 8)
	digits[0] = r.IntN(9) + 1
	for i := 1; i < 8; i++ {
		digits[i] = r.IntN(10)
	}

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
	return cbc.Code(digitsToString(digits) + strconv.Itoa(cd))
}

// --- DK: mod-11, multipliers [2,7,6,5,4,3,2,1] ---

func generateDKTaxID(r *rand.Rand) cbc.Code {
	mults := []int{2, 7, 6, 5, 4, 3, 2}
	for {
		digits := randomDigits(r, 7)
		digits[0] = r.IntN(9) + 1 // avoid leading zero

		sum := 0
		for i, m := range mults {
			sum += digits[i] * m
		}
		// 8th digit (weight 1) must make total divisible by 11
		check := (11 - (sum % 11)) % 11
		if check > 9 {
			continue
		}
		return cbc.Code(digitsToString(digits) + strconv.Itoa(check))
	}
}

// --- ES: Spain (National, Org, Foreigner) ---

func generateESNationalTaxID(r *rand.Rand) cbc.Code {
	n := r.IntN(99999998) + 1
	check := esCheckLetters[n%23]
	return cbc.Code(fmt.Sprintf("%08d%c", n, check))
}

func generateESOrgTaxID(r *rand.Rand) cbc.Code {
	typeLetter := esOrgTypes[r.IntN(len(esOrgTypes))]
	digits := randomDigits(r, 7)

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
	sb.WriteByte(check)
	return cbc.Code(sb.String())
}

// --- FR: SIREN Luhn + VAT mod-97 prefix ---

func generateFRTaxID(r *rand.Rand) cbc.Code {
	// Generate 8-digit SIREN body
	digits := randomDigits(r, 8)
	digits[0] = r.IntN(9) + 1

	// Compute Luhn check for 9th digit
	check := luhnCheckDigit(digits)
	digits = append(digits, check)

	// Compute VAT prefix
	siren, _ := strconv.Atoi(digitsToString(digits))
	vatCheck := (siren*100 + 12) % 97
	return cbc.Code(fmt.Sprintf("%02d%s", vatCheck, digitsToString(digits)))
}

// --- GB: weighted sum with dual check ---

func generateGBTaxID(r *rand.Rand) cbc.Code {
	mults := []int{8, 7, 6, 5, 4, 3, 2}
	for {
		// Generate 7 digits > 1000000
		digits := randomDigits(r, 7)
		digits[0] = r.IntN(9) + 1
		num, _ := strconv.Atoi(digitsToString(digits))
		if num <= 1000000 {
			continue
		}

		sum := 0
		for i, m := range mults {
			sum += digits[i] * m
		}
		// Subtract 97 repeatedly until <= 0
		checkDigit := sum
		for checkDigit > 0 {
			checkDigit -= 97
		}
		if checkDigit < 0 {
			checkDigit = -checkDigit
		}
		// Apply "new method" offset
		if checkDigit >= 55 {
			checkDigit -= 55
		} else {
			checkDigit += 42
		}
		if checkDigit > 99 {
			continue
		}
		return cbc.Code(fmt.Sprintf("%s%02d", digitsToString(digits), checkDigit))
	}
}

// --- GR: powers-of-2 weighted ---

func generateGRTaxID(r *rand.Rand) cbc.Code {
	digits := randomDigits(r, 8)
	digits[0] = r.IntN(9) + 1

	sum := 0
	for i := 0; i < 8; i++ {
		sum += digits[i] * (1 << uint(8-i))
	}
	check := (sum % 11) % 10
	return cbc.Code(digitsToString(digits) + strconv.Itoa(check))
}

// --- IE: new format (7 digits + check letter) ---

func generateIETaxID(r *rand.Rand) cbc.Code {
	digits := randomDigits(r, 7)
	digits[0] = r.IntN(9) + 1

	sum := 0
	weights := []int{8, 7, 6, 5, 4, 3, 2}
	for i, w := range weights {
		sum += digits[i] * w
	}
	check := ieCheckChars[sum%23]
	return cbc.Code(digitsToString(digits) + string(check))
}

// --- IN: GSTIN, mod-36 weighted ---

func generateINTaxID(r *rand.Rand) cbc.Code {
	const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const nonzeroAlnum = "123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Format: 2 digits + 5 upper + 4 digits + 1 upper + 1 nonzero-alnum + 'Z' + check
	state := fmt.Sprintf("%02d", r.IntN(37)+1) // 01-37
	pan := randomAlpha(r, 5)
	entity := fmt.Sprintf("%04d", r.IntN(10000))
	alpha := string(upper[r.IntN(26)])
	nza := string(nonzeroAlnum[r.IntN(len(nonzeroAlnum))])
	base := state + pan + entity + alpha + nza + "Z"

	// Compute mod-36 check
	sum := 0
	for i, ch := range base {
		val := charToVal36(ch)
		mult := 1
		if i%2 != 0 {
			mult = 2
		}
		product := val * mult
		sum += product/36 + product%36
	}
	check := (36 - (sum % 36)) % 36
	return cbc.Code(base + string(valToChar36(check)))
}

func charToVal36(ch rune) int {
	if ch >= '0' && ch <= '9' {
		return int(ch - '0')
	}
	return int(ch-'A') + 10
}

func valToChar36(v int) byte {
	if v < 10 {
		return byte('0' + v)
	}
	return byte('A' + v - 10)
}

// --- IT: Luhn, 11 digits ---

func generateITTaxID(r *rand.Rand) cbc.Code {
	digits := randomDigits(r, 10)
	digits[0] = r.IntN(9) + 1
	check := luhnCheckDigit(digits)
	return cbc.Code(digitsToString(digits) + strconv.Itoa(check))
}

// --- MX: format only ---

func generateMXCompanyTaxID(r *rand.Rand) cbc.Code {
	prefix := randomAlpha(r, 3)
	date := randomMXDate(r)
	suffix := randomAlphaNum(r, 3)
	return cbc.Code(prefix + date + suffix)
}

func generateMXPersonTaxID(r *rand.Rand) cbc.Code {
	prefix := randomAlpha(r, 4)
	date := randomMXDate(r)
	suffix := randomAlphaNum(r, 3)
	return cbc.Code(prefix + date + suffix)
}

func randomMXDate(r *rand.Rand) string {
	year := r.IntN(40) + 60
	month := r.IntN(12) + 1
	day := r.IntN(28) + 1
	return fmt.Sprintf("%02d%02d%02d", year, month, day)
}

// --- NL: mod-11, 9 digits + B + 2 digits ---

func generateNLTaxID(r *rand.Rand) cbc.Code {
	for {
		digits := randomDigits(r, 8)
		// Compute mod-11 for 9th digit (old method)
		// Multipliers 9,8,7,6,5,4,3,2 applied right-to-left on first 8 digits
		sum := 0
		for i := 0; i < 8; i++ {
			sum += digits[i] * (9 - i)
		}
		check := sum % 11
		if check > 9 {
			continue // invalid, retry
		}
		suffix := fmt.Sprintf("%02d", r.IntN(100))
		return cbc.Code(digitsToString(digits) + strconv.Itoa(check) + "B" + suffix)
	}
}

// --- PL: weighted mod-11, weights [6,5,7,2,3,4,5,6,7] ---

func generatePLTaxID(r *rand.Rand) cbc.Code {
	weights := []int{6, 5, 7, 2, 3, 4, 5, 6, 7}
	for {
		digits := randomDigits(r, 9)
		digits[0] = r.IntN(9) + 1
		// Ensure positions 1-2 aren't both zero
		if digits[1] == 0 && digits[2] == 0 {
			digits[1] = r.IntN(9) + 1
		}

		sum := 0
		for i, w := range weights {
			sum += digits[i] * w
		}
		check := sum % 11
		if check > 9 {
			continue
		}
		return cbc.Code(digitsToString(digits) + strconv.Itoa(check))
	}
}

// --- PT: weighted mod-11, valid prefixes ---

func generatePTTaxID(r *rand.Rand) cbc.Code {
	prefixes := []string{"1", "2", "3", "5", "6", "8"}
	prefix := pick(r, prefixes)

	for {
		// Fill remaining digits to make 8-digit body
		remaining := randomDigits(r, 8-len(prefix))
		body := prefix + digitsToString(remaining)

		sum := 0
		for i := 0; i < 8; i++ {
			d, _ := strconv.Atoi(string(body[i]))
			sum += d * (9 - i)
		}
		rem := sum % 11
		check := 0
		if rem >= 2 {
			check = 11 - rem
		}
		if check > 9 {
			continue
		}
		return cbc.Code(body + strconv.Itoa(check))
	}
}

// --- SE: Luhn + "01" suffix ---

func generateSETaxID(r *rand.Rand) cbc.Code {
	digits := randomDigits(r, 9)
	digits[0] = r.IntN(9) + 1
	check := luhnCheckDigit(digits)
	return cbc.Code(digitsToString(digits) + strconv.Itoa(check) + "01")
}

// --- SG: format only ---

func generateSGTaxID(r *rand.Rand) cbc.Code {
	// UEN ROC format: 4-digit year + 5 digits + letter
	const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	year := 1990 + r.IntN(36) // 1990-2025
	num := r.IntN(100000)
	letter := upper[r.IntN(26)]
	return cbc.Code(fmt.Sprintf("%d%05d%c", year, num, letter))
}

// --- Shared string helpers ---

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
