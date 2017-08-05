package api

import (
	"github.com/baopham/gotime/github"
	"github.com/baopham/gotime/gotime"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

func getGithubRepoResponseTime(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	repo, service, err := getGithubRequestInfo(r, vars)

	if err != nil {
		handleError(err, w, service)
		return
	}

	responseTime, err := service.GetResponseTime(repo)

	if err != nil {
		handleError(err, w, service)
		return
	}

	respond(responseTime.Duration.String(), w)
}

func getGithubRepoLatestActivity(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	repo, service, err := getGithubRequestInfo(r, vars)

	if err != nil {
		handleError(err, w, service)
		return
	}

	activity, err := service.GetLatestActivity(repo)

	if err != nil {
		handleError(err, w, service)
		return
	}

	respond(activity, w)
}

func getGithubRequestInfo(r *http.Request, vars map[string]string) (*gotime.Repo, *github.Service, error) {
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

	return repo, service, err
}
