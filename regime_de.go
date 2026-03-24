package mock

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

func deConfig() *regimeConfig {
	return &regimeConfig{
		Country:  l10n.DE.Tax(),
		Currency: currency.EUR,
		SupplierNames: []string{
			"Müller & Schmidt GmbH",
			"Technologie Innovationen AG",
			"Berater Gruppe München GmbH",
			"Fischer Elektronik GmbH",
			"Weber Maschinenbau GmbH",
			"Braun Digital Solutions GmbH",
			"Schneider Logistik AG",
			"Hoffmann Consulting GmbH",
			"Koch Industrietechnik GmbH",
			"Bauer Software GmbH",
		},
		CustomerNames: []string{
			"Deutsche Handelskompanie GmbH",
			"Norddeutsche Energie AG",
			"Bayern Automotive GmbH",
			"Sachsen Elektronik GmbH",
			"Rhein-Ruhr Logistics GmbH",
			"Berliner Technologie GmbH",
			"Hamburg Port Services AG",
			"Schwaben Engineering GmbH",
			"Dresden Pharma GmbH",
			"Frankfurt Financial Services AG",
		},
		Cities: []cityData{
			{Name: "Berlin", Region: "Berlin", Code: "10115"},
			{Name: "München", Region: "Bayern", Code: "80331"},
			{Name: "Hamburg", Region: "Hamburg", Code: "20095"},
			{Name: "Frankfurt am Main", Region: "Hessen", Code: "60311"},
			{Name: "Köln", Region: "Nordrhein-Westfalen", Code: "50667"},
			{Name: "Stuttgart", Region: "Baden-Württemberg", Code: "70173"},
			{Name: "Düsseldorf", Region: "Nordrhein-Westfalen", Code: "40213"},
			{Name: "Dresden", Region: "Sachsen", Code: "01067"},
		},
		Streets: []string{
			"Friedrichstraße",
			"Kurfürstendamm",
			"Maximilianstraße",
			"Königstraße",
			"Bahnhofstraße",
			"Hauptstraße",
			"Berliner Allee",
			"Schillerstraße",
		},
		Products: []itemData{
			{Name: "Laptop Business Pro", Price: "1299.00"},
			{Name: "Bürostuhl Ergonomisch", Price: "449.00"},
			{Name: "Schreibtisch Höhenverstellbar", Price: "699.00"},
			{Name: "Drucker Multifunktion", Price: "389.00"},
			{Name: "Netzwerk Switch 24-Port", Price: "279.00"},
			{Name: "USV Anlage 1500VA", Price: "459.00"},
			{Name: "Dokumentenscanner", Price: "349.00"},
			{Name: "Server-Rack 42HE", Price: "1890.00"},
		},
		Services: []itemData{
			{Name: "IT-Beratung", Price: "120.00", Unit: "h"},
			{Name: "Softwareentwicklung", Price: "95.00", Unit: "h"},
			{Name: "Systemadministration", Price: "85.00", Unit: "h"},
			{Name: "Projektmanagement", Price: "110.00", Unit: "h"},
			{Name: "Datenbankoptimierung", Price: "130.00", Unit: "h"},
			{Name: "Schulung IT-Sicherheit", Price: "250.00", Unit: "service"},
			{Name: "Webhosting Premium", Price: "29.90"},
			{Name: "Cloud-Infrastruktur Monatlich", Price: "499.00"},
		},
		TaxRates: []taxRate{
			{Category: tax.CategoryVAT, Rate: tax.KeyStandard},
			{Category: tax.CategoryVAT, Rate: tax.RateReduced},
		},
		PaymentKey: pay.MeansKeyCreditTransfer.With(pay.MeansKeySEPA),
		IBANPrefix: "DE",
	}
}

func deXRechnungConfig(base *regimeConfig) *regimeConfig {
	cfg := *base
	cfg.SupplierPeople = true
	cfg.SupplierInboxes = true
	cfg.SupplierPhones = true
	cfg.CustomerInboxes = true
	cfg.RequiresOrdering = true
	return &cfg
}
