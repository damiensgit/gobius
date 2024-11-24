//go:build tools
// +build tools

package main

import (
	_ "github.com/ethereum/go-ethereum/cmd/abigen"
)

// This file ensures tool dependencies are kept in go.mod
// The blank import above prevents the dependency from being removed
