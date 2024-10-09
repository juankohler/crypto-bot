package models

import (
	"encoding/json"
	"strings"

	"github.com/juankohler/crypto-bot/libs/go/errors"
)

var (
	ErrInvalidCurrency = errors.Define("currency.invalid")
)

type Currency string

const (
	ARS Currency = "ARS"
	BRL Currency = "BRL"
	CLP Currency = "CLP"
	COP Currency = "COP"
	EUR Currency = "EUR"
	MXN Currency = "MXN"
	PEN Currency = "PEN"
	USD Currency = "USD"
	UYU Currency = "UYU"
)

var allowedCurrencies = map[string]Currency{
	ARS.String(): ARS,
	BRL.String(): BRL,
	CLP.String(): CLP,
	COP.String(): COP,
	EUR.String(): EUR,
	MXN.String(): MXN,
	PEN.String(): PEN,
	USD.String(): USD,
	UYU.String(): UYU,
}

func NewCurrency(v string) (Currency, error) {
	if v, ok := allowedCurrencies[strings.ToUpper(v)]; ok {
		return v, nil
	}

	return "", errors.New(
		ErrInvalidCurrency,
		"invalid currency",
		errors.WithMetadata("currency", v),
	)
}

func (c Currency) String() string {
	return string(c)
}

// Serialization
func (c Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Currency) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	bCurrency, err := NewCurrency(s)
	if err != nil {
		return err
	}

	*c = bCurrency

	return nil
}
