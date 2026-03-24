package mock

import (
	"fmt"
	"math/rand/v2"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

// CFDI fiscal regime codes commonly used by businesses.
var mxFiscalRegimes = []string{
	"601", // General de Ley Personas Morales
	"603", // Personas Morales con Fines no Lucrativos
	"612", // Personas Físicas con Actividades Empresariales y Profesionales
	"620", // Sociedades Cooperativas de Producción
	"621", // Régimen de Incorporación Fiscal
	"626", // Régimen Simplificado de Confianza
}

// CFDI use codes.
var mxCFDIUses = []string{
	"G01", // Adquisición de mercancías
	"G02", // Devoluciones, descuentos o bonificaciones
	"G03", // Gastos en general
	"S01", // Sin efectos fiscales
}

// CFDI product/service codes (ClaveProdServ).
var mxProdServCodes = []string{
	"10101500", // Live animals
	"43211500", // Computers
	"43211700", // Computer accessories
	"80111600", // Temporary staffing
	"80141628", // Marketing management
	"84111500", // Accounting services
	"84111600", // Tax preparation
	"81112100", // Internet services
	"43232600", // Multimedia software
	"82101500", // Advertising
}

// CFDI postal codes (Mexican zip codes).
var mxPostalCodes = []string{
	"06600", // Ciudad de México - Juárez
	"01000", // Ciudad de México - Álvaro Obregón
	"44100", // Guadalajara
	"64000", // Monterrey
	"72000", // Puebla
	"76000", // Querétaro
	"20000", // Aguascalientes
	"86991", // Various
}

func mxConfig() *regimeConfig {
	return &regimeConfig{
		Country:  l10n.MX.Tax(),
		Currency: currency.MXN,
		SupplierNames: []string{
			"Soluciones Tecnológicas del Norte SA de CV",
			"Grupo Industrial Azteca SA de CV",
			"Comercializadora del Pacífico SA de CV",
			"Servicios Profesionales Reforma SC",
			"Industrias Metálicas Monterrey SA de CV",
			"Desarrollo Digital Cancún SA de CV",
			"Consultoría Empresarial CDMX SC",
			"Alimentos y Bebidas Jalisco SA de CV",
			"Logística Express Guadalajara SA de CV",
			"Tecnología Avanzada Querétaro SA de CV",
		},
		CustomerNames: []string{
			"Distribuidora Nacional SA de CV",
			"Manufacturas de Exportación SA de CV",
			"Servicios Integrales del Bajío SA de CV",
			"Comercio Electrónico México SA de CV",
			"Grupo Farmacéutico del Valle SA de CV",
			"Ingeniería y Construcción Maya SA de CV",
			"Textiles Modernos Puebla SA de CV",
			"Energía Renovable Solar SA de CV",
			"Automotriz Premium México SA de CV",
			"Finanzas y Contabilidad SC",
		},
		Cities: []cityData{
			{Name: "Ciudad de México", Region: "CMX", Code: "06600"},
			{Name: "Guadalajara", Region: "JAL", Code: "44100"},
			{Name: "Monterrey", Region: "NLE", Code: "64000"},
			{Name: "Puebla", Region: "PUE", Code: "72000"},
			{Name: "Querétaro", Region: "QUE", Code: "76000"},
			{Name: "Mérida", Region: "YUC", Code: "97000"},
			{Name: "Aguascalientes", Region: "AGU", Code: "20000"},
			{Name: "León", Region: "GUA", Code: "37000"},
		},
		Streets: []string{
			"Avenida Insurgentes Sur",
			"Paseo de la Reforma",
			"Avenida Revolución",
			"Calle Hidalgo",
			"Boulevard Manuel Ávila Camacho",
			"Avenida Universidad",
			"Calle Morelos",
			"Avenida Juárez",
		},
		Products: []itemData{
			{Name: "Computadora de escritorio", Price: "15000.00"},
			{Name: "Monitor LED 24 pulgadas", Price: "4500.00"},
			{Name: "Impresora multifuncional", Price: "6800.00"},
			{Name: "Licencia de software anual", Price: "3200.00"},
			{Name: "Servidor rack 2U", Price: "45000.00"},
			{Name: "Switch de red 48 puertos", Price: "12000.00"},
			{Name: "UPS 3000VA", Price: "8500.00"},
			{Name: "Disco duro externo 2TB", Price: "1800.00"},
		},
		Services: []itemData{
			{Name: "Servicio de consultoría", Price: "1500.00", Unit: "h"},
			{Name: "Desarrollo de software", Price: "1200.00", Unit: "h"},
			{Name: "Soporte técnico", Price: "800.00", Unit: "h"},
			{Name: "Diseño gráfico", Price: "900.00", Unit: "h"},
			{Name: "Auditoría fiscal", Price: "2500.00", Unit: "h"},
			{Name: "Capacitación profesional", Price: "3500.00", Unit: "service"},
			{Name: "Hospedaje web mensual", Price: "500.00"},
			{Name: "Mantenimiento preventivo", Price: "2000.00"},
		},
		TaxRates: []taxRate{
			{Category: tax.CategoryVAT, Rate: tax.KeyStandard},
		},
		PaymentKey:   pay.MeansKeyCreditTransfer,
		DefaultAddon: "mx-cfdi-v4",
	}
}

func mxCFDIConfig(base *regimeConfig) *regimeConfig {
	cfg := *base
	cfg.SupplierExt = func(r *rand.Rand) tax.Extensions {
		return tax.Extensions{
			"mx-cfdi-fiscal-regime": cbc.Code(pick(r, mxFiscalRegimes)),
		}
	}
	cfg.CustomerExt = func(r *rand.Rand) tax.Extensions {
		return tax.Extensions{
			"mx-cfdi-fiscal-regime": cbc.Code(pick(r, mxFiscalRegimes)),
			"mx-cfdi-use":           cbc.Code(pick(r, mxCFDIUses)),
		}
	}
	cfg.ItemExt = func(r *rand.Rand) tax.Extensions {
		return tax.Extensions{
			"mx-cfdi-prod-serv": cbc.Code(pick(r, mxProdServCodes)),
		}
	}
	cfg.InvoiceTaxExt = func(r *rand.Rand) tax.Extensions {
		return tax.Extensions{
			"mx-cfdi-issue-place": cbc.Code(pick(r, mxPostalCodes)),
		}
	}
	cfg.CustomerPostalCode = func(r *rand.Rand) cbc.Code {
		return cbc.Code(pick(r, mxPostalCodes))
	}
	return &cfg
}

// mxCorrectionStamp returns a fake SAT UUID stamp for credit note preceding references.
func mxCorrectionStamp() string {
	return fmt.Sprintf("a1b2c3d4-e5f6-7890-abcd-%012d", 1)
}
