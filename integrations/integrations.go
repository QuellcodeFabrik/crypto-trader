package integrations

type ExchangeIntegration interface {
    Init()
    GetAvailableCurrencies() []string
    GetCurrencyValue(token string) CurrencySnapshot
}

type CurrencySnapshot struct {
    Low         float64   `json:"low,string"`
    High        float64   `json:"high,string"`
    Current     float64   `json:"last,string"`
    Timestamp   int64     `json:"timestamp,string"`
    Average     float64   `json:"vwap,string"`
}