package acceptance_test

import (
	"encoding/json"
	"fmt"
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

	It("fails without the required flags", func() {
		cmd := exec.Command(pathToServerBinary)
		session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(SatisfyAll(
			Say("the required flags `--cert', `--key' and `--port' were not specified"),
		))
	})

	It("fails when a flag cannot be parsed", func() {
		cmd := exec.Command(pathToServerBinary, "--port=not-an-int")
		session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(SatisfyAll(
			Say("invalid argument for flag `--port'"),
		))
	})

	It("creates and stores secrets", func() {
		var key string
		By("accepting a valid create secret request", func() {
			res, err := httpsClient.Post(serverUrl("v1/ciphers"), "application/json", strings.NewReader(`{
				"id": "client-secret-id",
				"data": "some plain text"
			}`))
			Expect(err).NotTo(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

			var body createSecretResponseBody
			err = json.NewDecoder(res.Body).Decode(&body)
			Expect(err).NotTo(HaveOccurred())
			Expect(body.Key).NotTo(BeEmpty())
			key = body.Key
		})

		By("accepting a valid get secret request", func() {
			reqBody := fmt.Sprintf(`{
				"key": "%s"
			}`, key)
			req, err := http.NewRequest("GET", serverUrl("v1/ciphers/client-secret-id"), strings.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			res, err := httpsClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

			body, err := ioutil.ReadAll(res.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(MatchJSON(`{
				"data": "some plain text"
			}`))
		})
	})
})

type createSecretResponseBody struct {
	Key string `json:"key"`
}
