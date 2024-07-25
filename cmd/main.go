package main

import (
	"net/http"

	"tendercall-website.com/main/database"
	"tendercall-website.com/main/service/router"
)

func main() {
	database.Initdb()
	router.Route()

	http.ListenAndServe(":8080", nil)
}
