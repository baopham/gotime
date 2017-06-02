package gotime

import (
	"time"
)

type RepoProvider uint8

const (
	GITHUB RepoProvider = iota
)

type Repo struct {
	Owner, Name, URL string
	Members          []string
	Provider         RepoProvider
}

type IssueInfo struct {
	Repo      *Repo
	Number    int
	CreatedAt *time.Time
	ClosedAt  *time.Time
	// Owner, repo's member response times
	EarliestResponse *time.Time
}
