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
    "sync"
    "encoding/asn1"
)

var once sync.Once

type Bitstamp struct {
    supportedCurrencies []string
    apiAccessData ApiAccess
    lastNonce int
    timeSlot *TimeSlot
}

// CreateBuyOrder creates a buy order on the Bitstamp platform and gets it
// executed once the market prices match the given price value.
func (bitstamp *Bitstamp) CreateBuyOrder(currency string, amount float32, price float64) *Order {
    log.Println("integrations::Bitstamp::CreateBuyOrder()")

    if ! helper.IsElementInArray(currency, bitstamp.supportedCurrencies) {
        log.Printf("This currency is not supported: %s \n", currency)
        os.Exit(1)
    }

    currency = strings.ToLower(currency) + "eur"

    nonce := bitstamp.getNextNonce()
    signature := createSignature(nonce, bitstamp.apiAccessData)

    url := "https://www.bitstamp.net/api/v2/buy/" + currency + "/"

    postBody := []byte(createAuthRequestParameters(nonce, bitstamp.apiAccessData.ApiKey, signature) +
        "&amount=" + strconv.FormatFloat(float64(amount), 'f', -1, 32) +
        "&price=" + strconv.FormatFloat(price, 'f', -1, 64))

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    var order Order
    client := &http.Client{}

    // wait until free request slot is available
    for !bitstamp.HasFreeRequestSlot() {}

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
        } else if resp.StatusCode > 299 {
            log.Printf("Failed request: Response code %d: %s\n", resp.StatusCode, string(body))
            return nil
        }
        json.Unmarshal(body, &order)
    }
    return &order
}

// CreateSellOrder creates a sell order on the Bitstamp platform and gets it
// executed once the market prices match the given price value.
func (bitstamp *Bitstamp) CreateSellOrder(currency string, amount float32, price float64) *Order {
    log.Println("integrations::Bitstamp::CreateSellOrder()")

    if ! helper.IsElementInArray(currency, bitstamp.supportedCurrencies) {
        log.Printf("This currency is not supported: %s \n", currency)
        os.Exit(1)
    }

    currency = strings.ToLower(currency) + "eur"

    nonce := bitstamp.getNextNonce()
    signature := createSignature(nonce, bitstamp.apiAccessData)

    url := "https://www.bitstamp.net/api/v2/sell/" + currency + "/"

    postBody := []byte(createAuthRequestParameters(nonce, bitstamp.apiAccessData.ApiKey, signature) +
        "&amount=" + strconv.FormatFloat(float64(amount), 'f', -1, 32) +
        "&price=" + strconv.FormatFloat(price, 'f', -1, 64))

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    var order Order
    client := &http.Client{}

    // wait until free request slot is available
    for !bitstamp.HasFreeRequestSlot() {}

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
        } else if resp.StatusCode > 299 {
            log.Printf("Failed request: Response code %d: %s\n", resp.StatusCode, string(body))
            return nil
        }
        json.Unmarshal(body, &order)
    }
    return &order
}

// GetAccountBalance returns the overall account balance or nil in case
// something went wrong with the request.
func (bitstamp *Bitstamp) GetAccountBalance() *AccountBalance {
    log.Println("integrations::Bitstamp::GetAccountBalance()")

    nonce := bitstamp.getNextNonce()
    signature := createSignature(nonce, bitstamp.apiAccessData)

    url := "https://www.bitstamp.net/api/v2/balance/"

    postBody := []byte(createAuthRequestParameters(nonce, bitstamp.apiAccessData.ApiKey, signature))
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    var balance AccountBalance
    client := &http.Client{}

    // wait until free request slot is available
    for !bitstamp.HasFreeRequestSlot() {}

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
        } else if resp.StatusCode > 299 {
            log.Printf("Failed request: Response code %d: %s\n", resp.StatusCode, string(body))
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

    // wait until free request slot is available
    for !bitstamp.HasFreeRequestSlot() {}

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

    nonce := bitstamp.getNextNonce()
    signature := createSignature(nonce, bitstamp.apiAccessData)

    url := "https://www.bitstamp.net/api/v2/open_orders/" + currency + "/"

    postBody := []byte(createAuthRequestParameters(nonce, bitstamp.apiAccessData.ApiKey, signature))
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    orders := make([]Order, 0)
    client := &http.Client{}

    // wait until free request slot is available
    for !bitstamp.HasFreeRequestSlot() {}

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
    log.Println("integrations::Bitstamp::GetSupportedCurrencies()")

    return bitstamp.supportedCurrencies
}

// HasFreeRequestSlot uses the timeslot implementation to check if there is a
// free request timeslot. It has to wait as long as there is not free slot.
func (bitstamp *Bitstamp) HasFreeRequestSlot() bool {
    return bitstamp.timeSlot.IsFree()
}

// Init sets the supported currencies for the Bitstamp integration and prints
// them to the console. It additionally assigns the API access data that is
// required to send private request to the Bitstamp API.
func (bitstamp *Bitstamp) Init(apiAccessData ApiAccess) {
    log.Println("integrations::Bitstamp::Init()")

    bitstamp.timeSlot = &TimeSlot{}
    once.Do(func() {
        bitstamp.timeSlot.Init(1500)
    })

    bitstamp.apiAccessData = apiAccessData
    bitstamp.supportedCurrencies = []string{
        "BTC", "XRP", "LTC", "ETH", "BCH" }

    log.Println("Supported currencies:", bitstamp.supportedCurrencies)
}

// getNextNonce return the next nonce to be used in a request to the Bitstamp
// API. The method avoids to use the same nonce twice and thereby risking a
// failing request.
func (bitstamp *Bitstamp) getNextNonce() int {
    if bitstamp.lastNonce == 0 {
        bitstamp.lastNonce = int(time.Now().Unix())
    } else {
        bitstamp.lastNonce = bitstamp.lastNonce + 1
    }
    return bitstamp.lastNonce
}

// createAuthRequestParameters creates the string that is required in every
// private request to the Bitstamp API. All additional parameters can be easily
// appended with an '&' to the string which is returned from this function.
func createAuthRequestParameters(nonce int, apiKey string, signature string) string {
    return fmt.Sprintf(`key=%s&signature=%s&nonce=%d`, apiKey, signature, nonce)
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