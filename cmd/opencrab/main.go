// Package main is the entry point for OpenCrab - a secure personal AI assistant gateway.
// OpenCrab is a Go rewrite of OpenClaw with enhanced security and hardening.
package main

import (
	"os"

	"github.com/h3h3da/opencrab/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
