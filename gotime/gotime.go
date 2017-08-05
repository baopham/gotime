package gotime

import (
	"time"
)

type RepoProvider string
type Responsiveness uint8

const (
	GITHUB RepoProvider = "github"
)

const (
	VERY_RESPONSIVE Responsiveness = iota
	RESPONSIVE
	NOT_RESPONSIVE
)

type GoTimer interface {
	GetResponseTime(repo *Repo) (*ResponseTime, error)
	GetOwnRepo(owner, repoName string) (*Repo, error)
	GetOtherRepo(owner, repoName string) (*Repo, error)
	IsRateLimitError(err error) bool
	GetLatestActivity(repo *Repo) (*LatestActivity, error)
}

type Repo struct {
	Owner, Name, URL string
	Members          *[]string
	Provider         RepoProvider
}

type IssueInfo struct {
	Repo      *Repo
	Number    int
	CreatedAt *time.Time
	// Not nil only if it's closed by someone else other than the author
	OtherClosedAt *time.Time
	// Owner, repo's member response times
	EarliestResponse *time.Time
}

type ResponseTime struct {
	time.Duration
}

type LatestActivity struct {
	Time time.Time
	Type string
	URL  string
}

func (r ResponseTime) GetResponsiveness() Responsiveness {
	if d, err := time.ParseDuration("48h"); err == nil && r.Duration <= d {
		return VERY_RESPONSIVE
	}

	if d, err := time.ParseDuration("96h"); err == nil && r.Duration <= d {
		return RESPONSIVE
	}

	return NOT_RESPONSIVE
}
