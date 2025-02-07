package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Path to wordlist (txt)
var wordlistPath string

// Path to output (outputdir/outputname)
var outputPath string

// Path to output directory
var outputDir string

// Output name
var outputName string

// Output file formats
var outputFormats []string

// Number of threads
var numThreads int

var rootCmd = &cobra.Command{
	Use: "bluefox",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		outputDir = filepath.Dir(outputDir)
		baseName := filepath.Base(outputPath)
		outputName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&wordlistPath, "wordlist", "w", "", "wordlist (required)")

	rootCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "", "output path")

	rootCmd.PersistentFlags().StringSliceVarP(&outputFormats, "output-format", "f", make([]string, 0), "output formats")

	rootCmd.PersistentFlags().IntVarP(&numThreads, "threads", "t", 5, "number of threads")

	rootCmd.AddCommand(dnsCmd)
	rootCmd.AddCommand(dirCmd)
}

func requireWordlist() error {
	if wordlistPath == "" {
		return fmt.Errorf("missing wordlist")
	}
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func GenDoc() {
	doc.GenMarkdownTree(rootCmd, "./doc")
}
