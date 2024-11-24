package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	requiredGoVersion    = "1.22.9"
	defaultAbigenVersion = "v1.13.15" // Default version if not specified in env
)

func getAbigenVersion() string {
	if version := os.Getenv("ABIGEN_VERSION"); version != "" {
		// Ensure version starts with 'v'
		if !strings.HasPrefix(version, "v") {
			return "v" + version
		}
		return version
	}
	return defaultAbigenVersion
}

func ensureCorrectDirectory() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	if strings.HasSuffix(wd, filepath.Join("tools", "abigen")) {
		return fmt.Errorf("this script must be run from the project root directory using: go run tools/abigen/install.go")
	}

	toolsDir := filepath.Join("tools", "abigen")
	if _, err := os.Stat(toolsDir); err != nil {
		return fmt.Errorf("tools/abigen directory not found - are you in the project root directory?")
	}

	return nil
}

func ensureGoVersion() error {
	goCmd := fmt.Sprintf("go%s", requiredGoVersion)
	if _, err := exec.LookPath(goCmd); err != nil {
		fmt.Printf("Installing Go %s...\n", requiredGoVersion)
		cmd := exec.Command("go", "install", fmt.Sprintf("golang.org/dl/go%s@latest", requiredGoVersion))
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to install Go %s: %v\n%s", requiredGoVersion, err, output)
		}

		dlCmd := exec.Command(goCmd, "download")
		if output, err := dlCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to download Go %s: %v\n%s", requiredGoVersion, err, output)
		}
	}
	return nil
}

func main() {
	if err := ensureCorrectDirectory(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if err := ensureGoVersion(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	abigenVersion := getAbigenVersion()
	fmt.Printf("Using abigen version: %s\n", abigenVersion)

	binDir := "bin"
	if err := os.MkdirAll(binDir, 0755); err != nil {
		fmt.Printf("Failed to create bin directory: %v\n", err)
		os.Exit(1)
	}

	absPath, err := filepath.Abs(binDir)
	if err != nil {
		fmt.Printf("Failed to get absolute path: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Installing abigen to: %s\n", absPath)

	goCmd := fmt.Sprintf("go%s", requiredGoVersion)
	cmd := exec.Command(goCmd, "install",
		fmt.Sprintf("github.com/ethereum/go-ethereum/cmd/abigen@%s", abigenVersion))

	cmd.Env = append(os.Environ(),
		"GO111MODULE=on",
		fmt.Sprintf("GOBIN=%s", absPath),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to install abigen: %v\n", err)
		os.Exit(1)
	}

	expectedPath := filepath.Join(absPath, "abigen")
	if runtime.GOOS == "windows" {
		expectedPath += ".exe"
	}

	if _, err := os.Stat(expectedPath); err != nil {
		fmt.Printf("Failed to verify abigen installation at %s: %v\n", expectedPath, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully installed abigen %s to %s using Go %s\n",
		abigenVersion, expectedPath, requiredGoVersion)
}
