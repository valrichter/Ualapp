package util

var Currencies = map[string]string{
	"USD": "USD",
	"ARS": "ARS",
}

func IsValidCurrency(currency string) bool {
	_, ok := Currencies[currency]
	return ok
}
