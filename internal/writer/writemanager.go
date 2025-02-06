package writer

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/ezio-noir/bluefox/internal/message"
)

type Writer interface {
	Run(chan message.ResultMessage, *sync.WaitGroup)
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

func (m *WriteManager) Run(inChan chan message.ResultMessage) chan struct{} {
	done := make(chan struct{})
	wg := new(sync.WaitGroup)
	channels := map[Writer]chan message.ResultMessage{}

	wg2 := new(sync.WaitGroup)

	for _, writer := range m.writers {
		writerChan := make(chan message.ResultMessage)
		channels[writer] = writerChan

		wg.Add(1)
		writer.Run(writerChan, wg)
	}

	go func() {
		for entry := range inChan {
			wg2.Add(len(channels))
			for _, channel := range channels {
				go func(c chan message.ResultMessage) {
					channel <- entry
					wg2.Done()
				}(channel)
			}
		}

		wg2.Wait()
		for _, channel := range channels {
			close(channel)
		}

		wg.Wait()
		close(done)
	}()

	return done
}
