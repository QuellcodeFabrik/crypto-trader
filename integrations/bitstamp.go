package integrations

import (
    "log"
    "os"
    "io/ioutil"
    "net/http"
    "net/url"
    "encoding/json"
)

const PROXY = "" // http://proxy.intra.dmc.de:3128"

type CurrencySnapshot struct {
    Low         float64   `json:"low,string"`
    High        float64   `json:"high,string"`
    Current     float64   `json:"last,string"`
    Timestamp   int64     `json:"timestamp,string"`
    Average     float64   `json:"vwap,string"`
}

type Bitstamp struct {
    supportedCurrencies []string
}

func (bitstamp *Bitstamp) Init() {
    log.Println("Bitstamp::Init()")

    bitstamp.supportedCurrencies = []string{
        "btcusd", "btceur",
        "xrpusd", "xrpeur",
        "ltcusd", "ltceur",
        "ethusd", "etheur",
        "bchusd", "bcheur" }

    log.Println("Supported currencies:", bitstamp.supportedCurrencies)

    if PROXY != "" {
        proxyUrl, _ := url.Parse(PROXY)
        http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
    }
}

func (bitstamp *Bitstamp) GetCurrencyValue() CurrencySnapshot {
    log.Println("Bitstamp::GetCurrencyValue()")

    var snapshot CurrencySnapshot
    response, err := http.Get("https://www.bitstamp.net/api/v2/ticker/xrpeur")
    if err != nil {
        log.Println(err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        log.Println(string(contents))
        if err != nil {
            log.Println(err)
            os.Exit(1)
        }
        json.Unmarshal(contents, &snapshot)
    }
    return snapshot
}