package mock

import (
	"math/rand/v2"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func generateLines(r *rand.Rand, cfg *regimeConfig, count int) []*bill.Line {
	lines := make([]*bill.Line, count)
	for i := range lines {
		lines[i] = generateLine(r, cfg)
	}
	return lines
}

func generateLine(r *rand.Rand, cfg *regimeConfig) *bill.Line {
	// Mix products and services.
	var item itemData
	if r.IntN(2) == 0 && len(cfg.Products) > 0 {
		item = pick(r, cfg.Products)
	} else {
		item = pick(r, cfg.Services)
	}

	quantity := num.MakeAmount(int64(r.IntN(20)+1), 0)
	price, _ := num.AmountFromString(item.Price)

	lineItem := &org.Item{
		Name:  item.Name,
		Price: &price,
	}

	if item.Unit != "" {
		lineItem.Unit = org.Unit(item.Unit)
	}

	if cfg.ItemExt != nil {
		lineItem.Ext = cfg.ItemExt(r)
	}

	// Assign tax rate (mostly standard, occasionally reduced).
	rate := cfg.TaxRates[0]
	if len(cfg.TaxRates) > 1 && r.IntN(4) == 0 {
		rate = cfg.TaxRates[r.IntN(len(cfg.TaxRates))]
	}

	line := &bill.Line{
		Quantity: quantity,
		Item:     lineItem,
		Taxes: tax.Set{
			{
				Category: rate.Category,
				Rate:     rate.Rate,
			},
		},
	}

	// Occasionally add a line discount.
	if r.IntN(5) == 0 {
		pct, _ := num.PercentageFromString("10%")
		line.Discounts = []*bill.LineDiscount{
			{
				Percent: &pct,
				Reason:  "Discount",
			},
		}
	}

	return line
}
