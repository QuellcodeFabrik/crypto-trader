// The analyzer takes data from the connected database and
// derives secondary data from it. The secondary data will
// be used to take buying and selling decisions.
package analyzer

import (
    db "../database"
    "log"
)

const tradeMargin = 0.05

// GetOptimalPurchasePrice retrieves historic data for the given currency from
// the database and calculates characteristics that help to identify the
// optimal purchase price for that currency to make significant gains in the
// near future.
func GetOptimalPurchasePrice(positions []db.TradingPosition) float64 {
    var referenceValue float64
    for _, position := range positions {
        if referenceValue == 0 || position.Value < referenceValue {
            referenceValue = position.Value
        }
    }

    if referenceValue == 0 {
        // TODO get reference value from historical data
    }

    log.Printf("The reference value is set to %f.", referenceValue)

    // TODO calculate the risk factor and take into account the following:
    // - volatility
    // - recent price development
    riskFactor := 1.0

    return referenceValue * (1 - tradeMargin * riskFactor)
}

// IsCurrencyEligibleForDivestment retrieves historic data for the given currency
// from the database and calculates characteristics that help to identify if now
// is the right time to sell it for the given value.
func IsPositionEligibleForDivestment(position db.TradingPosition, currentValue float64) bool {

    // TODO calculate the risk factor and take into account the following:
    // - volatility
    // - recent price development
    riskFactor := 1.0

    return (position.Value * (1.0 + tradeMargin * riskFactor)) <= currentValue
}
