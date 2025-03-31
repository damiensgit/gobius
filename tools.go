//go:build tools

package main

// TODO: remove this file once we upgrade to go 1.24 which has
// native support for tools:
// https://tip.golang.org/doc/modules/managing-dependencies#tools
import (
	_ "github.com/go-task/task/v3/cmd/task"
)
