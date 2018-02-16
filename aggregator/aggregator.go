package aggregator

import (
    "log"
    "time"
    "../integrations"
)

func Init(integration *integrations.Bitstamp) {
    log.Println("aggregator::Init()")

    integration.Init()
    currencySnapshot := integration.GetCurrencyValue()
    log.Printf("%+v\n", currencySnapshot)

    startDataAggregationTimer()
}

func startDataAggregationTimer() {
    timeout := time.After(60 * time.Second)
    tick := time.Tick(10 * time.Second)

    for {
        select {
        case <-timeout:
            log.Println("Timed out")
        case <-tick:
            doSomething()
        }
    }
}

func doSomething() {
    log.Println("aggregator::doSomething()")
}
