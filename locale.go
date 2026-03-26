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

func getLocale(country l10n.TaxCountryCode) *localeData {
	if l, ok := locales[country]; ok {
		return l
	}
	return defaultLocale
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

var sepa = pay.MeansKeyCreditTransfer.With(pay.MeansKeySEPA)

var locales = map[l10n.TaxCountryCode]*localeData{
	"AE": {
		SupplierNames: []string{"Al Futtaim Trading L.L.C.", "Gulf Horizon Logistics L.L.C.", "Emirates Precision Engineering L.L.C.", "Desert Palm Technologies FZE", "Oasis Digital Solutions FZ-LLC"},
		CustomerNames: []string{"Al Maha Consulting FZ-LLC", "Nakheel Food Industries L.L.C.", "Al Reem Building Materials L.L.C.", "Arabian Gulf Services L.L.C.", "Dubai Commerce Group FZE"},
		Cities:        []cityData{{Name: "Dubai", Code: "00000"}, {Name: "Abu Dhabi", Code: "00000"}, {Name: "Sharjah", Code: "00000"}},
		Streets:       []string{"Sheikh Zayed Road", "Al Maktoum Street", "Khalifa Bin Zayed Street", "Corniche Road", "Al Wasl Road"},
		Products:      []itemData{{Name: "Industrial equipment", Price: "5000.00"}, {Name: "Office furniture set", Price: "3200.00"}, {Name: "Server hardware", Price: "8500.00"}},
		Services:      []itemData{{Name: "Engineering consultation", Price: "200.00", Unit: "h"}, {Name: "IT infrastructure setup", Price: "180.00", Unit: "h"}, {Name: "Project management", Price: "150.00", Unit: "h"}},
		PaymentKey:    pay.MeansKeyCreditTransfer,
		IBANPrefix:    "AE",
	},
	"AR": {
		SupplierNames: []string{"Pampa Logística S.A.", "Soluciones del Sur S.R.L.", "Industrias Rioplatenses S.A.", "Tecnología Austral S.A.S.", "Metalúrgica Patagonia S.A."},
		CustomerNames: []string{"Consultora Federal S.R.L.", "Agropecuaria Los Álamos S.A.", "Cervecería Andina S.R.L.", "Distribuidora Porteña S.A.", "Electrónica Nacional S.R.L."},
		Cities:        []cityData{{Name: "Buenos Aires", Region: "Ciudad Autónoma de Buenos Aires", Code: "C1001"}, {Name: "Córdoba", Region: "Córdoba", Code: "X5000"}, {Name: "Rosario", Region: "Santa Fe", Code: "S2000"}, {Name: "Mendoza", Region: "Mendoza", Code: "M5500"}},
		Streets:       []string{"Avenida Rivadavia", "Calle San Martín", "Avenida Corrientes", "Calle Belgrano", "Avenida Santa Fe"},
		Products:      []itemData{{Name: "Equipo industrial", Price: "45000.00"}, {Name: "Material de oficina", Price: "12000.00"}, {Name: "Computadora portátil", Price: "85000.00"}},
		Services:      []itemData{{Name: "Consultoría empresarial", Price: "8000.00", Unit: "h"}, {Name: "Desarrollo de software", Price: "6500.00", Unit: "h"}, {Name: "Auditoría contable", Price: "9000.00", Unit: "h"}},
		PaymentKey:    pay.MeansKeyCreditTransfer,
	},
	"AT": {
		SupplierNames: []string{"Alpentechnik GmbH", "Wiener Präzisionswerkzeuge AG", "Donau Spedition GmbH", "Tiroler Holzverarbeitung OG", "Linzer Softwarelösungen GmbH"},
		CustomerNames: []string{"Steirische Lebensmittel GmbH", "Salzburger Beratung GmbH", "Kärntner Maschinenbau AG", "Oberösterreichische Elektronik GmbH", "Vorarlberger Textil AG"},
		Cities:        []cityData{{Name: "Wien", Region: "Wien", Code: "1010"}, {Name: "Graz", Region: "Steiermark", Code: "8010"}, {Name: "Linz", Region: "Oberösterreich", Code: "4020"}, {Name: "Salzburg", Region: "Salzburg", Code: "5020"}, {Name: "Innsbruck", Region: "Tirol", Code: "6020"}},
		Streets:       []string{"Kärntner Straße", "Mariahilfer Straße", "Hauptplatz", "Landstraßer Hauptstraße", "Graben"},
		Products:      []itemData{{Name: "Präzisionswerkzeug", Price: "890.00"}, {Name: "Büroausstattung", Price: "650.00"}, {Name: "Industriemotor", Price: "2400.00"}},
		Services:      []itemData{{Name: "Unternehmensberatung", Price: "130.00", Unit: "h"}, {Name: "Softwareentwicklung", Price: "110.00", Unit: "h"}, {Name: "Technische Planung", Price: "95.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "AT",
	},
	"BE": {
		SupplierNames: []string{"Vlaamse Textielfabriek BV", "Brabantse Consultants BV", "Antwerpse Havendiensten NV", "Gentse Softwareoplossingen BV", "Waalse Voeding SA"},
		CustomerNames: []string{"Bruxelles Logistique SA", "Leuvens Ingenieursbureau BV", "Ardennen Houtindustrie SA", "Limburgse Technologie BV", "Namur Services SA"},
		Cities:        []cityData{{Name: "Bruxelles", Region: "Bruxelles-Capitale", Code: "1000"}, {Name: "Antwerpen", Region: "Vlaanderen", Code: "2000"}, {Name: "Gent", Region: "Vlaanderen", Code: "9000"}, {Name: "Liège", Region: "Wallonie", Code: "4000"}},
		Streets:       []string{"Rue de la Loi", "Koningstraat", "Grote Markt", "Avenue Louise", "Meir"},
		Products:      []itemData{{Name: "Industriële machines", Price: "3500.00"}, {Name: "Kantoormeubelen", Price: "980.00"}, {Name: "Elektronische apparatuur", Price: "1200.00"}},
		Services:      []itemData{{Name: "Advies en consultancy", Price: "125.00", Unit: "h"}, {Name: "IT-ondersteuning", Price: "95.00", Unit: "h"}, {Name: "Logistiek beheer", Price: "85.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "BE",
	},
	"BR": {
		SupplierNames: []string{"Industrias Paulista Ltda.", "Soluções Cariocas Ltda.", "Tecnologia Tropical S.A.", "Logística Meridional Ltda.", "Construtora Atlântica S.A."},
		CustomerNames: []string{"Alimentos do Cerrado Ltda.", "Consultoria Brasileira Ltda.", "Máquinas Nordeste Ltda.", "Comércio Digital S.A.", "Engenharia Costeira Ltda."},
		Cities:        []cityData{{Name: "São Paulo", Region: "SP", Code: "01310-100"}, {Name: "Rio de Janeiro", Region: "RJ", Code: "20040-020"}, {Name: "Belo Horizonte", Region: "MG", Code: "30130-000"}, {Name: "Curitiba", Region: "PR", Code: "80010-010"}},
		Streets:       []string{"Rua São José", "Avenida Paulista", "Rua Tiradentes", "Avenida Brasil", "Rua Sete de Setembro"},
		Products:      []itemData{{Name: "Computador de mesa", Price: "4500.00"}, {Name: "Impressora multifuncional", Price: "2800.00"}, {Name: "Equipamento de rede", Price: "3200.00"}},
		Services:      []itemData{{Name: "Consultoria de TI", Price: "350.00", Unit: "h"}, {Name: "Desenvolvimento de software", Price: "280.00", Unit: "h"}, {Name: "Suporte técnico", Price: "180.00", Unit: "h"}},
		PaymentKey:    pay.MeansKeyCreditTransfer,
	},
	"CA": {
		SupplierNames: []string{"Northern Shield Resources Ltd.", "Maple Leaf Machining Inc.", "Pacific Coast Software Ltd.", "Cascade Manufacturing Corp.", "Rocky Mountain Foods Ltd."},
		CustomerNames: []string{"Prairie Wind Logistics Inc.", "Great Lakes Engineering Inc.", "Boréal Consultation Inc.", "Atlantic Trade Services Ltd.", "Harbour City Technologies Inc."},
		Cities:        []cityData{{Name: "Toronto", Region: "ON", Code: "M5H 2N2"}, {Name: "Vancouver", Region: "BC", Code: "V6B 1A1"}, {Name: "Montréal", Region: "QC", Code: "H3B 1S6"}, {Name: "Calgary", Region: "AB", Code: "T2P 3M3"}},
		Streets:       []string{"Bay Street", "Rue Sainte-Catherine", "Granville Street", "Albert Street", "Stephen Avenue"},
		Products:      []itemData{{Name: "Industrial sensor kit", Price: "890.00"}, {Name: "Office workstation", Price: "1400.00"}, {Name: "Network equipment", Price: "2200.00"}},
		Services:      []itemData{{Name: "Management consulting", Price: "175.00", Unit: "h"}, {Name: "Software development", Price: "150.00", Unit: "h"}, {Name: "Environmental assessment", Price: "200.00", Unit: "h"}},
		PaymentKey:    pay.MeansKeyCreditTransfer,
	},
	"CH": {
		SupplierNames: []string{"Zürcher Feinmechanik AG", "Berner Softwarelösungen GmbH", "Genève Conseil SA", "Basler Pharmaindustrie AG", "Lausanne Logistique Sàrl"},
		CustomerNames: []string{"Luzerner Lebensmittel GmbH", "Tessiner Handelsbetrieb SA", "Winterthurer Maschinenbau AG", "St. Galler Elektronik AG", "Thurgauer Consulting GmbH"},
		Cities:        []cityData{{Name: "Zürich", Region: "ZH", Code: "8001"}, {Name: "Genève", Region: "GE", Code: "1201"}, {Name: "Basel", Region: "BS", Code: "4001"}, {Name: "Bern", Region: "BE", Code: "3001"}, {Name: "Lausanne", Region: "VD", Code: "1003"}},
		Streets:       []string{"Bahnhofstrasse", "Rue du Rhône", "Freie Strasse", "Marktgasse", "Avenue de la Gare"},
		Products:      []itemData{{Name: "Präzisionsinstrument", Price: "2400.00"}, {Name: "Laborausrüstung", Price: "5600.00"}, {Name: "Medizintechnik", Price: "8900.00"}},
		Services:      []itemData{{Name: "Finanzberatung", Price: "250.00", Unit: "h"}, {Name: "Pharmaforschung", Price: "300.00", Unit: "h"}, {Name: "IT-Sicherheitsaudit", Price: "220.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "CH",
	},
	"CO": {
		SupplierNames: []string{"Soluciones Andinas S.A.S.", "Tecnología Cafetera S.A.S.", "Industrias del Pacífico Ltda.", "Logística Colombiana S.A.", "Consultoría Bogotana S.A.S."},
		CustomerNames: []string{"Alimentos del Caribe S.A.S.", "Construcciones del Eje S.A.S.", "Ingeniería Orinoquia Ltda.", "Comercio Medellín S.A.S.", "Servicios Cali S.A."},
		Cities:        []cityData{{Name: "Bogotá", Region: "Cundinamarca", Code: "110111"}, {Name: "Medellín", Region: "Antioquia", Code: "050001"}, {Name: "Cali", Region: "Valle del Cauca", Code: "760001"}, {Name: "Barranquilla", Region: "Atlántico", Code: "080001"}},
		Streets:       []string{"Carrera 7", "Calle 72", "Avenida El Dorado", "Carrera 13", "Calle 100"},
		Products:      []itemData{{Name: "Equipo de cómputo", Price: "3500000.00"}, {Name: "Material de construcción", Price: "1200000.00"}, {Name: "Maquinaria industrial", Price: "8500000.00"}},
		Services:      []itemData{{Name: "Consultoría empresarial", Price: "250000.00", Unit: "h"}, {Name: "Desarrollo de software", Price: "180000.00", Unit: "h"}, {Name: "Auditoría fiscal", Price: "350000.00", Unit: "h"}},
		PaymentKey:    pay.MeansKeyCreditTransfer,
	},
	"DE": {
		SupplierNames: []string{"Müller & Schmidt GmbH", "Technologie Innovationen AG", "Fischer Elektronik GmbH", "Weber Maschinenbau GmbH", "Schneider Logistik AG"},
		CustomerNames: []string{"Deutsche Handelskompanie GmbH", "Bayern Automotive GmbH", "Berliner Technologie GmbH", "Hamburg Port Services AG", "Schwaben Engineering GmbH"},
		Cities:        []cityData{{Name: "Berlin", Region: "Berlin", Code: "10115"}, {Name: "München", Region: "Bayern", Code: "80331"}, {Name: "Hamburg", Region: "Hamburg", Code: "20095"}, {Name: "Frankfurt am Main", Region: "Hessen", Code: "60311"}, {Name: "Köln", Region: "Nordrhein-Westfalen", Code: "50667"}},
		Streets:       []string{"Friedrichstraße", "Kurfürstendamm", "Maximilianstraße", "Bahnhofstraße", "Hauptstraße"},
		Products:      []itemData{{Name: "Laptop Business Pro", Price: "1299.00"}, {Name: "Bürostuhl Ergonomisch", Price: "449.00"}, {Name: "Netzwerk Switch", Price: "279.00"}},
		Services:      []itemData{{Name: "IT-Beratung", Price: "120.00", Unit: "h"}, {Name: "Softwareentwicklung", Price: "95.00", Unit: "h"}, {Name: "Projektmanagement", Price: "110.00", Unit: "h"}, {Name: "Cloud-Infrastruktur", Price: "499.00"}},
		PaymentKey:    sepa,
		IBANPrefix:    "DE",
	},
	"DK": {
		SupplierNames: []string{"Nordisk Teknologi ApS", "Københavns Rådgivning A/S", "Jyllands Maskinværksted ApS", "Skandinavisk Software A/S", "Fynske Metalindustri ApS"},
		CustomerNames: []string{"Sjællandsk Fødevarer A/S", "Østersøens Logistik ApS", "Bornholms Konsulenter ApS", "Aarhus Handel A/S", "Odense Elektronik ApS"},
		Cities:        []cityData{{Name: "København", Region: "Hovedstaden", Code: "1050"}, {Name: "Aarhus", Region: "Midtjylland", Code: "8000"}, {Name: "Odense", Region: "Syddanmark", Code: "5000"}, {Name: "Aalborg", Region: "Nordjylland", Code: "9000"}},
		Streets:       []string{"Vesterbrogade", "Nørregade", "Østergade", "Kongens Nytorv", "Frederiksberggade"},
		Products:      []itemData{{Name: "Vindmøllekomponent", Price: "8500.00"}, {Name: "Kontormøbler", Price: "4200.00"}, {Name: "Elektronisk udstyr", Price: "3100.00"}},
		Services:      []itemData{{Name: "IT-rådgivning", Price: "950.00", Unit: "h"}, {Name: "Softwareudvikling", Price: "850.00", Unit: "h"}, {Name: "Projektledelse", Price: "1100.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "DK",
	},
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
	"FR": {
		SupplierNames: []string{"Aquitaine Mécanique SAS", "Parisienne de Conseil SARL", "Lyonnaise d'Informatique SA", "Bretagne Agroalimentaire SAS", "Normandie Industries SA"},
		CustomerNames: []string{"Provence Logistique SARL", "Alsacienne de Construction SAS", "Bordeaux Solutions Numériques SARL", "Côte d'Azur Technologies SA", "Loire Valley Commerce SAS"},
		Cities:        []cityData{{Name: "Paris", Region: "Île-de-France", Code: "75008"}, {Name: "Lyon", Region: "Auvergne-Rhône-Alpes", Code: "69002"}, {Name: "Marseille", Region: "Provence-Alpes-Côte d'Azur", Code: "13001"}, {Name: "Toulouse", Region: "Occitanie", Code: "31000"}, {Name: "Bordeaux", Region: "Nouvelle-Aquitaine", Code: "33000"}},
		Streets:       []string{"Rue du Faubourg Saint-Honoré", "Avenue des Champs-Élysées", "Boulevard Haussmann", "Rue de Rivoli", "Rue Pasteur"},
		Products:      []itemData{{Name: "Équipement industriel", Price: "2400.00"}, {Name: "Mobilier de bureau", Price: "890.00"}, {Name: "Matériel informatique", Price: "1350.00"}},
		Services:      []itemData{{Name: "Conseil en stratégie", Price: "150.00", Unit: "h"}, {Name: "Développement logiciel", Price: "110.00", Unit: "h"}, {Name: "Audit financier", Price: "180.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "FR",
	},
	"GB": {
		SupplierNames: []string{"Thornfield Engineering Ltd", "Blackmore & Whitfield Consulting Ltd", "Pennine Manufacturing Ltd", "Highland Software Solutions Ltd", "Brightstone Logistics PLC"},
		CustomerNames: []string{"Cambridge Analytics Ltd", "Cotswold Food Suppliers Ltd", "Meridian Technical Services Ltd", "Thames Valley Instruments Ltd", "Yorkshire Steel Works Ltd"},
		Cities:        []cityData{{Name: "London", Region: "England", Code: "EC2M 1NH"}, {Name: "Manchester", Region: "England", Code: "M1 4BT"}, {Name: "Edinburgh", Region: "Scotland", Code: "EH1 1YZ"}, {Name: "Birmingham", Region: "England", Code: "B1 1BB"}, {Name: "Bristol", Region: "England", Code: "BS1 4DJ"}},
		Streets:       []string{"King William Street", "Deansgate", "Princes Street", "Colmore Row", "Queen Square"},
		Products:      []itemData{{Name: "Precision instruments", Price: "1850.00"}, {Name: "Office equipment", Price: "620.00"}, {Name: "Industrial sensors", Price: "940.00"}},
		Services:      []itemData{{Name: "Management consulting", Price: "150.00", Unit: "h"}, {Name: "Software engineering", Price: "120.00", Unit: "h"}, {Name: "Financial advisory", Price: "200.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "GB",
	},
	"GR": {
		SupplierNames: []string{"Aigaio Technologies I.K.E.", "Olympiaki Consulting A.E.", "Athinaiki Viomichania E.P.E.", "Thessaloniki Logistics I.K.E.", "Mesogeiaka Michanimata A.E."},
		CustomerNames: []string{"Kritiki Pliroforiki I.K.E.", "Peloponnisiaki Kataskevastiki E.P.E.", "Nisiotiki Emporiki I.K.E.", "Makedonia Trade A.E.", "Epirus Engineering E.P.E."},
		Cities:        []cityData{{Name: "Athens", Region: "Attica", Code: "10564"}, {Name: "Thessaloniki", Region: "Central Macedonia", Code: "54622"}, {Name: "Patras", Region: "Western Greece", Code: "26221"}, {Name: "Heraklion", Region: "Crete", Code: "71201"}},
		Streets:       []string{"Ermou Street", "Vasilissis Sofias Avenue", "Stadiou Street", "Panepistimiou Street", "Akademias Street"},
		Products:      []itemData{{Name: "Marine equipment", Price: "3200.00"}, {Name: "Olive oil processing machine", Price: "5800.00"}, {Name: "Solar panel system", Price: "4500.00"}},
		Services:      []itemData{{Name: "Shipping consultation", Price: "100.00", Unit: "h"}, {Name: "Tourism management", Price: "80.00", Unit: "h"}, {Name: "Software development", Price: "90.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "GR",
	},
	"IE": {
		SupplierNames: []string{"Claddagh Software Ltd", "Emerald Isle Engineering Ltd", "Liffey Consulting Ltd", "Shannon Logistics Ltd", "Galway Bay Foods Ltd"},
		CustomerNames: []string{"Wicklow Precision Manufacturing Ltd", "Cork Harbour Technologies Ltd", "Boyne Valley Supplies Ltd", "Dingle Peninsula Services Ltd", "Kildare Instruments Ltd"},
		Cities:        []cityData{{Name: "Dublin", Region: "Leinster", Code: "D02 YX88"}, {Name: "Cork", Region: "Munster", Code: "T12 W026"}, {Name: "Galway", Region: "Connacht", Code: "H91 E2K5"}, {Name: "Limerick", Region: "Munster", Code: "V94 T9PX"}},
		Streets:       []string{"Grafton Street", "Patrick Street", "Eyre Square", "O'Connell Street", "George's Street"},
		Products:      []itemData{{Name: "Pharmaceutical equipment", Price: "4200.00"}, {Name: "Dairy processing unit", Price: "6800.00"}, {Name: "Medical devices", Price: "3500.00"}},
		Services:      []itemData{{Name: "Financial consulting", Price: "160.00", Unit: "h"}, {Name: "Pharmaceutical research", Price: "200.00", Unit: "h"}, {Name: "IT services", Price: "130.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "IE",
	},
	"IN": {
		SupplierNames: []string{"Ganesh Precision Tools Pvt. Ltd.", "Tara Infosystems Pvt. Ltd.", "Bharat Heavy Fabricators Ltd.", "Kaveri Consulting Services LLP", "Deccan Software Solutions Pvt. Ltd."},
		CustomerNames: []string{"Lakshmi Food Processing Pvt. Ltd.", "Narmada Logistics Pvt. Ltd.", "Ashoka Manufacturing Pvt. Ltd.", "Godavari Chemicals Ltd.", "Himalaya Textiles Pvt. Ltd."},
		Cities:        []cityData{{Name: "Mumbai", Region: "Maharashtra", Code: "400001"}, {Name: "New Delhi", Region: "Delhi", Code: "110001"}, {Name: "Bengaluru", Region: "Karnataka", Code: "560001"}, {Name: "Chennai", Region: "Tamil Nadu", Code: "600001"}, {Name: "Hyderabad", Region: "Telangana", Code: "500001"}},
		Streets:       []string{"Mahatma Gandhi Road", "Jawaharlal Nehru Marg", "Sardar Patel Road", "Rajaji Salai", "Subhash Chandra Bose Road"},
		Products:      []itemData{{Name: "Industrial machinery", Price: "250000.00"}, {Name: "Textile equipment", Price: "180000.00"}, {Name: "IT hardware", Price: "95000.00"}},
		Services:      []itemData{{Name: "IT consulting", Price: "5000.00", Unit: "h"}, {Name: "Software development", Price: "4000.00", Unit: "h"}, {Name: "Engineering design", Price: "3500.00", Unit: "h"}},
		PaymentKey:    pay.MeansKeyCreditTransfer,
	},
	"IT": {
		SupplierNames: []string{"Officine Meccaniche Lombarde S.r.l.", "Consulenza Romana S.p.A.", "Alimentari del Mezzogiorno S.r.l.", "Logistica Adriatica S.r.l.", "Tecnologie Fiorentine S.p.A."},
		CustomerNames: []string{"Industria Torinese S.r.l.", "Software Veneto S.r.l.", "Costruzioni Emiliane S.p.A.", "Tessile Milanese S.r.l.", "Ceramiche Toscane S.r.l."},
		Cities:        []cityData{{Name: "Roma", Region: "Lazio", Code: "00186"}, {Name: "Milano", Region: "Lombardia", Code: "20121"}, {Name: "Napoli", Region: "Campania", Code: "80133"}, {Name: "Torino", Region: "Piemonte", Code: "10121"}, {Name: "Firenze", Region: "Toscana", Code: "50123"}},
		Streets:       []string{"Via Roma", "Corso Vittorio Emanuele II", "Via Giuseppe Garibaldi", "Via Dante Alighieri", "Via Guglielmo Marconi"},
		Products:      []itemData{{Name: "Macchinario industriale", Price: "4500.00"}, {Name: "Arredamento ufficio", Price: "1200.00"}, {Name: "Strumentazione elettronica", Price: "2800.00"}},
		Services:      []itemData{{Name: "Consulenza aziendale", Price: "100.00", Unit: "h"}, {Name: "Sviluppo software", Price: "85.00", Unit: "h"}, {Name: "Progettazione ingegneristica", Price: "95.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "IT",
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
	"NL": {
		SupplierNames: []string{"Hollandse Werktuigbouw B.V.", "Amsterdamse Consultancy B.V.", "Rotterdamse Logistiek B.V.", "Utrechtse Software B.V.", "Brabantse Voedingsmiddelen N.V."},
		CustomerNames: []string{"Zeeuwse Techniek B.V.", "Groningse Machinefabriek B.V.", "Haagse Ingenieursdiensten B.V.", "Limburgse Chemie N.V.", "Friese Handel B.V."},
		Cities:        []cityData{{Name: "Amsterdam", Region: "Noord-Holland", Code: "1012"}, {Name: "Rotterdam", Region: "Zuid-Holland", Code: "3011"}, {Name: "Den Haag", Region: "Zuid-Holland", Code: "2511"}, {Name: "Utrecht", Region: "Utrecht", Code: "3511"}, {Name: "Eindhoven", Region: "Noord-Brabant", Code: "5611"}},
		Streets:       []string{"Keizersgracht", "Coolsingel", "Lange Voorhout", "Oudegracht", "Kalverstraat"},
		Products:      []itemData{{Name: "Waterbeheerinstallatie", Price: "5600.00"}, {Name: "Kantoorinrichting", Price: "1800.00"}, {Name: "Industriële pomp", Price: "3400.00"}},
		Services:      []itemData{{Name: "Management advies", Price: "140.00", Unit: "h"}, {Name: "Software ontwikkeling", Price: "110.00", Unit: "h"}, {Name: "Logistiek beheer", Price: "95.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "NL",
	},
	"PL": {
		SupplierNames: []string{"Warszawskie Zakłady Precyzyjne Sp. z o.o.", "Krakowska Firma Doradcza Sp. z o.o.", "Gdańskie Technologie S.A.", "Poznańska Logistyka Sp. z o.o.", "Wrocławski Serwis Informatyczny Sp. z o.o."},
		CustomerNames: []string{"Śląskie Artykuły Spożywcze Sp. z o.o.", "Łódzkie Zakłady Mechaniczne S.A.", "Lubelska Firma Budowlana Sp. z o.o.", "Szczecin Handel Sp. z o.o.", "Bydgoska Elektronika Sp. z o.o."},
		Cities:        []cityData{{Name: "Warszawa", Region: "Mazowieckie", Code: "00-001"}, {Name: "Kraków", Region: "Małopolskie", Code: "30-001"}, {Name: "Wrocław", Region: "Dolnośląskie", Code: "50-001"}, {Name: "Gdańsk", Region: "Pomorskie", Code: "80-001"}, {Name: "Poznań", Region: "Wielkopolskie", Code: "60-001"}},
		Streets:       []string{"ul. Marszałkowska", "ul. Floriańska", "ul. Świdnicka", "ul. Długa", "ul. Piotrkowska"},
		Products:      []itemData{{Name: "Maszyna przemysłowa", Price: "12000.00"}, {Name: "Sprzęt komputerowy", Price: "4500.00"}, {Name: "Wyposażenie biurowe", Price: "3200.00"}},
		Services:      []itemData{{Name: "Doradztwo biznesowe", Price: "350.00", Unit: "h"}, {Name: "Rozwój oprogramowania", Price: "280.00", Unit: "h"}, {Name: "Audyt finansowy", Price: "400.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "PL",
	},
	"PT": {
		SupplierNames: []string{"Indústrias do Tejo Lda.", "Consultores de Lisboa Lda.", "Tecnologias do Porto S.A.", "Alimentos do Algarve Lda.", "Logística Atlântica Lda."},
		CustomerNames: []string{"Engenharia do Minho Lda.", "Fábrica Beirã S.A.", "Soluções Digitais do Douro Lda.", "Comércio Alentejano Lda.", "Têxteis de Guimarães S.A."},
		Cities:        []cityData{{Name: "Lisboa", Region: "Lisboa", Code: "1100-148"}, {Name: "Porto", Region: "Porto", Code: "4000-007"}, {Name: "Braga", Region: "Braga", Code: "4700-006"}, {Name: "Coimbra", Region: "Coimbra", Code: "3000-001"}, {Name: "Faro", Region: "Faro", Code: "8000-001"}},
		Streets:       []string{"Rua Augusta", "Avenida dos Aliados", "Rua de Santa Catarina", "Rua Ferreira Borges", "Rua do Comércio"},
		Products:      []itemData{{Name: "Equipamento industrial", Price: "3800.00"}, {Name: "Mobiliário de escritório", Price: "1100.00"}, {Name: "Material eletrónico", Price: "2200.00"}},
		Services:      []itemData{{Name: "Consultoria de gestão", Price: "85.00", Unit: "h"}, {Name: "Desenvolvimento de software", Price: "70.00", Unit: "h"}, {Name: "Auditoria fiscal", Price: "95.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "PT",
	},
	"SE": {
		SupplierNames: []string{"Nordströms Verkstadsteknik AB", "Stockholms Konsultbyrå AB", "Göteborgs Livsmedel AB", "Malmö Digitala Lösningar AB", "Norrlands Maskinfabrik AB"},
		CustomerNames: []string{"Uppsalas Logistiktjänster AB", "Skånska Byggkonsulter AB", "Östgöta Industriservice AB", "Värmlands Teknik AB", "Hallands Handel AB"},
		Cities:        []cityData{{Name: "Stockholm", Region: "Stockholms län", Code: "111 21"}, {Name: "Göteborg", Region: "Västra Götalands län", Code: "411 01"}, {Name: "Malmö", Region: "Skåne län", Code: "211 18"}, {Name: "Uppsala", Region: "Uppsala län", Code: "753 10"}},
		Streets:       []string{"Drottninggatan", "Kungsgatan", "Storgatan", "Vasagatan", "Sveavägen"},
		Products:      []itemData{{Name: "Industriverktyg", Price: "8500.00"}, {Name: "Kontorsutrustning", Price: "4200.00"}, {Name: "Elektronikkomponenter", Price: "3100.00"}},
		Services:      []itemData{{Name: "IT-konsulttjänster", Price: "1100.00", Unit: "h"}, {Name: "Mjukvaruutveckling", Price: "950.00", Unit: "h"}, {Name: "Projektledning", Price: "1200.00", Unit: "h"}},
		PaymentKey:    sepa,
		IBANPrefix:    "SE",
	},
	"SG": {
		SupplierNames: []string{"Temasek Engineering Pte. Ltd.", "Lion City Logistics Pte. Ltd.", "Raffles Consulting Services Pte. Ltd.", "Marina Bay Technologies Pte. Ltd.", "Merlion Precision Manufacturing Pte. Ltd."},
		CustomerNames: []string{"Orchid Food Industries Pte. Ltd.", "Harbour Front Software Pte. Ltd.", "Jurong Heavy Equipment Pte. Ltd.", "Sentosa Trading Pte. Ltd.", "Changi Electronics Pte. Ltd."},
		Cities:        []cityData{{Name: "Singapore", Code: "048616"}, {Name: "Singapore", Code: "179101"}, {Name: "Singapore", Code: "609916"}},
		Streets:       []string{"Robinson Road", "Raffles Place", "Shenton Way", "Jurong Gateway Road", "Ang Mo Kio Street"},
		Products:      []itemData{{Name: "Electronic components", Price: "2400.00"}, {Name: "Precision machinery", Price: "8500.00"}, {Name: "Office equipment", Price: "1600.00"}},
		Services:      []itemData{{Name: "Financial consulting", Price: "250.00", Unit: "h"}, {Name: "Software development", Price: "180.00", Unit: "h"}, {Name: "Logistics management", Price: "150.00", Unit: "h"}},
		PaymentKey:    pay.MeansKeyCreditTransfer,
	},
	"US": {
		SupplierNames: []string{"Ironwood Manufacturing LLC", "Summit Ridge Consulting Inc.", "Bayshore Logistics Corp.", "Pinehurst Technologies Inc.", "Copper Basin Engineering LLC"},
		CustomerNames: []string{"Prairie Systems Inc.", "Redstone Food Distributors LLC", "Lakeview Software Solutions Inc.", "Granite State Instruments Inc.", "Magnolia Health Services LLC"},
		Cities:        []cityData{{Name: "New York", Region: "NY", Code: "10017"}, {Name: "Los Angeles", Region: "CA", Code: "90012"}, {Name: "Chicago", Region: "IL", Code: "60601"}, {Name: "Houston", Region: "TX", Code: "77002"}, {Name: "Phoenix", Region: "AZ", Code: "85004"}},
		Streets:       []string{"Madison Avenue", "Wilshire Boulevard", "North Michigan Avenue", "Main Street", "Central Avenue"},
		Products:      []itemData{{Name: "Industrial sensor kit", Price: "1200.00"}, {Name: "Office workstation", Price: "850.00"}, {Name: "Server hardware", Price: "3400.00"}},
		Services:      []itemData{{Name: "Management consulting", Price: "200.00", Unit: "h"}, {Name: "Software engineering", Price: "175.00", Unit: "h"}, {Name: "Legal advisory", Price: "300.00", Unit: "h"}},
		PaymentKey:    pay.MeansKeyCreditTransfer,
	},
}
