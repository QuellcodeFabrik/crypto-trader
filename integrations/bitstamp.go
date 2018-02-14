package integrations

import (
    "log"
    "os"
    "io/ioutil"
    "net/http"
    "net/url"
)

const PROXY = "http://proxy.intra.dmc.de:3128"

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

    proxyUrl, _ := url.Parse(PROXY)
    http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
}

func (bitstamp *Bitstamp) GetCurrencyValue() {
    log.Println("Bitstamp::GetCurrencyValue()")

    response, err := http.Get("https://www.bitstamp.net/api/v2/ticker/xrpeur")
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
        log.Println(string(contents))
    }
}