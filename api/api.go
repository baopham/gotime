package api

import (
	"github.com/baopham/gotime/gotime"
	"github.com/gorilla/mux"
	"net/http"
)

func GetResponseTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	switch gotime.RepoProvider(vars["provider"]) {
	case gotime.GITHUB:
		getGithubRepoResponseTime(&w, r, vars)
	}
}
