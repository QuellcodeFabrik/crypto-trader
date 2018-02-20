// The integration package implements the available crypto-currency exchange
// platforms and offers convenience functions to make use of their
// capabilities.
package integrations

import (
    "log"
    "os"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "../helper"
    "strings"
)

type Bitstamp struct {
    supportedCurrencies []string
}

func (bitstamp *Bitstamp) Init() {
    log.Println("integrations::Bitstamp::Init()")

    bitstamp.supportedCurrencies = []string{
        "BTC", "XRP", "LTC", "ETH", "BCH" }

    log.Println("Supported currencies:", bitstamp.supportedCurrencies)
}

func (bitstamp *Bitstamp) GetAvailableCurrencies() []string {
    return bitstamp.supportedCurrencies
}

func (bitstamp *Bitstamp) GetCurrencyValue(currency string) CurrencySnapshot {
    log.Println("integrations::Bitstamp::GetCurrencyValue()")

    currency = strings.ToLower(currency) + "eur"

    if ! helper.IsElementInArray(currency, bitstamp.supportedCurrencies) {
        log.Println("This currency is not supported: " + currency)
        os.Exit(1)
    }

    var snapshot CurrencySnapshot
    response, err := http.Get("https://www.bitstamp.net/api/v2/ticker/" + currency)
    if err != nil {
        log.Println(err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            log.Println(err)
            os.Exit(1)
        }
        json.Unmarshal(contents, &snapshot)
    }
    return snapshot
}