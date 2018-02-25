// The aggregator collects data from the given data sources and stores this
// data in the connected database.
package aggregator

import (
    "log"
    "time"
    "../integrations"
    db "../database"
    "encoding/json"
)

// Init initializes the aggregator by starting different timers. Those timers
// are used to collect data from the given crypto-currency exchange platform
// and to store this data in the connected database.
func Init(exchange integrations.ExchangeIntegration) {
    log.Println("aggregator::Init()")

    exchange.Init()

    channel := make(chan string)

    go startDataAggregationTimer(channel, exchange)
    go startStatisticsTimer(channel)

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

func startDataAggregationTimer(c chan<- string, exchange integrations.ExchangeIntegration) {
    log.Println("aggregator::startDataAggregationTimer()")

    availableCurrencies := exchange.GetSupportedCurrencies()
    numberOfCurrencies := len(availableCurrencies)
    lastCurrencyUsed := 0

    tick := time.Tick(6 * time.Second)

    for {
        select {
        case <-tick:
            snapshot := exchange.GetCurrencySnapshot(availableCurrencies[lastCurrencyUsed])

            if err := db.StoreSnapshot(availableCurrencies[lastCurrencyUsed], snapshot); err != nil {
                log.Printf("Could not store snapshot: %s", err.Error())
                return
            }

            out, err := json.Marshal(snapshot)
            if err != nil {
                panic (err)
            }
            c <- availableCurrencies[lastCurrencyUsed] + ": " + string(out)

            lastCurrencyUsed = (lastCurrencyUsed + 1) % numberOfCurrencies
        }
    }
}

func startStatisticsTimer(c chan<- string) {
    log.Println("aggregator::startStatisticsTimer()")

    tick := time.Tick(10 * time.Second)

    for {
        select {
        case <-tick:
            c <- "Statistics"
            // TODO implement and trigger statistics evaluation here
        }
    }
}