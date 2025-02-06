package cmd

import (
	"fmt"

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
		domain := args[0]
		fmt.Printf("Running the dns command for domain %s...\n", domain)

		txtReader := reader.NewTXTReader(wordlistPath)

		bruteforcer := dns.NewDNSBruteforcer(domain, 5)

		writeManager := writer.WriteManager{}
		writeManager.AddWriter(writer.NewConsoleWriter())

		fmt.Printf("Output path: %s\n", outputDir)
		if outputDir != "" {
			for _, ext := range outputFormats {
				writeManager.AddFileWriter(outputDir, outputName, ext)
			}
		}

		readChan := make(chan string, 5)
		writeChan := make(chan message.ResultMessage)

		// // txtReaderDone := txtReader.Run(readChan)
		// txtReaderDone := runner.Run(txtReader.Emitter(readChan))
		// // bruteforcerDone := bruteforcer.Run(readChan, writeChan)
		// bruteforcerDone := runner.Run(bruteforcer.Runner(readChan, writeChan))
		// // writeManagerDone := writeManager.Run(writeChan)
		// writeManagerDone := runner.Run(writeManager.Receiver(writeChan))

		// <-txtReaderDone
		// <-bruteforcerDone
		// <-writeManagerDone

		runner.WaitAll(
			runner.Run(txtReader.Emitter(readChan)),
			runner.Run(bruteforcer.Runner(readChan, writeChan)),
			runner.Run(writeManager.Receiver(writeChan)),
		)
	},
}
