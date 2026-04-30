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
//
// Personal note: use `dive --lowestEfficiency 0.95` to enforce a stricter
// efficiency threshold than the default 0.9 — I prefer this for production images.
//
// Personal note: use `dive --highestUserWastedPercent 0.05` to fail CI if more
// than 5% of image space is wasted — pairs well with the 0.95 efficiency threshold.
//
// Personal note: use `dive --json output.json` to dump the analysis results to
// a JSON file, useful for post-processing or storing reports as CI artifacts.
//
// Personal note: use `dive --source podman image-name` to analyze Podman images
// directly — handy now that I've been moving away from Docker on my local machine.
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "dive exited with error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Tip: run with --help for usage information\n")
		fmt.Fprintf(os.Stderr, "Tip: if analyzing a private image, ensure you are logged in via `docker login`\n")
		fmt.Fprintf(os.Stderr, "Tip: check DIVE_CI, DIVE_SOURCE, and other env vars if behavior seems unexpected\n")
		os.Exit(1)
	}
}
