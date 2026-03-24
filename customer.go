package mock

import (
	"fmt"
	"math/rand/v2"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func generateCustomer(r *rand.Rand, cfg *regimeConfig) *org.Party {
	city := pick(r, cfg.Cities)

	customer := &org.Party{
		Name: pick(r, cfg.CustomerNames),
		TaxID: &tax.Identity{
			Country: cfg.Country,
			Code:    generateTaxID(r, cfg.Country, true),
		},
		Addresses: []*org.Address{
			{
				Street:   pick(r, cfg.Streets),
				Number:   fmt.Sprintf("%d", r.IntN(200)+1),
				Locality: city.Name,
				Region:   city.Region,
				Code:     city.Code,
				Country:  l10n.ISOCountryCode(cfg.Country),
			},
		},
	}

	// Override postal code if addon requires a specific one.
	if cfg.CustomerPostalCode != nil {
		customer.Addresses[0].Code = cfg.CustomerPostalCode(r)
	}

	if cfg.CustomerExt != nil {
		customer.Ext = cfg.CustomerExt(r)
	}

	if cfg.CustomerInboxes {
		customer.Inboxes = []*org.Inbox{
			{Email: "billing@customer-example.com"},
		}
	}

	return customer
}
