package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dive [IMAGE]",
	Short: "A tool for exploring each layer in a docker image",
	Long: `dive is a tool for exploring a docker image, layer contents,
and discovering ways to shrink the size of your Docker/OCI image.

Usage:
  dive <image-tag> [flags]
  dive <subcommand>`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return runDive(args[0])
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dive.yaml)")
	rootCmd.PersistentFlags().String("source", "docker", "the source of the image to analyze. Available sources: docker, podman, docker-archive, oci-archive, oci-dir")
	rootCmd.PersistentFlags().Bool("ci", false, "skip the interactive TUI and validate against CI rules")
	rootCmd.PersistentFlags().String("ci-config", ".dive-ci", "path to the CI config file")
	rootCmd.PersistentFlags().Bool("json", false, "output results as JSON (for CI mode)")
	rootCmd.PersistentFlags().Bool("lowestEfficiency", false, "only report layers with the lowest efficiency")

	// Bind flags to viper
	_ = viper.BindPFlag("source", rootCmd.PersistentFlags().Lookup("source"))
	_ = viper.BindPFlag("ci", rootCmd.PersistentFlags().Lookup("ci"))
	_ = viper.BindPFlag("ci-config", rootCmd.PersistentFlags().Lookup("ci-config"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".dive" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".dive")
	}

	viper.SetEnvPrefix("DIVE")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// runDive is the main entry point for analyzing a container image.
func runDive(imageRef string) error {
	source := viper.GetString("source")
	ciMode := viper.GetBool("ci")

	fmt.Printf("Analyzing image: %s (source: %s)\n", imageRef, source)

	if ciMode {
		fmt.Println("Running in CI mode...")
		// CI mode: analyze and validate without TUI
		return runCI(imageRef, source)
	}

	// Interactive TUI mode
	return runUI(imageRef, source)
}

// runCI performs image analysis and validates against CI rules without a TUI.
func runCI(imageRef, source string) error {
	// TODO: implement CI analysis pipeline
	_ = source
	_ = imageRef
	return nil
}

// runUI launches the interactive terminal UI for image exploration.
func runUI(imageRef, source string) error {
	// TODO: implement TUI launch
	_ = source
	_ = imageRef
	return nil
}
