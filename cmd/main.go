package main

import (
	"github.com/cjdriod/go-credit-card-verifier/cmd/api"
	"github.com/cjdriod/go-credit-card-verifier/database"
)

func init() {
	database.ConnectDatabase()
	database.SyncDatabase()
}
func main() {

	app := api.InitServer()
	app.Serve()
}
