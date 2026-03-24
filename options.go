package mock

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
)

type options struct {
	regime     l10n.TaxCountryCode
	addon      cbc.Key
	invType    cbc.Key
	lines      int
	simplified bool
	seed       int64
	hasSeed    bool
	template   *bill.Invoice
}

// Option configures mock invoice generation.
type Option func(*options)

// WithRegime sets the tax regime country code.
func WithRegime(code l10n.TaxCountryCode) Option {
	return func(o *options) { o.regime = code }
}

// WithAddon specifies an addon to apply.
func WithAddon(key cbc.Key) Option {
	return func(o *options) { o.addon = key }
}

// WithType sets the invoice type (standard, credit-note, corrective, debit-note, proforma).
func WithType(t cbc.Key) Option {
	return func(o *options) { o.invType = t }
}

// WithLines sets the number of line items.
func WithLines(n int) Option {
	return func(o *options) { o.lines = n }
}

// WithCredit is shorthand for WithType(bill.InvoiceTypeCreditNote).
func WithCredit() Option {
	return func(o *options) { o.invType = bill.InvoiceTypeCreditNote }
}

// WithSimplified generates a simplified invoice (no customer).
func WithSimplified() Option {
	return func(o *options) { o.simplified = true }
}

// WithSeed sets a deterministic random seed for reproducible output.
func WithSeed(seed int64) Option {
	return func(o *options) { o.seed = seed; o.hasSeed = true }
}

// WithTemplate provides a partial invoice whose non-nil fields override
// the generated values. Use this to test specific scenarios while letting
// gobl.mock fill in the rest.
func WithTemplate(inv *bill.Invoice) Option {
	return func(o *options) { o.template = inv }
}

func defaultOptions() *options {
	return &options{
		regime: l10n.ES.Tax(),
		lines:  3,
	}
}
