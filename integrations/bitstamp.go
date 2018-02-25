// The integration package implements the available crypto-currency exchange
// platforms and offers convenience functions to make use of their
// capabilities.
package integrations

import (
    "log"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "encoding/hex"
    "../helper"
    "strings"
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
    "bytes"
    "strconv"
    "time"
    "os"
)

type Bitstamp struct {
    supportedCurrencies []string
}

// CreateBuyOrder BUY LIMIT ORDER
func (bitstamp *Bitstamp) CreateBuyOrder(amount float32, value float64) *Order {
    // TODO implement
    // https://www.bitstamp.net/api/v2/buy/{currency_pair}/
    //
    // amount	    Amount
    // price	    Price
    // limit_price	If the order gets executed, a new sell order will be placed, with "limit_price" as its price.

    return nil
}

// CreateSellOrder SELL LIMIT ORDER
func (bitstamp *Bitstamp) CreateSellOrder(amount float32, value float64) *Order {
    // TODO implement
    // https://www.bitstamp.net/api/v2/sell/{currency_pair}/
    //
    // amount	    Amount
    // price	    Price
    // limit_price	If the order gets executed, a new buy order will be placed, with "limit_price" as its price

    return nil
}

// GetAccountBalance returns the overall account balance or nil in case
// something went wrong with the request.
func (bitstamp *Bitstamp) GetAccountBalance() *AccountBalance {
    log.Println("integrations::Bitstamp::GetAccountBalance()")

    nonce := int(time.Now().Unix())
    signature := createSignature(nonce, ApiAccessData)

    url := "https://www.bitstamp.net/api/v2/balance/"

    postBody := []byte(createAuthRequestParameters(nonce, ApiAccessData.ApiKey, signature))
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    var balance AccountBalance
    client := &http.Client{}
    resp, err := client.Do(req)

    defer resp.Body.Close()
    if err != nil {
        log.Printf("Could not send request: %s\n", err.Error())
        return nil
    } else {
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Could not read response data: %s\n", err.Error())
            return nil
        }
        json.Unmarshal(body, &balance)
    }
    return &balance
}

// GetCurrencyValue returns the characteristics of the given currency at the
// moment of retrieval and returns a CurrencySnapshot data struct or nil in
// case of an error.
func (bitstamp *Bitstamp) GetCurrencySnapshot(currency string) *CurrencySnapshot {
    log.Println("integrations::Bitstamp::GetCurrencyValue()")

    if ! helper.IsElementInArray(currency, bitstamp.supportedCurrencies) {
        log.Printf("This currency is not supported: %s \n", currency)
        os.Exit(1)
    }

    currency = strings.ToLower(currency) + "eur"

    var snapshot CurrencySnapshot
    resp, err := http.Get("https://www.bitstamp.net/api/v2/ticker/" + currency)

    defer resp.Body.Close()
    if err != nil {
        log.Printf("Could not send request: %s\n", err.Error())
        return nil
    } else {
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Could not read response data: %s\n", err.Error())
            return nil
        } else if resp.StatusCode > 299 {
            log.Printf("Failed request: Response code %d: %s\n", resp.StatusCode, string(body))
            return nil
        }
        json.Unmarshal(body, &snapshot)
    }
    return &snapshot
}

// GetOpenOrders retrieves all open orders for the given currency. If no open
// orders are available an empty array is returned. If something goes wrong the
// the error object is returned.
func (bitstamp *Bitstamp) GetOpenOrders(currency string) (error, []Order) {
    log.Println("integrations::Bitstamp::GetOpenOrders()")

    if ! helper.IsElementInArray(currency, bitstamp.supportedCurrencies) {
        log.Printf("This currency is not supported: %s \n", currency)
        os.Exit(1)
    }

    currency = strings.ToLower(currency) + "eur"
    nonce := int(time.Now().Unix())
    signature := createSignature(nonce, ApiAccessData)

    url := "https://www.bitstamp.net/api/v2/open_orders/" + currency + "/"

    postBody := []byte(createAuthRequestParameters(nonce, ApiAccessData.ApiKey, signature))
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    orders := make([]Order, 0)

    client := &http.Client{}
    resp, err := client.Do(req)

    defer resp.Body.Close()
    if err != nil {
        log.Printf("Could not send request: %s\n", err.Error())
        return err, nil
    } else {
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Could not read response data: %s\n", err.Error())
            return err, nil
        } else if resp.StatusCode > 299 {
            log.Printf("Failed request: Response code %d: %s\n", resp.StatusCode, string(body))
            return err, nil
        }
        json.Unmarshal(body, &orders)
    }
    return nil, orders
}

// GetSupportedCurrencies returns a list of all currencies that are supported
// by the respective trading platform integration. All currencies in that list
// are represented by their token and are kept in capital letters.
func (bitstamp *Bitstamp) GetSupportedCurrencies() []string {
    return bitstamp.supportedCurrencies
}

// Init sets the supported currencies for the Bitstamp integration and prints
// them to the console.
func (bitstamp *Bitstamp) Init() {
    log.Println("integrations::Bitstamp::Init()")

    bitstamp.supportedCurrencies = []string{
        "BTC", "XRP", "LTC", "ETH", "BCH" }

    log.Println("Supported currencies:", bitstamp.supportedCurrencies)
}

// createSignature takes a nonce, a customerId and an ApiAccess struct to
// create and return the signature string out of it that is necessary to send
// private requests to the Bitstamp API.
func createSignature(nonce int, apiCredentials ApiAccess) string {
    message := []byte(strconv.Itoa(nonce) + apiCredentials.CustomerId + apiCredentials.ApiKey)
    mac := hmac.New(sha256.New, []byte(apiCredentials.ApiSecret))
    mac.Write(message)
    return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}

// createAuthRequestParameters creates the string that is required in every
// private request to the Bitstamp API. All additional parameters can be easily
// appended with an '&' to the string which is returned from this function.
func createAuthRequestParameters(nonce int, apiKey string, signature string) string {
    return fmt.Sprintf(`key=%s&signature=%s&nonce=%d`, apiKey, signature, nonce)
}