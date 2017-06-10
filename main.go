package main

import (
	"github.com/baopham/gotime/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	r.Path("/response-time/{provider}/{owner}/{repo}").
		HandlerFunc(api.GetResponseTime)

	port := "3000"

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Fatalln(http.ListenAndServe(":"+port, r))
}
