package writer

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/ezio-noir/bluefox/internal/message"
	"github.com/ezio-noir/bluefox/internal/runner"
)

type Writer interface {
	Receiver(chan message.ResultMessage) func()
}

type WriteManager struct {
	DirPath string
	Name    string
	Formats ExtBitset

	writers []Writer
}

func NewWriteManager(dirPath string, name string) *WriteManager {
	if info, err := os.Stat(dirPath); err != nil || !info.IsDir() {
		return nil
	}

	return &WriteManager{
		DirPath: dirPath,
		Name:    name,
		Formats: 0,
	}
}

func (m *WriteManager) hasFormat(ext FileExt) bool {
	return m.Formats&(1<<ext) != 0
}

func (m *WriteManager) addWriter(w Writer) {
	m.writers = append(m.writers, w)
}

func (m *WriteManager) AddFormat(ext string) error {
	extBit := extFromString(ext)

	if extBit == None {
		return fmt.Errorf("invalid format %s", ext)
	}
	if m.hasFormat(extBit) {
		return fmt.Errorf("already has format %s", ext)
	}

	m.Formats &= (1 << extBit)

	switch extBit {
	case Standard:
		m.addWriter(NewConsoleWriter())
	case TXT:
		m.addWriter(NewTXTWriter(path.Join(m.DirPath, fmt.Sprintf("%s.%s", m.Name, "txt"))))
	}

	return nil
}

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
