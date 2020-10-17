package main

import (
	"os"

	"github.com/GSabadini/go-transactions/infrastructure"
)

func main() {
	infrastructure.NewApplication().Start(os.Getenv("APP_PORT"))
}
