package integration_test

import (
	"context"
	. "github.com/baopham/gotime/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Github integration", func() {
	Context("when calling ResponseTime() on other repo", func() {
		It("should return the average response time", func() {
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
			Expect(responseTime).Should(BeNumerically(">", time.Duration(0)))

			log.Println(responseTime)
		})
	})
})
