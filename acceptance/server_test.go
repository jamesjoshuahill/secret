package acceptance_test

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	It("can print usage", func() {
		cmd := exec.Command(pathToServerBinary, "--help")
		session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(SatisfyAll(
			Say("--port= Port to serve HTTPS"),
			Say("--cert= Path to TLS certificate file"),
			Say("--key=  Path to TLS private key file"),
		))
	})

	It("fails when flags cannot be parsed", func() {
		cmd := exec.Command(pathToServerBinary, "--port=not-an-int")
		session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(SatisfyAll(
			Say("invalid argument for flag `--port'"),
		))
	})

	It("accepts a valid create cipher request", func() {
		res, err := client.Post(serverUrl("v1/ciphers"), "application/json", strings.NewReader(`{
			"id": "client-cipher-id",
			"data": "some plain text"
		}`))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(MatchJSON(`{
			"key": "key for client-cipher-id"
		}`))
	})

	It("rejects a malformed create cipher request", func() {
		res, err := client.Post(serverUrl("v1/ciphers"), "application/json", strings.NewReader("not json"))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(SatisfyAll(
			ContainSubstring("error"),
			ContainSubstring("decoding request body"),
		))
	})

	It("accepts a valid get cipher request", func() {
		req, err := http.NewRequest("GET", serverUrl("v1/ciphers/client-cipher-id"), strings.NewReader(`{
			"key": "key for client-cipher-id"
		}`))
		Expect(err).NotTo(HaveOccurred())

		res, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(MatchJSON(`{
			"data": "some plain text"
		}`))
	})

	It("rejects a malformed get cipher request", func() {
		req, err := http.NewRequest("GET", serverUrl("v1/ciphers/client-cipher-id"), strings.NewReader("not json"))
		Expect(err).NotTo(HaveOccurred())

		res, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(SatisfyAll(
			ContainSubstring("error"),
			ContainSubstring("decoding request body"),
		))
	})

	It("rejects a get cipher request with the wrong key", func() {
		req, err := http.NewRequest("GET", serverUrl("v1/ciphers/client-cipher-id"), strings.NewReader(`{
			"key": "wrong key"
		}`))
		Expect(err).NotTo(HaveOccurred())

		res, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusUnauthorized))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(SatisfyAll(
			ContainSubstring("error"),
			ContainSubstring("wrong key"),
		))
	})
})
