package main

import (
	"github.com/GSabadini/go-transactions/infrastructure"
)

func main() {
	infrastructure.NewHTTPServer().Start()
}
