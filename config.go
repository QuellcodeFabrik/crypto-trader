package main

type ApiAccess struct {
    ApiKey string
    ApiSecret string
    CustomerId string
}

// This file is not tracked by Git anymore. To track it again do the following
// git update-index --no-assume-unchanged integrations/config.go

var ApiAccessData = ApiAccess{
    ApiKey: "",
    ApiSecret: "",
    CustomerId: "" }

var TestApiAccessData = ApiAccess{
    ApiKey: "apiKeyDEMO",
    ApiSecret: "apiSecretDEMO",
    CustomerId: "12345" }