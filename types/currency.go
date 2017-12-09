package types

import "errors"

type Currency int32

// MarshalText implements encoding/text TextMarshaler interface.
func (c Currency) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements encoding/text TextUnmarshaler interface.
func (c *Currency) UnmarshalText(text []byte) error {
	val, ok := currenciesByAlpha[string(text)]
	if !ok {
		return errors.New("Invalid currency: " + string(text))
	}
	*c = val
	return nil
}

// String returns the three-letter currency code.
func (c Currency) String() string {
	return currencyInfos[c].alpha
}

// Decimals returns the number of decimals of this currency.
func (c Currency) Decimals() int {
	return currencyInfos[c].decimals
}

type currencyInfo struct {
	numeric  string
	alpha    string
	decimals int
}

var currencyInfos = []currencyInfo{
	{},
	{alpha: "EUR", numeric: "978", decimals: 2},
	{alpha: "USD", numeric: "840", decimals: 2},
	{alpha: "GBP", numeric: "826", decimals: 2},
}

var currenciesByAlpha map[string]Currency

func init() {
	currenciesByAlpha = make(map[string]Currency)
	for i, info := range currencyInfos {
		if info.alpha != "" {
			currenciesByAlpha[info.alpha] = Currency(i)
		}
	}
}
