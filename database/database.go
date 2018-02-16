package database

import (
    _ "github.com/go-sql-driver/mysql"
	"database/sql"
    "log"
)

type CryptoCurrency struct {
    Id      int     `json:"id"`
    Name    string  `json:"name"`
    Token   string  `json:"token"`
}

// package global database handle
var db *sql.DB = nil

func Deinit() {
    db.Close()
    db = nil
}

func Init() {
    if db == nil {
        var err error = nil
        db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/crypto_trader")

        if err != nil {
            log.Println("Could not open database");
            panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
        }

        // Validate DSN data:
        err = db.Ping()
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }
    }
}

func GetAvailableCurrency() CryptoCurrency {
    log.Println("Database::GetAvailableCurrencies()")

    stmtOut, err := db.Prepare("SELECT * FROM Currency WHERE id = ?")
    if err != nil {
        log.Println("Error on statement preparation.")
        panic(err.Error())
    }
    // close the statement after function is done
    defer stmtOut.Close()

    var currency CryptoCurrency // we "scan" the result in here

    err = stmtOut.QueryRow(1).Scan(&(currency.Id), &(currency.Name), &(currency.Token)) // WHERE id = 1
    if err != nil {
        log.Println("Error on statement execution.")
        panic(err.Error())
    }

    return currency
}

func GetAvailableCurrencies() []CryptoCurrency {
    log.Println("Database::GetAvailableCurrencies()")

    stmtOut, err := db.Prepare("SELECT * FROM Currency")
    if err != nil {
        log.Println("Error on statement preparation.")
        panic(err.Error())
    }
    // close the statement after function is done
    defer stmtOut.Close()

    var currencies []CryptoCurrency

    rows, err := stmtOut.Query()
    if err != nil {
        log.Println("Error on statement execution.")
        panic(err.Error())
    }

    defer rows.Close()
    for rows.Next() {
        var currency CryptoCurrency

        if err := rows.Scan(&(currency.Id), &(currency.Name), &(currency.Token)); err != nil {
            log.Fatal(err)
        }

        currencies = append(currencies, currency)
    }

    if rows.Err() != nil {
        log.Println("Error on iterating over rows.")
        panic(err.Error())
    }

    return currencies
}