package analyzer

import (
    db "../database"
    "testing"
    "time"
)

var mockPositions = []db.InvestmentPosition{
    { Id:1, Amount:1, Value:1.0, Timestamp:time.Now() },
    { Id:2, Amount:5, Value:3.0, Timestamp:time.Now() },
    { Id:3, Amount:5, Value:5.0, Timestamp:time.Now() }}

func TestGetOptimalPurchasePrice(t *testing.T) {
    optimalPurchasePrice := GetOptimalPurchasePrice(mockPositions)
    if optimalPurchasePrice != 0.95 {
        t.Errorf("The optimal purchase price calculation is wrong, got: %f, want: %f.",
            optimalPurchasePrice, 0.95)
    }
}

func TestIsPositionEligibleForDivestment(t *testing.T) {
    shallPositionBeSold := IsPositionEligibleForDivestment(mockPositions[1], 8.0)
    t.Log(shallPositionBeSold)

    if shallPositionBeSold == false {
        t.Fail()
    }

    shallPositionBeSold = IsPositionEligibleForDivestment(mockPositions[0], 1.1)
    t.Log(shallPositionBeSold)

    if shallPositionBeSold == true {
        t.Fail()
    }
}