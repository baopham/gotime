package gotime_test

import (
	. "github.com/baopham/gotime/gotime"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gotime", func() {
	Context("ResponseTime", func() {
		Context("when calling GetResponsiveness()", func() {
			It("should return VERY_RESPONSIVE if the response time <= 48h", func() {
				d, err := time.ParseDuration("32h")
				Expect(err).To(BeNil())
				responseTime := ResponseTime{d}
				Expect(responseTime.GetResponsiveness()).To(Equal(VERY_RESPONSIVE))
			})

			It("should return RESPONSIVE if the response time is between 48h and 96h", func() {
				d, err := time.ParseDuration("50h")
				Expect(err).To(BeNil())
				responseTime := ResponseTime{d}
				Expect(responseTime.GetResponsiveness()).To(Equal(RESPONSIVE))
			})

			It("should return NOT_RESPONSIVE if the response time > 96h", func() {
				d, err := time.ParseDuration("97h")
				Expect(err).To(BeNil())
				responseTime := ResponseTime{d}
				Expect(responseTime.GetResponsiveness()).To(Equal(NOT_RESPONSIVE))
			})
		})
	})
})
