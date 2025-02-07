package cmd

import (
	"github.com/ezio-noir/bluefox/internal/dir"
	"github.com/ezio-noir/bluefox/internal/message"
	"github.com/ezio-noir/bluefox/internal/reader"
	"github.com/ezio-noir/bluefox/internal/runner"
	"github.com/ezio-noir/bluefox/internal/writer"
	"github.com/spf13/cobra"
)

// Timeout for requests (seconds)
var timeoutSec uint

var dirCmd = &cobra.Command{
	Use:   "dir <baseurl>",
	Short: "Bruteforce directories",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := requireWordlist()
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		txtReader := reader.NewTXTReader(wordlistPath)

		baseURL := args[0]
		bruteforcer := dir.NewDirBruteforcer(baseURL, timeoutSec, numThreads)

		writeManager := writer.NewWriteManager(outputDir, outputName)
		if outputPath == "" {
			writeManager.AddFormat("std")
		} else {
			for _, format := range outputFormats {
				writeManager.AddFormat(format)
			}
		}

		readChan := make(chan string, numThreads)
		writeChan := make(chan message.ResultMessage)

		runner.WaitAll(
			runner.Run(txtReader.Emitter(readChan)),
			runner.Run(bruteforcer.Runner(readChan, writeChan)),
			runner.Run(writeManager.Receiver(writeChan)),
		)
	},
}

func init() {
	dirCmd.PersistentFlags().UintVar(&timeoutSec, "timeout", 5, "timeout for HTTP request (default 5 seconds)")
}
