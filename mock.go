// Package mock generates realistic, valid GOBL invoices for testing and development.
package mock

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"

	// Ensure all regimes and addons are registered.
	_ "github.com/invopop/gobl/addons"
	_ "github.com/invopop/gobl/regimes"
)

// Invoice generates a random valid GOBL invoice.
func Invoice(opts ...Option) (*bill.Invoice, error) {
	env, err := Envelope(opts...)
	if err != nil {
		return nil, err
	}
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("unexpected document type in envelope")
	}
	return inv, nil
}

// Envelope generates a random valid GOBL envelope containing an invoice.
func Envelope(opts ...Option) (*gobl.Envelope, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	r := newRand(o)

	// Resolve addon: use explicit, or default for regime.
	addon := o.Addon
	if addon == "" {
		if base, ok := regimeConfigs[o.Regime]; ok && base.DefaultAddon != "" {
			addon = base.DefaultAddon
		}
	}

	cfg := getRegimeConfig(o.Regime, addon)
	if cfg == nil {
		return nil, fmt.Errorf("unsupported regime: %s", o.Regime)
	}

	inv := &bill.Invoice{
		Series:    "MOCK",
		Code:      cbc.Code(fmt.Sprintf("%05d", r.IntN(99999)+1)),
		IssueDate: cal.Today(),
		Currency:  cfg.Currency,
	}

	inv.SetRegime(o.Regime)

	if addon != "" {
		inv.SetAddons(addon)
	}

	// Tags.
	if o.Simplified {
		inv.SetTags(tax.TagSimplified)
	}

	// Tax extensions (addon-specific).
	if cfg.InvoiceTaxExt != nil {
		inv.Tax = &bill.Tax{
			Ext: cfg.InvoiceTaxExt(r),
		}
	}

	// Supplier.
	inv.Supplier = generateSupplier(r, cfg)

	// Customer (unless simplified).
	if !o.Simplified {
		inv.Customer = generateCustomer(r, cfg)
	}

	// Lines.
	inv.Lines = generateLines(r, cfg, o.Lines)

	// Payment.
	inv.Payment = generatePayment(r, cfg)

	// Ordering (if required by addon).
	if cfg.RequiresOrdering {
		inv.Ordering = &bill.Ordering{
			Code: cbc.Code(fmt.Sprintf("PO-%05d", r.IntN(99999)+1)),
		}
	}

	// Credit note handling.
	if o.Credit {
		inv.Type = bill.InvoiceTypeCreditNote
		yesterday := cal.Today().Add(0, 0, -1)
		preceding := &org.DocumentRef{
			Identify:  uuid.Identify{UUID: uuid.V7()},
			Type:      bill.InvoiceTypeStandard,
			Series:    "MOCK",
			Code:      cbc.Code(fmt.Sprintf("%05d", r.IntN(99999)+1)),
			IssueDate: &yesterday,
		}
		// Addon-specific correction extensions.
		if addon == "es-facturae-v3" {
			preceding.Ext = esCorrectionExtensions(r)
		}
		if addon == "mx-cfdi-v4" {
			preceding.Stamps = []*head.Stamp{
				{
					Provider: "sat-uuid",
					Value:    mxCorrectionStamp(),
				},
			}
		}
		inv.Preceding = []*org.DocumentRef{preceding}
	}

	// Wrap in envelope (triggers Calculate + Validate).
	env, err := gobl.Envelop(inv)
	if err != nil {
		return nil, fmt.Errorf("building envelope: %w", err)
	}

	return env, nil
}

func generatePayment(r *rand.Rand, cfg *regimeConfig) *bill.PaymentDetails {
	instructions := &pay.Instructions{
		Key: cfg.paymentKey(),
	}

	// Add credit transfer details for bank transfer payments.
	if cfg.IBANPrefix != "" {
		instructions.CreditTransfer = []*pay.CreditTransfer{
			{
				IBAN: fakeIBAN(r, cfg.IBANPrefix),
				Name: "Bank Account",
			},
		}
	}

	if cfg.PaymentExt != nil {
		instructions.Ext = cfg.PaymentExt(r)
	}

	return &bill.PaymentDetails{
		Instructions: instructions,
		Terms: &pay.Terms{
			Key:   pay.TermKeyInstant,
			Notes: "Payment due upon receipt.",
		},
	}
}

// fakeIBAN generates a plausible-looking IBAN for the given country prefix.
func fakeIBAN(r *rand.Rand, prefix string) string {
	digits := make([]byte, 20)
	for i := range digits {
		digits[i] = byte('0' + r.IntN(10))
	}
	return prefix + string(digits[:20])
}

func newRand(o *Options) *rand.Rand {
	if o.HasSeed {
		return rand.New(rand.NewPCG(uint64(o.Seed), uint64(o.Seed)))
	}
	return rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()+1)))
}
