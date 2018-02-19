// The analyzer takes data from the connected database and
// derives secondary data from it. The secondary data will
// be used to take buying and selling decisions.
package analyzer

import (
    db "../database"
)

const TRADE_MARGIN = 0.05

// GetOptimalPurchasePrice retrieves historic data for the given currency from
// the database and calculates characteristics that help to identify the
// optimal purchase price for that currency to make significant gains in the
// near future.
func GetOptimalPurchasePrice(currency db.CryptoCurrency, positions []db.InvestmentPosition) float64 {
    return 0.0
}

// IsCurrencyEligibleForDivestment retrieves historic data for the given currency
// from the database and calculates characteristics that help to identify if now
// is the right time to sell it for the given value.
func IsCurrencyEligibleForDivestment(currency db.CryptoCurrency, value float64) bool {
    return false
}
