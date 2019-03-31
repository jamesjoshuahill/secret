package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	It("listens on $PORT", func() {
		Eventually(dialServer).Should(Succeed())
	})
})
