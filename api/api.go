package api

import (
	"github.com/baopham/gotime/gotime"
	"github.com/gorilla/mux"
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
