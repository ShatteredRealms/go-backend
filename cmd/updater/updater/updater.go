package updater

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	version = "0.0.1"
	rootCmd = &cobra.Command{
		Use:   "updater",
		Short: "Updater - a simple cli to create metadata about a release patch for Unreal Engine games",
		Long: `Updater is a simple CLI that works in conjunction with Unreal Engine to generate metadata files
that can be hosted on a CDN to allow clients to know what files are needed for downloading a release
or applying patches.`,
		Version: version,
	}
)

func Execute() {
	initGenerate()
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
