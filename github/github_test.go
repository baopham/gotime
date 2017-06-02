package github_test

import (
	"context"
	"github.com/baopham/gotime"
	. "github.com/baopham/gotime/github"
	"golang.org/x/oauth2"
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Github", func() {
	Context("when calling ResponseTime()", func() {
		It("should return list of issues", func() {
			// TODO: update test with mocks
			req := &Request{
				HTTPClient: oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
				)),
				Ctx: context.Background(),
			}
			repo := &gotime.Repo{
				Owner:    "baopham",
				Name:     "gotime",
				Members:  []string{"baopham"},
				Provider: gotime.GITHUB,
			}
			responseTime, err := ResponseTime(req, repo)
			Expect(err).To(BeNil())
			log.Println(responseTime)
		})
	})
})
