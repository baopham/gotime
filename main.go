package main

import (
	"github.com/baopham/gotime/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.Path("/response-time/{provider}/{owner}/{repo}").
		HandlerFunc(api.GetResponseTime)

	log.Fatalln(http.ListenAndServe(":8000", r))
}
