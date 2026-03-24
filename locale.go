package mock

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pay"
)

type cityData struct {
	Name   string
	Region string
	Code   cbc.Code
}

type itemData struct {
	Name  string
	Price string
	Unit  string
}

// localeData holds locale-specific content for realistic invoice generation.
type localeData struct {
	SupplierNames []string
	CustomerNames []string
	Cities        []cityData
	Streets       []string
	Products      []itemData
	Services      []itemData
	PaymentKey    cbc.Key
	IBANPrefix    string
}

var locales = map[l10n.TaxCountryCode]*localeData{
	"ES": {
		SupplierNames: []string{"Servicios Técnicos Avanzados S.L.", "Distribuciones García S.A.", "Consultores Ibéricos S.L.", "Ingeniería Solar Madrid S.L.", "Tecnología y Redes Barcelona S.L."},
		CustomerNames: []string{"Empresa Nacional de Turismo S.A.", "Innovación Retail S.L.", "Logística Peninsular S.A.", "Soluciones Cloud Iberia S.L.", "Industrias Químicas Levante S.L."},
		Cities:        []cityData{{Name: "Madrid", Region: "Madrid", Code: "28001"}, {Name: "Barcelona", Region: "Barcelona", Code: "08001"}, {Name: "Valencia", Region: "Valencia", Code: "46001"}, {Name: "Sevilla", Region: "Sevilla", Code: "41001"}, {Name: "Bilbao", Region: "Vizcaya", Code: "48001"}},
		Streets:       []string{"Calle Gran Vía", "Paseo de la Castellana", "Avenida Diagonal", "Calle Serrano", "Calle Mayor"},
		Products:      []itemData{{Name: "Ordenador portátil", Price: "899.00"}, {Name: "Monitor LED 27\"", Price: "349.00"}, {Name: "Impresora láser", Price: "299.00"}},
		Services:      []itemData{{Name: "Desarrollo de software", Price: "90.00", Unit: "h"}, {Name: "Consultoría empresarial", Price: "120.00", Unit: "h"}, {Name: "Soporte técnico", Price: "60.00", Unit: "h"}, {Name: "Mantenimiento web", Price: "500.00"}},
		PaymentKey:    pay.MeansKeyCreditTransfer,
		IBANPrefix:    "ES",
	},
	"DE": {
		SupplierNames: []string{"Müller & Schmidt GmbH", "Technologie Innovationen AG", "Fischer Elektronik GmbH", "Weber Maschinenbau GmbH", "Schneider Logistik AG"},
		CustomerNames: []string{"Deutsche Handelskompanie GmbH", "Bayern Automotive GmbH", "Berliner Technologie GmbH", "Hamburg Port Services AG", "Schwaben Engineering GmbH"},
		Cities:        []cityData{{Name: "Berlin", Region: "Berlin", Code: "10115"}, {Name: "München", Region: "Bayern", Code: "80331"}, {Name: "Hamburg", Region: "Hamburg", Code: "20095"}, {Name: "Frankfurt am Main", Region: "Hessen", Code: "60311"}, {Name: "Köln", Region: "Nordrhein-Westfalen", Code: "50667"}},
		Streets:       []string{"Friedrichstraße", "Kurfürstendamm", "Maximilianstraße", "Bahnhofstraße", "Hauptstraße"},
		Products:      []itemData{{Name: "Laptop Business Pro", Price: "1299.00"}, {Name: "Bürostuhl Ergonomisch", Price: "449.00"}, {Name: "Netzwerk Switch", Price: "279.00"}},
		Services:      []itemData{{Name: "IT-Beratung", Price: "120.00", Unit: "h"}, {Name: "Softwareentwicklung", Price: "95.00", Unit: "h"}, {Name: "Projektmanagement", Price: "110.00", Unit: "h"}, {Name: "Cloud-Infrastruktur", Price: "499.00"}},
		PaymentKey:    pay.MeansKeyCreditTransfer.With(pay.MeansKeySEPA),
		IBANPrefix:    "DE",
	},
	"MX": {
		SupplierNames: []string{"Soluciones Tecnológicas del Norte SA de CV", "Grupo Industrial Azteca SA de CV", "Servicios Profesionales Reforma SC", "Consultoría Empresarial CDMX SC", "Tecnología Avanzada Querétaro SA de CV"},
		CustomerNames: []string{"Distribuidora Nacional SA de CV", "Manufacturas de Exportación SA de CV", "Comercio Electrónico México SA de CV", "Energía Renovable Solar SA de CV", "Finanzas y Contabilidad SC"},
		Cities:        []cityData{{Name: "Ciudad de México", Region: "CMX", Code: "06600"}, {Name: "Guadalajara", Region: "JAL", Code: "44100"}, {Name: "Monterrey", Region: "NLE", Code: "64000"}, {Name: "Puebla", Region: "PUE", Code: "72000"}, {Name: "Querétaro", Region: "QUE", Code: "76000"}},
		Streets:       []string{"Avenida Insurgentes Sur", "Paseo de la Reforma", "Avenida Revolución", "Calle Hidalgo", "Boulevard Manuel Ávila Camacho"},
		Products:      []itemData{{Name: "Computadora de escritorio", Price: "15000.00"}, {Name: "Impresora multifuncional", Price: "6800.00"}, {Name: "Servidor rack 2U", Price: "45000.00"}},
		Services:      []itemData{{Name: "Servicio de consultoría", Price: "1500.00", Unit: "h"}, {Name: "Desarrollo de software", Price: "1200.00", Unit: "h"}, {Name: "Soporte técnico", Price: "800.00", Unit: "h"}, {Name: "Hospedaje web", Price: "500.00"}},
		PaymentKey:    pay.MeansKeyCreditTransfer,
	},
}

var defaultLocale = &localeData{
	SupplierNames: []string{"Acme Corporation", "Global Services Ltd.", "Tech Solutions Inc.", "Prime Industries", "Atlas Consulting"},
	CustomerNames: []string{"National Retail Corp.", "Metro Supply Chain Ltd.", "Pinnacle Holdings", "Horizon Technologies", "Sterling Manufacturing"},
	Cities:        []cityData{{Name: "Capital City", Code: "10001"}, {Name: "Commerce District", Code: "20001"}, {Name: "Business Center", Code: "30001"}},
	Streets:       []string{"Main Street", "Commerce Avenue", "Industrial Boulevard", "Market Street"},
	Products:      []itemData{{Name: "Professional laptop", Price: "1200.00"}, {Name: "Office desk", Price: "450.00"}, {Name: "LED monitor", Price: "400.00"}},
	Services:      []itemData{{Name: "IT consulting", Price: "120.00", Unit: "h"}, {Name: "Software development", Price: "100.00", Unit: "h"}, {Name: "Technical support", Price: "80.00", Unit: "h"}, {Name: "Monthly hosting", Price: "50.00"}},
	PaymentKey:    pay.MeansKeyCreditTransfer,
}

func getLocale(country l10n.TaxCountryCode) *localeData {
	if l, ok := locales[country]; ok {
		return l
	}
	return defaultLocale
}
