package github_test

import (
	"context"
	"fmt"
	. "github.com/baopham/gotime/github"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func respondWithFixture(path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile(path)
		fmt.Fprint(w, string(b))
	}
}

var _ = Describe("Github", func() {
	var (
		// server is a test HTTP server used to provide mock API responses.
		server *httptest.Server

		service *Service

		owner string

		repoName string
	)

	BeforeEach(func() {
		// test server
		mux := http.NewServeMux()
		server = httptest.NewServer(mux)
		url, _ := url.Parse(server.URL)
		req := &Request{
			HTTPClient: nil,
			Ctx:        context.Background(),
		}
		service = NewService(req)
		service.Client.BaseURL = url

		owner = "baopham"
		repoName = "gotime"

		mux.HandleFunc(fmt.Sprintf("/repos/%s/%s", owner, repoName), respondWithFixture("./fixtures/repo.json"))
		mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/issues", owner, repoName), respondWithFixture("./fixtures/issues.json"))
		mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/issues/1/comments", owner, repoName), respondWithFixture("./fixtures/comments_1.json"))
		mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/issues/2/comments", owner, repoName), respondWithFixture("./fixtures/comments_2.json"))
		mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/collaborators", owner, repoName), respondWithFixture("./fixtures/collaborators.json"))
	})

	AfterEach(func() {
		server.Close()
	})

	Context("when calling ResponseTime() on own repo", func() {
		It("should return the average response time", func() {
			repo, err := service.GetOwnRepo(owner, repoName)
			Expect(err).To(BeNil())

			responseTime, err := service.ResponseTime(repo)
			Expect(err).To(BeNil())

			Expect(responseTime.String()).To(Equal("1m58s"))
		})
	})

	Context("when calling ResponseTime() on other repo", func() {
		It("should return the average response time", func() {
			// TODO
		})
	})
})
