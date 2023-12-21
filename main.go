package main

import (
	"github.com/TiMattos/go-hokkaido/database"
	"github.com/TiMattos/go-hokkaido/routes"
)

func main() {
	database.ConectaComBancoDeDados()
	routes.HandleRequests()
}
