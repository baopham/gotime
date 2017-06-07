package api

import (
	"encoding/json"
	"github.com/baopham/gotime/github"
	"github.com/baopham/gotime/gotime"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

func GetResponseTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	switch gotime.RepoProvider(vars["provider"]) {
	case gotime.GITHUB:
		getGithubRepoResponseTime(&w, r, vars)
	}
}

func getGithubRepoResponseTime(w *http.ResponseWriter, r *http.Request, vars map[string]string) {
	var token oauth2.TokenSource
	owner, repoName := vars["owner"], vars["repo"]
	context := r.Context()
	req := &github.Request{
		Ctx: context,
	}

	t := strings.TrimSpace(r.URL.Query().Get("token"))

	if t != "" {
		token = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: t})
		req.HTTPClient = oauth2.NewClient(context, token)
	}

	service := github.NewService(req)

	var repo *gotime.Repo
	var err error

	if req.HTTPClient != nil {
		repo, err = service.GetOwnRepo(owner, repoName)
	} else {
		repo, err = service.GetOtherRepo(owner, repoName)
	}

	if err != nil {
		handleError(err, w, service)
		return
	}

	responseTime, err := service.GetResponseTime(repo)

	if err != nil {
		handleError(err, w, service)
		return
	}

	json.NewEncoder(*w).Encode(responseTime.Duration.String())
}

func handleError(err error, w *http.ResponseWriter, service gotime.GoTimer) {
	message := "Something went wrong"

	if service.IsRateLimitError(err) {
		message = "API rate limit. Supply a token to increase the limit"
	}

	http.Error(*w, message, http.StatusBadRequest)
}
