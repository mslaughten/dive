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
//
// Tip: combine with `docker build` using:
//   docker build -t my-image . && dive my-image
//
// Personal note: I also find it useful to alias this in ~/.bashrc:
//   alias dive='DIVE_CI=true dive'   # for quick non-interactive checks
//
// Personal note: use `dive --source docker-archive image.tar` to analyze
// exported tarballs offline without needing the Docker daemon running.
//
// Personal note: use `dive --ci-config .dive-ci.yml` to specify a custom
// CI config file path, handy when managing multiple projects with different
// efficiency thresholds in the same repo.
//
// Personal note: use `dive --ignore-errors` to continue even when some
// layers can't be fully analyzed (useful for scratch-based or distroless images).
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "dive exited with error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Tip: run with --help for usage information\n")
		fmt.Fprintf(os.Stderr, "Tip: if analyzing a private image, ensure you are logged in via `docker login`\n")
		os.Exit(1)
	}
}
