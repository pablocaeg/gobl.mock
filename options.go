package mock

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
)

// Options configures mock invoice generation.
type Options struct {
	// Regime is the tax country code (e.g., "ES", "DE", "MX").
	Regime l10n.TaxCountryCode
	// Addon is an optional addon key (e.g., "es-facturae-v3", "de-xrechnung-v3").
	Addon cbc.Key
	// Lines is the number of line items to generate. Defaults to 3.
	Lines int
	// Credit generates a credit note instead of a standard invoice.
	Credit bool
	// Simplified generates a simplified invoice (no customer).
	Simplified bool
	// Seed provides a deterministic random seed for reproducible output.
	Seed int64
	// HasSeed indicates whether a seed was explicitly provided.
	HasSeed bool
}

// Option is a functional option for configuring invoice generation.
type Option func(*Options)

// WithRegime sets the tax regime country code.
func WithRegime(code l10n.TaxCountryCode) Option {
	return func(o *Options) {
		o.Regime = code
	}
}

// WithAddon specifies an addon to apply.
func WithAddon(key cbc.Key) Option {
	return func(o *Options) {
		o.Addon = key
	}
}

// WithLines sets the number of line items.
func WithLines(n int) Option {
	return func(o *Options) {
		o.Lines = n
	}
}

// WithCredit generates a credit note.
func WithCredit() Option {
	return func(o *Options) {
		o.Credit = true
	}
}

// WithSimplified generates a simplified invoice.
func WithSimplified() Option {
	return func(o *Options) {
		o.Simplified = true
	}
}

// WithSeed sets a deterministic random seed.
func WithSeed(seed int64) Option {
	return func(o *Options) {
		o.Seed = seed
		o.HasSeed = true
	}
}

func defaultOptions() *Options {
	return &Options{
		Regime: l10n.ES.Tax(),
		Lines:  3,
	}
}
