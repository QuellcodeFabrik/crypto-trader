package database

import (
    "log"
    "time"
    "fmt"
)

type TradingPosition struct {
    Id         int64
    CurrencyId int64
    Amount     float64  `json:"amount,string"`
    Value      float64  `json:"value,string"`
    Timestamp  int64    `json:"timestamp,string"`
}

// GetTradingPositions returns all positions for the given
// currency or all positions if no currency is given.
func GetTradingPositions(currency *CryptoCurrency) []TradingPosition {
    log.Println("database::GetTradingPositions()")

    queryString := "SELECT * FROM Position"

    if currency != nil {
        queryString += " WHERE currencyId=" + string(currency.Id)
    }

    stmtOut, err := db.Prepare(queryString)
    if err != nil {
        log.Println("Error on statement preparation.")
        panic(err.Error())
    }
    defer stmtOut.Close()

    var tradingPositions []TradingPosition

    rows, err := stmtOut.Query()
    if err != nil {
        log.Println("Error on statement execution.")
        panic(err.Error())
    }

    defer rows.Close()
    for rows.Next() {
        var tradingPosition TradingPosition

        if err := rows.Scan(
            &(tradingPosition.Id), &(tradingPosition.Amount), &(tradingPosition.CurrencyId),
            &(tradingPosition.Timestamp), &(tradingPosition.Value)); err != nil {
            log.Fatal(err)
        }

        tradingPositions = append(tradingPositions, tradingPosition)
    }

    if rows.Err() != nil {
        log.Println("Error on iterating over rows.")
        panic(err.Error())
    }

    return tradingPositions
}

// AddTradingPosition adds a new position item into the database as soon as the
// trade was executed successfully.
func AddTradingPosition(currency *CryptoCurrency, position *TradingPosition) error {
    log.Println("database::AddTradingPosition()")

    insert, err := db.Prepare("INSERT INTO Position (timestamp, currencyId, amount, value) " +
        "VALUES( ?, ?, ?, ? )")

    if err != nil {
        return err
    }
    defer insert.Close()

    timestamp := time.Unix(position.Timestamp, 0)
    _, err = insert.Exec(timestamp, currency.Id, position.Amount, position.Value)
    return err
}

// RemoveTradingPosition removes the position item from the database as soon as
// it was sold successfully.
func RemoveTradingPosition(positionId int64) error {
    deleteStatement, err := db.Prepare(`DELETE FROM Position WHERE id = ?`)

    if err != nil {
        return err
    }
    defer deleteStatement.Close()

    result, err := deleteStatement.Exec(positionId)
    rowsAffected, err := result.RowsAffected()

    if err != nil {
        return fmt.Errorf("mysql: could not get rows affected: %v", err)
    } else if rowsAffected != 1 {
        return fmt.Errorf("mysql: expected 1 row affected, got %d", rowsAffected)
    }
    return err
}