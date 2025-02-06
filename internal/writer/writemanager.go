package writer

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/ezio-noir/bluefox/internal/message"
	"github.com/ezio-noir/bluefox/internal/runner"
)

type Writer interface {
	// Run(chan message.ResultMessage) <-chan struct{}
	Receiver(chan message.ResultMessage) func()
}

type WriteManager struct {
	writers []Writer
}

func NewWriteManager() *WriteManager {
	return &WriteManager{}
}

func (m *WriteManager) AddWriter(w Writer) {
	m.writers = append(m.writers, w)
}

func (m *WriteManager) AddFileWriter(dirPath string, fileName string, fileExt string) {
	fullPath := filepath.Join(dirPath, fmt.Sprintf("%s.%s", fileName, fileExt))
	switch fileExt {
	case "txt":
		m.AddWriter(NewTXTWriter(fullPath))
	case "csv":
		return
	}
}

// func (m *WriteManager) Run(inChannel chan message.ResultMessage) chan struct{} {
// 	done := make(chan struct{})
// 	wg := new(sync.WaitGroup)
// 	channels := map[Writer]chan message.ResultMessage{}

// 	wg2 := new(sync.WaitGroup)

// 	// writerDones := make([]<-chan struct{}, len(m.writers))

// 	// for i, writer := range m.writers {
// 	// 	writerChannel := make(chan message.ResultMessage)
// 	// 	channels[writer] = writerChannel

// 	// 	// wg.Add(1)
// 	// 	// writer.Run(writerChan, wg)

// 	// 	writerDones[i] = writer.Run(writerChannel)
// 	// }

// 	go func() {
// 		defer close(done)

// 		numWriters := len(m.writers)
// 		writerDones := make([]<-chan struct{}, numWriters)
// 		for i, writer := range m.writers {
// 			writerChannel := make(chan message.ResultMessage)
// 			channels[writer] = writerChannel
// 			writerDones[i] = writer.Run(writerChannel)
// 		}

// 		entryWritten := new(sync.WaitGroup)
// 		for entry := range inChannel {
// 			// wg2.Add(len(channels))
// 			entryWritten.Add(len(w))
// 			for _, channel := range channels {
// 				go func(c chan message.ResultMessage) {
// 					channel <- entry
// 					wg2.Done()
// 				}(channel)
// 			}
// 		}

// 		wg2.Wait()
// 		for _, channel := range channels {
// 			close(channel)
// 		}

// 		wg.Wait()

// 		for _, writerDone := range writerDones {
// 			<-writerDone
// 		}

// 		// close(done)
// 	}()

// 	return done
// }

// func (m *WriteManager) Run(inChannel chan message.ResultMessage) chan struct{} {
// 	done := make(chan struct{})

// 	go func() {
// 		defer close(done)

// 		numWriters := len(m.writers)
// 		channels := make([]chan message.ResultMessage, numWriters)
// 		writerDones := make([]<-chan struct{}, numWriters)

// 		for i := range numWriters {
// 			writerChannel := make(chan message.ResultMessage)
// 			channels[i] = writerChannel
// 			writerDones[i] = m.writers[i].Run(writerChannel)
// 		}

// 		entryWritten := new(sync.WaitGroup)
// 		for entry := range inChannel {
// 			entryWritten.Add(numWriters)
// 			for _, channel := range channels {
// 				go func(c chan message.ResultMessage) {
// 					channel <- entry
// 					entryWritten.Done()
// 				}(channel)
// 			}
// 			entryWritten.Wait()
// 		}

// 		for _, channel := range channels {
// 			close(channel)
// 		}

// 		for _, writerDone := range writerDones {
// 			<-writerDone
// 		}
// 	}()

// 	return done
// }

func (m *WriteManager) Receiver(inChannel chan message.ResultMessage) func() {
	return func() {
		numWriters := len(m.writers)
		channels := make([]chan message.ResultMessage, numWriters)
		writerDones := make([]<-chan struct{}, numWriters)

		for i := range numWriters {
			writerChannel := make(chan message.ResultMessage)
			channels[i] = writerChannel
			writerDones[i] = runner.Run(m.writers[i].Receiver(writerChannel))
		}

		entryWritten := new(sync.WaitGroup)
		for entry := range inChannel {
			entryWritten.Add(numWriters)
			for _, channel := range channels {
				go func(c chan message.ResultMessage) {
					channel <- entry
					entryWritten.Done()
				}(channel)
			}
			entryWritten.Wait()
		}

		for _, channel := range channels {
			close(channel)
		}

		for _, writerDone := range writerDones {
			<-writerDone
		}
	}
}
