package models

import (
	"encoding/json"
	"strings"

	"github.com/juankohler/crypto-bot/libs/go/errors"
)

var (
	ErrInvalidCountry = errors.Define("country.invalid")
)

type Country string

const (
	CO  Country = "CO"
	MX  Country = "MX"
	AR  Country = "AR"
	UY  Country = "UY"
	PE  Country = "PE"
	ES  Country = "ES"
	USA Country = "USA"
	CL  Country = "CL"
	BR  Country = "BR"
	WW  Country = "WW"
	LAN Country = "LAN"
	SA  Country = "SA"
	VE  Country = "VE"
)

var allowedCountries = map[string]Country{
	CO.String():  CO,
	MX.String():  MX,
	AR.String():  AR,
	UY.String():  UY,
	PE.String():  PE,
	ES.String():  ES,
	USA.String(): USA,
	CL.String():  CL,
	BR.String():  BR,
	WW.String():  WW,
	LAN.String(): LAN,
	SA.String():  SA,
	VE.String():  VE,
}

func NewCountry(v string) (Country, error) {
	if v, ok := allowedCountries[strings.ToUpper(v)]; ok {
		return v, nil
	}

	return "", errors.New(
		ErrInvalidCountry,
		"invalid country",
		errors.WithMetadata("country", v),
	)
}

func (c Country) String() string {
	return string(c)
}

// Serialization
func (c Country) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Country) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	bCountry, err := NewCountry(s)
	if err != nil {
		return err
	}

	*c = bCountry

	return nil
}
