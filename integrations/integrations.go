package integrations

import (
    "sync"
    "time"
)

type ExchangeIntegration interface {
    CreateBuyOrder(currency string, amount float32, price float64) *Order
    CreateSellOrder(currency string, amount float32, price float64) *Order
    GetAccountBalance() *AccountBalance
    GetCurrencySnapshot(currency string) *CurrencySnapshot
    GetOpenOrders(currency string) (error, []Order)
    GetSupportedCurrencies() []string
    HasFreeRequestSlot() bool
    Init()
}

type TimeSlot struct {
    mutex     sync.Mutex
    interval  int64
    last      int64
}

func (s *TimeSlot) Init(interval int64) {
    s.interval = interval
}

func (s *TimeSlot) IsFree() bool {
    s.mutex.Lock()
    if s.last <= time.Now().Unix() - s.interval {
        s.last = time.Now().Unix()
	s.mutex.Unlock()
        return true
    }
    s.mutex.Unlock()
    return false;
}

type AccountBalance struct {
    BCH     float64 `json:"bch_balance,string"`
    BTC     float64 `json:"btc_balance,string"`
    ETH     float64 `json:"eth_balance,string"`
    EUR     float64 `json:"eur_balance,string"`
    LTC     float64 `json:"ltc_balance,string"`
    XRP     float64 `json:"xrp_balance,string"`
}

type ApiAccess struct {
    ApiKey string
    ApiSecret string
    CustomerId string
}

type CurrencySnapshot struct {
    Id          int64
    Low         float64   `json:"low,string"`
    High        float64   `json:"high,string"`
    Current     float64   `json:"last,string"`
    Timestamp   int64     `json:"timestamp,string"`
    Average     float64   `json:"vwap,string"`
}

type Order struct {
    Id          int64   `json:"id,string"`
    Timestamp   int64   `json:"datetime,string"`
    Type        int     `json:"type,string"`
    Price	    float64 `json:"price,string"`
    Amount      float64 `json:"amount,string"`
}

func (order *Order) IsBuyingOrder() bool {
    return order.Type == 0
}

func (order *Order) IsSellingOrder() bool {
    return order.Type == 1
}
