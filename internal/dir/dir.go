package dir

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ezio-noir/bluefox/internal/message"
)

type DirBruteforcer struct {
	BaseURL    string
	TimeoutSec uint

	numWorkers int
}

func NewDirBruteforcer(baseUrl string, timeoutSec uint, numWorkers int) *DirBruteforcer {
	return &DirBruteforcer{
		BaseURL:    baseUrl,
		TimeoutSec: timeoutSec,
		numWorkers: numWorkers,
	}
}

func worker(baseURL string, timeoutSec uint, inChannel chan string, outChannel chan message.ResultMessage, wg *sync.WaitGroup) {
	defer wg.Done()

	workerClient := http.Client{
		Timeout: time.Second * time.Duration(timeoutSec),
	}

	for directory := range inChannel {
		url := fmt.Sprintf("%s/%s", baseURL, directory)

		resp, err := workerClient.Head(url)
		if err != nil {
			outChannel <- &DirResultMessage{
				StatusCode: 0,
				URL:        url,
			}
			continue
		}
		defer resp.Body.Close()

		outChannel <- &DirResultMessage{
			StatusCode: resp.StatusCode,
			URL:        url,
		}
	}
}

func (b *DirBruteforcer) Runner(inChannel chan string, outChannel chan message.ResultMessage) func() {
	return func() {
		defer close(outChannel)

		workersDone := new(sync.WaitGroup)
		workersDone.Add(b.numWorkers)
		for i := 1; i <= b.numWorkers; i++ {
			go worker(b.BaseURL, b.TimeoutSec, inChannel, outChannel, workersDone)
		}

		workersDone.Wait()
	}
}
