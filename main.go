package main

import (
	"crud/lib/api"
	"crud/lib/database"
)

func main() {
	database.DB()
	api.Routes()
}
