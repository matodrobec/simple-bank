package util

const (
	USD = "USD"
	EUR = "EUR"
	SK  = "SK"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, SK:
		return true
	}
	return false
}
