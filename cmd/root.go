package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var wordlistPath string
var outputPath string
var outputDir string
var outputName string
var outputFormats []string

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
	rootCmd.MarkPersistentFlagRequired("wordlist")

	rootCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", ".", "output path")

	rootCmd.PersistentFlags().StringSliceVarP(&outputFormats, "output-format", "f", make([]string, 0), "output formats")

	rootCmd.AddCommand(dnsCommand)
}

func Execute() error {
	return rootCmd.Execute()
}
