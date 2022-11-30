package util



const(
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupporredCurrency(currency string) bool{
	switch currency{
	case USD,EUR,CAD:
		return true
	}
	return false
}