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

// Run: lint test
func All() {
	mg.SerialDeps(Lint, Test)
}

// Run linters using golangci-lint
func Lint() error {
	mg.Deps(InstallGolangCILint)
	fmt.Println("=> Lint")
	return sh.Run("./tools/golangci-lint", "run")
}

// Run tests using ginkgo test runner
func Test() error {
	mg.Deps(InstallGinkgo)
	fmt.Println("=> Test")
	return sh.Run("./tools/ginkgo", "-race", "-cover", "-r")
}

// Install golangci-lint in tools directory
func InstallGolangCILint() error {
	needUpdate, err := target.Path("./tools/golangci-lint")
	if err != nil || !needUpdate {
		return err
	}

	fmt.Println("=> Install golangci-lint")

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
	u := fmt.Sprintf("https://github.com/golangci/golangci-lint/releases/download/v%s/%s.tar.gz", version, asset)

	err = sh.Run("curl", "-sSfL", u, "-o", tarball.Name())
	return sh.Run("tar", "xf", tarball.Name(), "-C", "tools", "--strip-components", "1", fmt.Sprintf("%s/golangci-lint", asset))
}

// Install ginkgo in tools directory
func InstallGinkgo() error {
	needUpdate, err := target.Path("./tools/ginkgo", "./go.mod")
	if err != nil || !needUpdate {
		return err
	}

	fmt.Println("=> Install ginkgo")
	d, err := filepath.Abs("./tools")
	if err != nil {
		return err
	}
	env := map[string]string{"GOBIN": d}
	return sh.RunWith(env, "go", "install", "github.com/onsi/ginkgo/ginkgo")
}

// Remove tools directory
func Clean() error {
	fmt.Println("=> Clean")
	return sh.Rm("./tools")
}
