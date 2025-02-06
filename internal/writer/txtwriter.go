package writer

import (
	"os"
	"sync"

	"github.com/ezio-noir/bluefox/internal/message"
)

type TXTWriter struct {
	filePath string
}

func NewTXTWriter(filePath string) *TXTWriter {
	return &TXTWriter{
		filePath: filePath,
	}
}

func (w *TXTWriter) Run(ch chan message.ResultMessage, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		file, err := os.Create(w.filePath)
		if err != nil {

		}
		defer file.Close()

		for message := range ch {
			file.WriteString(message.String())
			file.WriteString("\n")
		}
	}()
}
