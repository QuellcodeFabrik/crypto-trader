package database

import (
    "log"
    "fmt"
    "time"
    "../integrations"
)

// StoreSnapshot stores the given snapshot struct in the database taking into
// account the currency reference under which it is known in the connected
// currency exchange platform.
func StoreSnapshot(currencyToken string, snapshot *integrations.CurrencySnapshot) error {
    log.Println("database::StoreSnapshot()")

    insert, err := db.Prepare("INSERT INTO Snapshot (timestamp, currency, value, low, high, average) " +
        "VALUES( ?, ?, ?, ?, ?, ? )")

    if err != nil {
        return err
    }
    defer insert.Close()

    currency := getCryptoCurrencyByToken(currencyToken, availableCurrencies)

    if currency == nil {
        return &Error{
            fmt.Sprintf("Crypto-Currency '%s' is unkown. Cannot insert database.", currencyToken)}
    }

    timestamp := time.Unix(snapshot.Timestamp, 0)
    _, err = insert.Exec(timestamp, currency.Id, snapshot.Current, snapshot.Low, snapshot.High, snapshot.Average)
    return err
}