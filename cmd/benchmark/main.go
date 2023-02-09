package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	version = "0.0.1"
	rootCmd = &cobra.Command{
		Use:   "loady",
		Short: "generates loads on sro microservices",
		Long: `A CLI tool for generating loads on SRO microservices. Loads can be directed an specific services, or all
services. The load will make real requests, then attempt to clean up itself.`,
		Version: version,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
