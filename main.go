package main

import (
	"os"

	"github.com/wagoodman/dive/cmd"
)

// main is the entry point for the dive application.
// dive is a tool for exploring each layer in a docker image,
// investigating the contents, and discovering ways to shrink
// the size of your Docker/OCI image.
func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
