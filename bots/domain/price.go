package domain

import (
	"github.com/shopspring/decimal"
)

type Price struct {
	BaseCurrency  string
	QuoteCurrency string
	Price         decimal.Decimal
}

func NewPrice(
	baseCurrency string,
	quoteCurrency string,
	price decimal.Decimal,
) (*Price, error) {
	entity := &Price{
		BaseCurrency:  baseCurrency,
		QuoteCurrency: quoteCurrency,
		Price:         price,
	}

	return entity, nil
}
