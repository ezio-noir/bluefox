package cmd

import (
	"github.com/ezio-noir/bluefox/internal/dns"
	"github.com/ezio-noir/bluefox/internal/message"
	"github.com/ezio-noir/bluefox/internal/reader"
	"github.com/ezio-noir/bluefox/internal/runner"
	"github.com/ezio-noir/bluefox/internal/writer"
	"github.com/spf13/cobra"
)

var dnsCommand = &cobra.Command{
	Use:   "dns",
	Short: "Bruteforce subdomains",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		txtReader := reader.NewTXTReader(wordlistPath)

		domain := args[0]
		bruteforcer := dns.NewDNSBruteforcer(domain, numThreads)

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
