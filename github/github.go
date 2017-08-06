package github

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/baopham/gotime/concurrentslice"
	"github.com/baopham/gotime/gotime"
	"github.com/google/go-github/github"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Request struct {
	HTTPClient *http.Client
	Ctx        context.Context
}

type Service struct {
	Client *github.Client
	Ctx    context.Context
}

func NewService(req *Request) *Service {
	return &Service{
		Client: github.NewClient(req.HTTPClient),
		Ctx:    req.Ctx,
	}
}

// GetResponseTime gives the general time that is needed to respond to an issue
func (s *Service) GetResponseTime(repo *gotime.Repo) (*gotime.ResponseTime, error) {
	issues, err := s.getIssues(repo)

	if err != nil {
		return nil, err
	}

	infos := make(chan *gotime.IssueInfo, len(issues))
	duration := make(chan time.Duration, 1)

	go s.collectIssues(repo, issues, infos)
	go processIssues(infos, duration)

	return &gotime.ResponseTime{<-duration}, nil
}

// GetLatestActivity gives the latest activity in the repo (be it commit or response to an issue)
func (s *Service) GetLatestActivity(repo *gotime.Repo) (*gotime.Activity, error) {
	var wg sync.WaitGroup
	wg.Add(3)

	activities := make(chan *gotime.Activity)
	listOptions := github.ListOptions{PerPage: 10}

	// Commit type - get the latest one
	go func() {
		defer wg.Done()
		commits, _, err := s.Client.Repositories.ListCommits(
			s.Ctx,
			repo.Owner,
			repo.Name,
			&github.CommitsListOptions{ListOptions: listOptions},
		)
		if err == nil && len(commits) > 0 {
			commit := commits[0]
			activities <- &gotime.Activity{
				Type: "Commit",
				URL:  commit.HTMLURL,
				Time: commit.Commit.Committer.Date,
			}
		}
	}()

	// Search for latest comment activity of the owner
	go func() {
		defer wg.Done()
		result, _, err := s.Client.Search.Issues(
			s.Ctx,
			fmt.Sprintf("commenter:%s repo:%s/%s", repo.Owner, repo.Owner, repo.Name),
			&github.SearchOptions{ListOptions: listOptions},
		)
		if err == nil && result != nil && len(result.Issues) > 0 {
			issue := result.Issues[0]
			activities <- &gotime.Activity{
				Type: "Comment",
				URL:  issue.HTMLURL,
				// Not entirely correct that we use updated_at.
				// TODO: improve this
				Time: issue.UpdatedAt,
			}
		}
	}()

	// Search for latest issue activity of the owner
	go func() {
		defer wg.Done()
		result, _, err := s.Client.Search.Issues(
			s.Ctx,
			fmt.Sprintf("author:%s repo:%s/%s", repo.Owner, repo.Owner, repo.Name),
			&github.SearchOptions{Sort: "created", Order: "desc", ListOptions: listOptions},
		)
		if err == nil && result != nil && len(result.Issues) > 0 {
			issue := result.Issues[0]
			activities <- &gotime.Activity{
				Type: "Issue",
				URL:  issue.HTMLURL,
				Time: issue.CreatedAt,
			}
		}
	}()

	go func() {
		wg.Wait()
		close(activities)
	}()

	var latestActivity *gotime.Activity

	for activity := range activities {
		if latestActivity == nil || latestActivity.Time.Before(*activity.Time) {
			latestActivity = activity
		}
	}

	return latestActivity, nil
}

func (s *Service) GetOwnRepo(owner, repoName string) (*gotime.Repo, error) {
	repo, _, err := s.Client.Repositories.Get(s.Ctx, owner, repoName)

	if err != nil {
		return nil, err
	}

	members := []string{}
	page := 1

	for page > 0 {
		collabs, resp, err := s.Client.Repositories.ListCollaborators(
			s.Ctx,
			owner,
			repoName,
			&github.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		)

		if err != nil {
			log.Printf("!!!! failed to get collabs for %s, err: %s", repoName, err)
			return nil, err
		}

		for _, u := range collabs {
			members = append(members, u.GetLogin())
		}
		page = resp.NextPage
	}
	return &gotime.Repo{
		Owner:    owner,
		Name:     repoName,
		URL:      repo.GetURL(),
		Members:  &members,
		Provider: gotime.GITHUB,
	}, nil
}

func (s *Service) GetOtherRepo(owner, repoName string) (*gotime.Repo, error) {
	repo, _, err := s.Client.Repositories.Get(s.Ctx, owner, repoName)

	if err != nil {
		log.Printf("!!!! failed to get repo %s, err: %s", repoName, err)
		return nil, err
	}

	return &gotime.Repo{
		Owner:    owner,
		Name:     repoName,
		URL:      repo.GetURL(),
		Members:  nil,
		Provider: gotime.GITHUB,
	}, nil
}

func (s *Service) IsRateLimitError(err error) bool {
	_, ok := err.(*github.RateLimitError)
	return ok
}

// Get some latest sample issues
func (s *Service) getIssues(repo *gotime.Repo) ([]*github.Issue, error) {
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

func (s *Service) collectIssues(repo *gotime.Repo, issues []*github.Issue, c chan<- *gotime.IssueInfo) {
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

func (s *Service) getIssueInfo(repo *gotime.Repo, issue *github.Issue, page int) (*gotime.IssueInfo, error) {
	info := &gotime.IssueInfo{
		Repo:      repo,
		Number:    *issue.Number,
		CreatedAt: issue.CreatedAt,
	}

	if issue.ClosedBy != nil && issue.ClosedBy.GetLogin() != issue.User.GetLogin() {
		info.OtherClosedAt = issue.ClosedAt
	}

	comments, resp, err := s.Client.Issues.ListComments(
		s.Ctx,
		repo.Owner,
		repo.Name,
		issue.GetNumber(),
		&github.IssueListCommentsOptions{
			Sort:      "created",
			Direction: "asc",
			ListOptions: github.ListOptions{
				Page: page,
			},
		},
	)

	if err != nil {
		log.Printf("!!!! failed to get comments for %s, err: %s", *issue.Title, err)
		return nil, err
	}

	if repo.Members != nil {
		members := stringsToMap(*repo.Members)

		// Find the earliest comment made by one of the repo's members
		for _, comment := range comments {
			if _, isMember := members[comment.User.GetLogin()]; isMember {
				info.EarliestResponse = comment.CreatedAt
				break
			}
		}
	} else {
		doc, err := goquery.NewDocument(issue.GetHTMLURL())
		if err != nil {
			log.Printf("!!!! failed to parse HTML, err: %s", err)
			return nil, err
		}
		for _, comment := range comments {
			if valid, _ := commentMadeByMember(doc, comment, repo); valid {
				info.EarliestResponse = comment.CreatedAt
				break
			}
		}
	}

	if info.EarliestResponse == nil && resp.LastPage > page {
		return s.getIssueInfo(repo, issue, resp.NextPage)
	}

	return info, nil
}

func commentMadeByMember(doc *goquery.Document, c *github.IssueComment, repo *gotime.Repo) (bool, error) {
	if c.User.GetLogin() == repo.Owner {
		return true, nil
	}

	comment := doc.Find(fmt.Sprintf("#issuecomment-%d", c.GetID()))
	label := strings.TrimSpace(comment.Find(".timeline-comment-label").Text())

	return label == "Owner" || label == "Member" || label == "Contributor", nil
}

func processIssues(infos chan *gotime.IssueInfo, duration chan time.Duration) {
	slice := concurrentslice.New()

	for {
		info, more := <-infos
		if !more {
			duration <- getAverageResponseTime(slice)
			return
		}

		var earliestTime *time.Time

		if info.EarliestResponse == nil {
			earliestTime = info.OtherClosedAt
		} else if info.OtherClosedAt != nil && info.OtherClosedAt.Before(*info.EarliestResponse) {
			earliestTime = info.OtherClosedAt
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
