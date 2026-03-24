package mock

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
)

// Options configures mock invoice generation.
type Options struct {
	Regime     l10n.TaxCountryCode
	Addon      cbc.Key
	Type       cbc.Key // invoice type: standard, credit-note, corrective, debit-note, proforma
	Lines      int
	Simplified bool
	Seed       int64
	HasSeed    bool
	Template   *bill.Invoice // partial invoice to merge over generated fields
}

// Option is a functional option for configuring invoice generation.
type Option func(*Options)

// WithRegime sets the tax regime country code.
func WithRegime(code l10n.TaxCountryCode) Option {
	return func(o *Options) { o.Regime = code }
}

// WithAddon specifies an addon to apply.
func WithAddon(key cbc.Key) Option {
	return func(o *Options) { o.Addon = key }
}

// WithType sets the invoice type (standard, credit-note, corrective, debit-note, proforma).
func WithType(t cbc.Key) Option {
	return func(o *Options) { o.Type = t }
}

// WithLines sets the number of line items.
func WithLines(n int) Option {
	return func(o *Options) { o.Lines = n }
}

// WithCredit is shorthand for WithType(bill.InvoiceTypeCreditNote).
func WithCredit() Option {
	return func(o *Options) { o.Type = bill.InvoiceTypeCreditNote }
}

// WithSimplified generates a simplified invoice (no customer).
func WithSimplified() Option {
	return func(o *Options) { o.Simplified = true }
}

// WithSeed sets a deterministic random seed for reproducible output.
func WithSeed(seed int64) Option {
	return func(o *Options) { o.Seed = seed; o.HasSeed = true }
}

// WithTemplate provides a partial invoice whose non-nil fields override
// the generated values. Use this to test specific scenarios while letting
// gobl.mock fill in the rest.
func WithTemplate(inv *bill.Invoice) Option {
	return func(o *Options) { o.Template = inv }
}

func defaultOptions() *Options {
	return &Options{
		Regime: l10n.ES.Tax(),
		Lines:  3,
	}
}
