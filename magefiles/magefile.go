//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/magefile/mage/sh"
)

// Test runs all tests with race detection and coverage.
func Test() error {
	fmt.Println("Running tests with race detection and coverage...")
	return sh.RunV("go", "test", "-race", "-coverprofile=coverage.out", "./...")
}

// Coverage runs tests and opens the HTML coverage report in a browser.
func Coverage() error {
	if err := Test(); err != nil {
		return err
	}

	fmt.Println("Generating HTML coverage report...")
	if err := sh.RunV("go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html"); err != nil {
		return err
	}

	fmt.Println("Opening coverage report...")
	return openBrowser("coverage.html")
}

// Lint runs golangci-lint.
func Lint() error {
	fmt.Println("Running golangci-lint...")
	return sh.RunV("golangci-lint", "run", "./...")
}

// Build verifies the module compiles.
func Build() error {
	fmt.Println("Verifying module compiles...")
	return sh.RunV("go", "build", "./...")
}

// openBrowser opens the specified file in the default browser.
func openBrowser(file string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", file)
	case "darwin":
		cmd = exec.Command("open", file)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", file)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}
