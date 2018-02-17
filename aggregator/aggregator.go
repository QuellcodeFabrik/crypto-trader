package aggregator

import (
    "log"
    "time"
    "../integrations"
    db "../database"
    "encoding/json"
)

func Init(integration integrations.ExchangeIntegration) {
    log.Println("aggregator::Init()")

    integration.Init()

    channel := make(chan string)

    go startDataAggregationTimer(channel, integration)
    go startStatisticsTimer(channel)

    for {
        select {
        case message := <- channel:
            log.Println(message)
        }
    }
}

func startDataAggregationTimer(c chan<- string, exchange integrations.ExchangeIntegration) {
    log.Println("aggregator::startDataAggregationTimer()")

    availableCurrencies := exchange.GetAvailableCurrencies()
    numberOfCurrencies := len(availableCurrencies)
    lastCurrencyUsed := 0

    tick := time.Tick(5 * time.Second)

    for {
        select {
        case <-tick:
            snapshot := exchange.GetCurrencyValue(availableCurrencies[lastCurrencyUsed])

            if err := db.StoreSnapshot(availableCurrencies[lastCurrencyUsed], &snapshot); err != nil {
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
            // do something here
        }
    }
}