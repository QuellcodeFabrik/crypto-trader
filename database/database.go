package Database

import (
    _ "github.com/go-sql-driver/mysql"
	"database/sql"
    "log"
    "fmt"
)

type CryptoCurrency struct {
    id      int
    name    string
    token   string
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

func GetAvailableCurrencies() {
    stmtOut, err := db.Prepare("SELECT * FROM Currency WHERE id = ?")
    if err != nil {
        log.Println("Error on statement preparation.")
        panic(err.Error())
    }
    // close the statement after function is done
    defer stmtOut.Close()

    var currency CryptoCurrency // we "scan" the result in here

    err = stmtOut.QueryRow(1).Scan(&(currency.id), &(currency.name), &(currency.token)) // WHERE id = 1
    if err != nil {
        log.Println("Error on statement execution.")
        panic(err.Error())
    }

    fmt.Printf("The currency with ID 1 is: %+v \n", currency)
}