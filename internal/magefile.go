// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = All

// Run the default targets
func All() {
	mg.SerialDeps(Lint, Test)
}

// Build Build the server binary
func Build() error {
	logStep("Build")
	return sh.Run("go", "build", "-o", "./bin/secret-server", "./cmd/secret-server")
}

// Run the linters
func Lint() error {
	mg.Deps(InstallGolangCILint)
	logStep("Lint")
	return sh.RunV("./tools/golangci-lint", "run")
}

// Run all tests
func Test() {
	mg.SerialDeps(TestUnit, TestAcceptance)
}

// Run unit tests
func TestUnit() error {
	mg.Deps(InstallGinkgo)
	logStep("Test Unit")
	return sh.RunV("./tools/ginkgo", "-race", "-cover", "-r", "-skipPackage", "acceptance_test")
}

// Run acceptance tests
func TestAcceptance() error {
	mg.Deps(Build, InstallGinkgo)
	logStep("Test Acceptance")
	return sh.RunV("./tools/ginkgo", "acceptance_test")
}

// Install golangci-lint
func InstallGolangCILint() error {
	needUpdate, err := target.Path("./tools/golangci-lint")
	if err != nil || !needUpdate {
		return err
	}

	logStep("Install golangci-lint")

	version := "1.38.0"

	err = os.MkdirAll("./tools", 0755)
	if err != nil {
		return err
	}

	tarball, err := ioutil.TempFile("", "golangci-lint")
	if err != nil {
		return err
	}
	defer sh.Rm(tarball.Name())

	asset := fmt.Sprintf("golangci-lint-%s-%s-amd64", version, runtime.GOOS)
	u := fmt.Sprintf("https://github.com/golangci/golangci-lint/releases/download/v%s/%s.tar.gz",
		version, asset)

	err = sh.Run("curl", "-sSfL", u, "-o", tarball.Name())
	return sh.Run("tar", "xf", tarball.Name(),
		"-C", "tools",
		"--strip-components", "1",
		fmt.Sprintf("%s/golangci-lint", asset))
}

// Install ginkgo
func InstallGinkgo() error {
	needUpdate, err := target.Path("./tools/ginkgo", "./go.mod")
	if err != nil || !needUpdate {
		return err
	}

	logStep("Install ginkgo")
	d, err := filepath.Abs("./tools")
	if err != nil {
		return err
	}
	env := map[string]string{"GOBIN": d}
	return sh.RunWith(env, "go", "install", "github.com/onsi/ginkgo/ginkgo")
}

// Start Start the server
func Start() error {
	mg.Deps(Build)
	logStep("Start server")
	return sh.RunV("./bin/secret-server",
		"--port", "8080",
		"--cert", "acceptance_test/testdata/cert.pem",
		"--key", "acceptance_test/testdata/key.pem")
}

// Remove binaries and tools
func Clean() error {
	logStep("Clean")
	err := sh.Rm("./bin")
	if err != nil {
		return err
	}
	return sh.Rm("./tools")
}

func logStep(step string) {
	fmt.Printf("==> %s\n", step)
}
