package mock

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"

	_ "github.com/invopop/gobl/addons"
	_ "github.com/invopop/gobl/regimes"
)

func TestProbeAddonExtensions(t *testing.T) {
	tests := []struct {
		addon   string
		wantKey string   // object type key in result
		wantExt []string // extension keys expected
	}{
		{"es-facturae-v3", "invoice.tax", []string{"es-facturae-doc-type", "es-facturae-invoice-class"}},
		{"es-tbai-v1", "invoice.tax", []string{"es-tbai-region"}},
		{"fr-choruspro-v1", "party", []string{"fr-choruspro-scheme"}},
		{"gr-mydata-v1", "pay", []string{"gr-mydata-payment-means"}},
	}

	for _, tc := range tests {
		t.Run(tc.addon, func(t *testing.T) {
			ad := tax.AddonForKey(cbc.Key(tc.addon))
			if ad == nil {
				t.Skipf("addon %s not found", tc.addon)
			}
			mapping := probeAddonExtensions(ad)
			keys, ok := mapping[tc.wantKey]
			if !ok {
				t.Fatalf("expected %s mapping for %s, got: %v", tc.wantKey, tc.addon, mapping)
			}
			for _, expected := range tc.wantExt {
				found := false
				for _, k := range keys {
					if string(k) == expected {
						found = true
						break
					}
				}
				assert.True(t, found, "expected extension %s in %s for %s", expected, tc.wantKey, tc.addon)
			}
		})
	}
}

func TestProbeAddonExtensions_AllAddons(t *testing.T) {
	// Verify probing doesn't panic on any addon.
	for _, ad := range tax.AllAddonDefs() {
		t.Run(string(ad.Key), func(t *testing.T) {
			mapping := probeAddonExtensions(ad)
			t.Logf("%s: discovered %d object types", ad.Key, len(mapping))
			for obj, keys := range mapping {
				t.Logf("  %s: %v", obj, keys)
			}
		})
	}
}
