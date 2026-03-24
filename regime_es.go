package mock

import (
	"math/rand/v2"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

func esConfig() *regimeConfig {
	return &regimeConfig{
		Country:  l10n.ES.Tax(),
		Currency: currency.EUR,
		SupplierNames: []string{
			"Servicios Técnicos Avanzados S.L.",
			"Distribuciones García S.A.",
			"Consultores Ibéricos S.L.",
			"Ingeniería Solar Madrid S.L.",
			"Alimentación Mediterránea S.A.",
			"Construcciones Modernas S.L.",
			"Tecnología y Redes Barcelona S.L.",
			"Transportes Rápidos del Sur S.L.",
			"Farmacia Digital Valencia S.L.",
			"Diseño Creativo Andaluz S.L.",
		},
		CustomerNames: []string{
			"Empresa Nacional de Turismo S.A.",
			"Innovación Retail S.L.",
			"Gestión Documental Express S.L.",
			"Logística Peninsular S.A.",
			"Energía Verde Castilla S.L.",
			"Servicios Financieros del Norte S.A.",
			"Industrias Químicas Levante S.L.",
			"Soluciones Cloud Iberia S.L.",
			"Marketing Digital Sur S.L.",
			"Automoción Premium España S.A.",
		},
		Cities: []cityData{
			{Name: "Madrid", Region: "Madrid", Code: "28001"},
			{Name: "Barcelona", Region: "Barcelona", Code: "08001"},
			{Name: "Valencia", Region: "Valencia", Code: "46001"},
			{Name: "Sevilla", Region: "Sevilla", Code: "41001"},
			{Name: "Bilbao", Region: "Vizcaya", Code: "48001"},
			{Name: "Málaga", Region: "Málaga", Code: "29001"},
			{Name: "Zaragoza", Region: "Zaragoza", Code: "50001"},
			{Name: "Murcia", Region: "Murcia", Code: "30001"},
		},
		Streets: []string{
			"Calle Gran Vía",
			"Paseo de la Castellana",
			"Avenida Diagonal",
			"Calle Serrano",
			"Calle Mayor",
			"Paseo de Gracia",
			"Rambla de Catalunya",
			"Calle Alcalá",
		},
		Products: []itemData{
			{Name: "Ordenador portátil", Price: "899.00"},
			{Name: "Monitor LED 27 pulgadas", Price: "349.00"},
			{Name: "Teclado mecánico", Price: "89.00"},
			{Name: "Ratón inalámbrico", Price: "45.00"},
			{Name: "Disco SSD 1TB", Price: "119.00"},
			{Name: "Impresora láser", Price: "299.00"},
			{Name: "Cable USB-C", Price: "12.00"},
			{Name: "Webcam HD", Price: "65.00"},
		},
		Services: []itemData{
			{Name: "Desarrollo de software", Price: "90.00", Unit: "h"},
			{Name: "Consultoría empresarial", Price: "120.00", Unit: "h"},
			{Name: "Diseño gráfico", Price: "75.00", Unit: "h"},
			{Name: "Soporte técnico", Price: "60.00", Unit: "h"},
			{Name: "Auditoría contable", Price: "150.00", Unit: "h"},
			{Name: "Formación profesional", Price: "200.00", Unit: "service"},
			{Name: "Mantenimiento web mensual", Price: "500.00"},
			{Name: "Servicio de hosting anual", Price: "240.00"},
		},
		TaxRates: []taxRate{
			{Category: tax.CategoryVAT, Rate: tax.KeyStandard},
			{Category: tax.CategoryVAT, Rate: tax.RateReduced},
			{Category: tax.CategoryVAT, Rate: tax.RateSuperReduced},
		},
		PaymentKey: pay.MeansKeyCreditTransfer,
		IBANPrefix: "ES",
	}
}

func esFacturaeConfig(base *regimeConfig) *regimeConfig {
	cfg := *base
	// FacturaE requires customer with tax_id, which is the default behavior.
	// Scenarios auto-set es-facturae-doc-type and es-facturae-invoice-class.
	return &cfg
}

// esCorrectionExtensions returns the correction extension for FacturaE preceding references.
func esCorrectionExtensions(_ *rand.Rand) tax.Extensions {
	return tax.Extensions{
		"es-facturae-correction": cbc.Code("01"), // pricing correction
	}
}
