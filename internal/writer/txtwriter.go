package writer

import (
	"os"

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

// func (w *TXTWriter) Run(ch chan message.ResultMessage, wg *sync.WaitGroup) {
// 	go func() {
// 		defer wg.Done()

// 		file, err := os.Create(w.filePath)
// 		if err != nil {

// 		}
// 		defer file.Close()

// 		for message := range ch {
// 			file.WriteString(message.String())
// 			file.WriteString("\n")
// 		}
// 	}()
// }

// func (w *TXTWriter) Run(inChannel chan message.ResultMessage) <-chan struct{} {
// 	done := make(chan struct{})

// 	go func() {
// 		defer close(done)

// 		file, err := os.Create(w.filePath)
// 		if err != nil {

// 		}
// 		defer file.Close()

// 		for message := range inChannel {
// 			file.WriteString(message.String())
// 			file.WriteString("\n")
// 		}
// 	}()

// 	return done
// }

func (w *TXTWriter) Receiver(inChannel chan message.ResultMessage) func() {
	return func() {
		file, err := os.Create(w.filePath)
		if err != nil {
			return
		}
		defer file.Close()

		for message := range inChannel {
			file.WriteString(message.String())
			file.WriteString("\n")
		}
	}
}
