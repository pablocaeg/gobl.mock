package mock

import (
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// probeAddonExtensions discovers which extensions an addon's validator
// requires on each object type by calling the validator with minimal
// instances and parsing the validation errors. This lets gobl.mock
// support new addons without hardcoded configs.
func probeAddonExtensions(ad *tax.AddonDef) map[string][]cbc.Key {
	if ad == nil || ad.Validator == nil {
		return nil
	}

	result := make(map[string][]cbc.Key)

	// Probe the invoice validator for tax.ext requirements.
	if keys := probeExtKeys(ad.Validator, minimalInvoice()); len(keys) > 0 {
		result["invoice.tax"] = keys
	}

	// Probe standalone type validators.
	if keys := probeExtKeys(ad.Validator, minimalParty()); len(keys) > 0 {
		result["party"] = keys
	}
	if keys := probeExtKeys(ad.Validator, minimalCombo()); len(keys) > 0 {
		result["combo"] = keys
	}
	if keys := probeExtKeys(ad.Validator, minimalItem()); len(keys) > 0 {
		result["item"] = keys
	}
	if keys := probeExtKeys(ad.Validator, minimalPayInstructions()); len(keys) > 0 {
		result["pay"] = keys
	}

	return result
}

// probeExtKeys calls a validator with a minimal object and extracts
// extension keys from "required" validation errors.
func probeExtKeys(validator func(any) error, obj any) []cbc.Key {
	err := validator(obj)
	if err == nil {
		return nil
	}
	return extractRequiredExtKeys(err)
}

// extractRequiredExtKeys recursively walks validation errors to find
// extension keys that are "required".
func extractRequiredExtKeys(err error) []cbc.Key {
	var keys []cbc.Key

	ve, ok := err.(validation.Errors)
	if !ok {
		return nil
	}

	for field, fieldErr := range ve {
		if field == "ext" {
			// The ext field contains key->error entries.
			if extErrs, ok := fieldErr.(validation.Errors); ok {
				for key, keyErr := range extErrs {
					if keyErr != nil && strings.Contains(keyErr.Error(), "required") {
						keys = append(keys, cbc.Key(key))
					}
				}
			}
			continue
		}
		// Recurse into nested fields (e.g., "tax" -> "ext").
		keys = append(keys, extractRequiredExtKeys(fieldErr)...)
	}

	return keys
}

func minimalInvoice() *bill.Invoice {
	return &bill.Invoice{
		Tax: &bill.Tax{Ext: tax.Extensions{}},
	}
}

func minimalParty() *org.Party {
	return &org.Party{
		Ext: tax.Extensions{},
	}
}

func minimalCombo() *tax.Combo {
	return &tax.Combo{
		Category: tax.CategoryVAT,
		Ext:      tax.Extensions{},
	}
}

func minimalItem() *org.Item {
	price := num.MakeAmount(100, 2)
	return &org.Item{
		Name:  "test",
		Price: &price,
		Ext:   tax.Extensions{},
	}
}

func minimalPayInstructions() *pay.Instructions {
	return &pay.Instructions{
		Key: pay.MeansKeyCreditTransfer,
		Ext: tax.Extensions{},
	}
}
