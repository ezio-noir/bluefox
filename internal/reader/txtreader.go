package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TXTReader struct {
	filePath string
}

func NewTXTReader(filePath string) *TXTReader {
	return &TXTReader{
		filePath: filePath,
	}
}

func (r *TXTReader) Run(outChan chan string) chan struct{} {
	done := make(chan struct{})

	go func() {
		defer func() {
			close(outChan)
			close(done)
		}()

		file, err := os.Open(r.filePath)
		if err != nil {
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.TrimSpace(scanner.Text())
			if word != "" {
				fmt.Printf("Reader emits %s\n", word)
				outChan <- word
			}
		}
	}()

	return done
}
