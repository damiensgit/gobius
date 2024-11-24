package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	defaultVersion = "0.8.19" // Default version if not specified in env
)

var solcURLs = map[string]string{
	"windows": "https://github.com/ethereum/solidity/releases/download/v%s/solc-windows.exe",
	"darwin":  "https://github.com/ethereum/solidity/releases/download/v%s/solc-macos",
	"linux":   "https://github.com/ethereum/solidity/releases/download/v%s/solc-static-linux",
}

func getVersion() string {
	// Check for environment variable override
	if version := os.Getenv("SOLC_VERSION"); version != "" {
		// Strip 'v' prefix if present
		return strings.TrimPrefix(version, "v")
	}
	return defaultVersion
}

func checkVersion(binPath string, version string) bool {
	out, err := exec.Command(binPath, "--version").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), version)
}

func main() {

	// Add command line flag handling
	for _, arg := range os.Args[1:] {
		if arg == "--print-version" {
			fmt.Print(getVersion())
			return
		}
	}

	version := getVersion()

	// Create bin directory if it doesn't exist
	binDir := "bin"
	if err := os.MkdirAll(binDir, 0755); err != nil {
		fmt.Printf("Failed to create bin directory: %v\n", err)
		os.Exit(1)
	}

	// Get the URL for the current OS
	url, ok := solcURLs[runtime.GOOS]
	if !ok {
		fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
		os.Exit(1)
	}
	url = fmt.Sprintf(url, version)

	// Define the output filename
	filename := filepath.Join(binDir, "solc")
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}

	// Check if correct version is already installed
	if _, err := os.Stat(filename); err == nil {
		if checkVersion(filename, version) {
			fmt.Printf("solc v%s is already installed\n", version)
			return
		}
	}

	// Download the file
	fmt.Printf("Downloading solc v%s for %s...\n", version, runtime.GOOS)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to download solc: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to download solc v%s: HTTP %d\n", version, resp.StatusCode)
		os.Exit(1)
	}

	// Create the output file
	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	// Copy the content
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Failed to write solc binary: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully installed solc v%s to %s\n", version, filename)
}
