package github_test

import (
	"context"
	. "github.com/baopham/gotime/github"
	"golang.org/x/oauth2"
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO: update test with mocks
var _ = Describe("Github", func() {
	Context("when calling ResponseTime()", func() {
		It("should return list of issues", func() {
			req := &Request{
				HTTPClient: oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
				)),
				Ctx: context.Background(),
			}
			service := NewService(req)

			repo, err := service.GetOtherRepo("go-kit", "kit")
			Expect(err).To(BeNil())

			responseTime, err := service.ResponseTime(repo)
			Expect(err).To(BeNil())

			log.Println(responseTime)
		})
	})

	Context("when calling GetOwnRepo", func() {
		It("should return the repo info", func() {
			req := &Request{
				HTTPClient: oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
				)),
				Ctx: context.Background(),
			}
			service := NewService(req)
			repo, err := service.GetOwnRepo("baopham", "gotime")
			Expect(err).To(BeNil())
			log.Println(*repo.Members)
		})
	})
})
