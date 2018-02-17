// The executor package handles all buying and selling transactions while
// using the analyzer package to prepare the decision for a specific
// transaction.
package executor

import (
    "log"
    db "../database"
    "../integrations"
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

func checkPurchasePotential(currency db.CryptoCurrency) float32 {
    log.Println("executor::checkPurchasePotential()")
    return 0
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