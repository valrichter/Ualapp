package main

import (
	"github.com/valrichter/Ualapp/api"
)

func main() {
	server := api.NewGinServer(".")
	server.Start(8080)
}
