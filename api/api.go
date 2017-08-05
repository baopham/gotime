package api

import (
	"encoding/json"
	"github.com/baopham/gotime/gotime"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func GetResponseTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")

	switch gotime.RepoProvider(vars["provider"]) {
	case gotime.GITHUB:
		getGithubRepoResponseTime(w, r, vars)
	}
}

func GetLatestActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")

	switch gotime.RepoProvider(vars["provider"]) {
	case gotime.GITHUB:
		getGithubRepoLatestActivity(w, r, vars)
	}
}

func respond(payload interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(payload)

	if err != nil {
		log.Println("failed to encode to JSON: ", err)
	}
}
