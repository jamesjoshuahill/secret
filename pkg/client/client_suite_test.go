package client_test

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

func readCloser(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}
