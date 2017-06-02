package github

import (
	"context"
	"github.com/baopham/gotime"
	"github.com/baopham/gotime/concurrentslice"
	"github.com/google/go-github/github"
	"log"
	"net/http"
	"sync"
	"time"
)

type Request struct {
	HTTPClient *http.Client
	Ctx        context.Context
}

type service struct {
	Client *github.Client
	Ctx    context.Context
}

// ResponseTime gives the general time that is needed to respond to an issue
func ResponseTime(req *Request, repo *gotime.Repo) (time.Duration, error) {
	s := service{
		Client: github.NewClient(req.HTTPClient),
		Ctx:    req.Ctx,
	}

	issues, err := s.getIssues(repo)

	if err != nil {
		return 0, err
	}

	infos := make(chan *gotime.IssueInfo, len(issues))
	duration := make(chan time.Duration, 1)

	go s.collect(repo, issues, infos)
	go process(infos, duration)

	return <-duration, nil
}

// Get some latest sample issues
func (s *service) getIssues(repo *gotime.Repo) ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		Sort:      "created",
		Direction: "desc",
		State:     "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	issues, _, err := s.Client.Issues.ListByRepo(s.Ctx, repo.Owner, repo.Name, opt)
	return issues, err
}

func (s *service) collect(repo *gotime.Repo, issues []*github.Issue, c chan<- *gotime.IssueInfo) {
	var wg sync.WaitGroup
	wg.Add(len(issues))

	for _, issue := range issues {
		go func(issue *github.Issue) {
			defer wg.Done()
			info, err := s.getIssueInfo(repo, issue, 1)
			if err != nil {
				return
			}
			c <- info
		}(issue)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
}

func (s *service) getIssueInfo(repo *gotime.Repo, issue *github.Issue, page int) (*gotime.IssueInfo, error) {
	info := &gotime.IssueInfo{
		Repo:      repo,
		Number:    *issue.Number,
		ClosedAt:  issue.ClosedAt,
		CreatedAt: issue.CreatedAt,
	}

	opt := &github.IssueListCommentsOptions{
		Sort:      "created",
		Direction: "asc",
		ListOptions: github.ListOptions{
			Page: page,
		},
	}

	comments, resp, err := s.Client.Issues.ListComments(
		s.Ctx,
		repo.Owner,
		repo.Name,
		*issue.Number,
		opt,
	)

	if err != nil {
		log.Printf("!!!! failed to get comments for %s, err: %s", *issue.Title, err)
		return nil, err
	}

	members := stringsToMap(repo.Members)

	// Find the earliest comment made by one of the repo's members
	for _, comment := range comments {
		if _, isMember := members[*comment.User.Login]; isMember {
			info.EarliestResponse = comment.CreatedAt
			break
		}
	}

	if info.EarliestResponse == nil && resp.LastPage > page {
		return s.getIssueInfo(repo, issue, resp.NextPage)
	}

	return info, nil
}

func process(infos chan *gotime.IssueInfo, duration chan time.Duration) {
	slice := concurrentslice.New()

	for {
		info, more := <-infos
		if !more {
			duration <- getAverageResponseTime(slice)
			return
		}

		var earliestTime *time.Time

		if info.EarliestResponse == nil {
			earliestTime = info.ClosedAt
		} else if info.ClosedAt != nil && info.ClosedAt.Before(*info.EarliestResponse) {
			earliestTime = info.ClosedAt
		} else {
			earliestTime = info.EarliestResponse
		}

		if earliestTime != nil {
			slice.Append(earliestTime.Sub(*info.CreatedAt))
		}
	}
}

func getAverageResponseTime(slice *concurrentslice.Slice) time.Duration {
	log.Printf("getting avg response time for %d items\n", slice.Size())

	if slice.Size() == 0 {
		return 0
	}

	times := make([]time.Duration, slice.Size())
	slice.Fill(func(i int, v interface{}) {
		times[i] = v.(time.Duration)
	})

	sum := times[0]
	for _, t := range times[1:] {
		sum += t
	}

	return sum / time.Duration(len(times))
}

func stringsToMap(arr []string) map[string]bool {
	m := make(map[string]bool)
	for _, value := range arr {
		m[value] = true
	}
	return m
}
