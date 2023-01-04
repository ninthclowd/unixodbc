package mariadb_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Context", func() {
	Describe("ExecContext", func() {
		It("should cancel if the context cancels", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
			defer cancel()
			start := time.Now()
			_, err := conn.ExecContext(ctx, "SELECT SLEEP(10)")
			elapsed := time.Since(start)
			Expect(elapsed.Seconds()).To(BeNumerically("<", 1))
			Expect(err).To(Equal(context.DeadlineExceeded))
		})
	})
	Describe("QueryContext", func() {
		It("should cancel if the context cancels", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
			defer cancel()
			start := time.Now()
			_, err := conn.QueryContext(ctx, "SELECT SLEEP(10)")
			elapsed := time.Since(start)
			Expect(elapsed.Seconds()).To(BeNumerically("<", 1))
			Expect(err).To(Equal(context.DeadlineExceeded))
		})
	})

})
