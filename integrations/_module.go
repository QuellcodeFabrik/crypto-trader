package integrations

type ExchangeIntegration interface {
    Init()
    GetCurrencyValue() CurrencySnapshot
}