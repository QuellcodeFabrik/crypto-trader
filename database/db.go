// The database package is a collection of database access methods making the
// storage and retrieval of data convenient and controlled.
package database

import (
    _ "github.com/go-sql-driver/mysql"
	"database/sql"
    "log"
    "fmt"
    "time"
)

type CryptoCurrency struct {
    Id        int     `json:"id"`
    Name      string  `json:"name"`
    Token     string  `json:"token"`
}

type InvestmentPosition struct {
    Id        int
    Currency  CryptoCurrency
    Amount    int
    Value     float64
    Timestamp time.Time
}

type Error struct {
    message string
}

func (e *Error) Error() string {
    return fmt.Sprintf("Error in database handler: %s", e.message)
}

var db *sql.DB = nil
var availableCurrencies []CryptoCurrency

func DeInit() {
    db.Close()
    db = nil
}

// Init initializes the database connection and does an initial test request to
// check if the connection is working as expected.
func Init() {
    log.Println("database::Init()")

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

        availableCurrencies = getAvailableCurrencies()
    }
}

//
// Private functions
//

func getAvailableCurrencies() []CryptoCurrency {
    log.Println("database::GetAvailableCurrencies()")

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

func getCryptoCurrencyByToken(token string, list []CryptoCurrency) *CryptoCurrency {
    for _, listItem := range list {
        if listItem.Token == token {
            return &listItem
        }
    }
    return nil
}