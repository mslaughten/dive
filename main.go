package main

import (
	"fmt"
	"os"

	"github.com/wagoodman/dive/cmd"
)

// main is the entry point for the dive application.
// dive is a tool for exploring each layer in a docker image,
// investigating the contents, and discovering ways to shrink
// the size of your Docker/OCI image.
//
// Personal fork: added non-zero exit code message for easier debugging
// when dive fails in CI/scripting contexts.
//
// Note: set DIVE_CI=true in your environment to run in CI mode,
// which skips the interactive TUI and just prints the image efficiency report.
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "dive exited with error: %v\n", err)
		os.Exit(1)
	}
}
