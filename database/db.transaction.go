package database

import (
    "log"
    "fmt"
    "time"
)

type Transaction struct {
    Id              int
    CurrencyId      int
    Amount          int
    Value           float64
    MarketTendency  string
    Timestamp       time.Time
    Type            string
}

// GetTransactions returns all transactions for the given
// currency or all transactions if no currency is given.
func GetTransactions(currency *CryptoCurrency) ([]Transaction, error) {
    log.Println("database::GetTransactions()")

    queryString := "SELECT * FROM Transaction"

    if currency != nil {
        queryString += " WHERE currencyId=" + string(currency.Id)
    }

    stmtOut, err := db.Prepare(queryString)
    if err != nil {
        log.Println("Error on preparing database statement.")
        return nil, err
    }
    defer stmtOut.Close()

    var transactions []Transaction

    rows, err := stmtOut.Query()
    if err != nil {
        log.Println("Error on statement execution.")
        return nil, err
    }

    defer rows.Close()
    for rows.Next() {
        var transaction Transaction

        if err := rows.Scan(
            &(transaction.Id), &(transaction.Amount), &(transaction.CurrencyId),
            &(transaction.MarketTendency), &(transaction.Timestamp), &(transaction.Type),
            &(transaction.Value)); err != nil {
            return nil, err
        }

        transactions = append(transactions, transaction)
    }

    if rows.Err() != nil {
        log.Println("Error on iterating over rows.")
        return nil, err
    }

    return transactions, nil
}

// StoreTransaction stores the given transaction struct in the database taking into
// account the currency reference under which it is known in the connected
// currency exchange platform.
func StoreTransaction(currencyToken string, transaction *Transaction) error {
    log.Println("database::StoreTransaction()")

    insert, err := db.Prepare("INSERT INTO Transaction (timestamp, currencyId, " +
        "marketTendency, type, amount, value) VALUES( ?, ?, ?, ?, ?, ? )")

    if err != nil {
        return err
    }
    defer insert.Close()

    currency := getCryptoCurrencyByToken(currencyToken, availableCurrencies)

    if currency == nil {
        return &Error{
            fmt.Sprintf("Crypto-Currency '%s' is unkown. " +
                "Cannot insert transaction into database.", currencyToken)}
    }

    timestamp := transaction.Timestamp.Unix()
    _, err = insert.Exec(timestamp, currency.Id, transaction.MarketTendency,
        transaction.Type, transaction.Amount, transaction.Value)
    return err
}