package writer

import (
	"fmt"
	"sync"

	"github.com/ezio-noir/bluefox/internal/message"
)

type ConsoleWriter struct{}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

func (w *ConsoleWriter) Run(ch chan message.ResultMessage, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		for message := range ch {
			fmt.Printf("%s\n", message.String())
		}
	}()
}
