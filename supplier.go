package mock

import (
	"fmt"
	"math/rand/v2"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func generateSupplier(r *rand.Rand, cfg *regimeConfig) *org.Party {
	city := pick(r, cfg.Cities)

	supplier := &org.Party{
		Name: pick(r, cfg.SupplierNames),
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
		Emails: []*org.Email{
			{Address: "billing@example.com"},
		},
	}

	if cfg.SupplierExt != nil {
		supplier.Ext = cfg.SupplierExt(r)
	}

	if cfg.SupplierPeople {
		supplier.People = []*org.Person{
			{
				Name: &org.Name{
					Given:   pick(r, givenNames(cfg.Country)),
					Surname: pick(r, surnames(cfg.Country)),
				},
				Emails: []*org.Email{
					{Address: "contact@example.com"},
				},
			},
		}
	}

	if cfg.SupplierInboxes {
		supplier.Inboxes = []*org.Inbox{
			{Email: "invoices@example.com"},
		}
	}

	if cfg.SupplierPhones {
		supplier.Telephones = []*org.Telephone{
			{Number: phonePrefix(cfg.Country) + "100200300"},
		}
	}

	return supplier
}

func givenNames(country l10n.TaxCountryCode) []string {
	switch country {
	case l10n.ES.Tax():
		return []string{"María", "Carlos", "Ana", "José", "Laura", "Miguel", "Carmen", "Pablo"}
	case l10n.DE.Tax():
		return []string{"Thomas", "Anna", "Michael", "Sarah", "Stefan", "Julia", "Andreas", "Katrin"}
	case l10n.MX.Tax():
		return []string{"Alejandro", "María", "Carlos", "Ana", "Roberto", "Claudia", "Fernando", "Lucía"}
	default:
		return []string{"John", "Jane"}
	}
}

func surnames(country l10n.TaxCountryCode) []string {
	switch country {
	case l10n.ES.Tax():
		return []string{"García", "Martínez", "López", "Hernández", "González", "Rodríguez", "Fernández", "Sánchez"}
	case l10n.DE.Tax():
		return []string{"Müller", "Schmidt", "Schneider", "Fischer", "Weber", "Meyer", "Wagner", "Bauer"}
	case l10n.MX.Tax():
		return []string{"Hernández", "García", "López", "Martínez", "González", "Rodríguez", "Pérez", "Sánchez"}
	default:
		return []string{"Smith", "Doe"}
	}
}

func phonePrefix(country l10n.TaxCountryCode) string {
	switch country {
	case l10n.ES.Tax():
		return "+34"
	case l10n.DE.Tax():
		return "+49"
	case l10n.MX.Tax():
		return "+52"
	default:
		return "+1"
	}
}
