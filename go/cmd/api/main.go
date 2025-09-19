package main

import (
	"goapp/internal/web"
	"log"
	"net/http"
)


func main() {
	r := web.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}

