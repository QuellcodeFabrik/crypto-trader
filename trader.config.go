package main

import (
    "./integrations"
)

// This file is not tracked by Git anymore. To track it again do the following
// git update-index --no-assume-unchanged trader.config.go

var ApiAccessData = integrations.ApiAccess{
    ApiKey: "",
    ApiSecret: "",
    CustomerId: "" }

var TestApiAccessData = integrations.ApiAccess{
    ApiKey: "apiKeyDEMO",
    ApiSecret: "apiSecretDEMO",
    CustomerId: "12345" }