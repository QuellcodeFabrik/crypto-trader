// The executor package handles all buying and selling transactions while
// using the analyzer package to prepare the decision for a specific
// transaction.
package executor

import (
    "log"
    db "../database"
    "../integrations"
    "../analyzer"
    "time"
)

func Init(exchange integrations.ExchangeIntegration) {
    log.Println("executor::Init()")

    exchange.Init()

    channel := make(chan string)

    go startExecutionTimer(channel, exchange)

    for {
        select {
        case message := <- channel:
            log.Println(message)
        }
    }
}

//
// Private functions
//

func startExecutionTimer(c chan<- string, exchange integrations.ExchangeIntegration) {
    log.Println("executor::startExecutionTimer()")

    tick := time.Tick(time.Minute)

    for {
        select {
        case <-tick:
            c <- "Execution"
            // TODO implement and trigger execution events here
        }
    }
}

func checkPurchasePotentialAndExecute(currency db.CryptoCurrency, snapshot integrations.CurrencySnapshot) {
    log.Println("executor::checkPurchasePotential()")

    // TODO
    // retrieve existing positions from database and check if there already
    // are positions for the given currency
    var existingPositions []db.TradingPosition

    existingPositions = append(existingPositions, db.TradingPosition{Id: 1, Amount: 1.0, Value: 10.0, Timestamp: time.Now().Unix()})
    existingPositions = append(existingPositions, db.TradingPosition{Id: 2, Amount: 1.0, Value: 100.0, Timestamp: time.Now().Unix()})
    existingPositions = append(existingPositions, db.TradingPosition{Id: 3, Amount: 1.0, Value: 120.0, Timestamp: time.Now().Unix()})

    if snapshot.Current > analyzer.GetOptimalPurchasePrice(existingPositions) {
        log.Printf("Purchase order for %s will be made now.", currency.Name)

        //
    }
}

func checkSalesPotential(currency db.CryptoCurrency) float32 {
    log.Println("executor::checkSalesPotential()")
    return 0
}

func executePurchaseOrder(currency db.CryptoCurrency) {
    log.Println("executor::executePurchaseOrder()")
}

func executeSalesOrder(currency db.CryptoCurrency) {
    log.Println("executor::executeSalesOrder()")
}