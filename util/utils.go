package util

// TODO: Refactor utils

type Currency struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Id string `json:"starter"`
}

var Currencies = map[string]Currency{
	"USD": {
		Code:    "USD",
		Name:    "United States Dollar",
		Id: "1",
	},
	"ARS": {
		Code:    "ARS",
		Name:    "Argentine Peso",
		Id: "2",
	},
}

func IsValidCurrency(currency string) bool {
	_, ok := Currencies[currency]
	return ok
}
