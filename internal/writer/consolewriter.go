package writer

import (
	"fmt"

	"github.com/ezio-noir/bluefox/internal/message"
)

type ConsoleWriter struct{}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

// func (w *ConsoleWriter) Run(ch chan message.ResultMessage, wg *sync.WaitGroup) {
// 	go func() {
// 		defer wg.Done()

// 		for message := range ch {
// 			fmt.Printf("%s\n", message.String())
// 		}
// 	}()
// }

// func (w *ConsoleWriter) Run(ch chan message.ResultMessage) <-chan struct{} {
// 	done := make(chan struct{})

// 	go func() {
// 		defer close(done)

// 		for message := range ch {
// 			fmt.Printf("%s\n", message.String())
// 		}
// 	}()

// 	return done
// }

func (w *ConsoleWriter) Receiver(inChannel chan message.ResultMessage) func() {
	return func() {
		for entry := range inChannel {
			fmt.Printf("%s\n", entry.String())
		}
	}
}
