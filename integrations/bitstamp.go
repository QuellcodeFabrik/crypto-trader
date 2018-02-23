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
    "encoding/hex"
    "../helper"
    "strings"
    "crypto/hmac"
    "crypto/sha256"
)

type Bitstamp struct {
    supportedCurrencies []string
}

func (bitstamp *Bitstamp) Init() {
    log.Println("integrations::Bitstamp::Init()")

    bitstamp.supportedCurrencies = []string{
        "BTC", "XRP", "LTC", "ETH", "BCH" }

    log.Println("Supported currencies:", bitstamp.supportedCurrencies)

    nonce := int64(1) // TODO create real nonce
    signature := createSignature(nonce, ApiAccessData)

    log.Printf("Bitstamp signature: %s", signature)
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

// createSignature takes a nonce, a customerId and an ApiAccess struct to
// create and return the signature string out of it that is necessary to send
// private requests to the Bitstamp API.
func createSignature(nonce int64, apiCredentials ApiAccess) string {
    message := []byte(string(nonce) + apiCredentials.CustomerId + apiCredentials.ApiKey)
    mac := hmac.New(sha256.New, []byte(apiCredentials.ApiSecret))
    mac.Write(message)
    return strings.ToUpper(hex.EncodeToString(message))
}