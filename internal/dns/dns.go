package dns

import (
	"fmt"
	"net"
	"sync"

	"github.com/ezio-noir/bluefox/internal/message"
)

type DNSBruteforcer struct {
	domain     string
	numWorkers int
}

func NewDNSBruteforcer(domain string, numWorkers int) *DNSBruteforcer {
	return &DNSBruteforcer{
		domain:     domain,
		numWorkers: numWorkers,
	}
}

func worker(domain string, inChan chan string, outChan chan message.ResultMessage, wg *sync.WaitGroup) {
	defer wg.Done()

	for subdomain := range inChan {
		host := fmt.Sprintf("%s.%s", subdomain, domain)
		// fmt.Printf("[DNS] tries %s\n", host)
		if addrs, err := net.LookupHost(host); err != nil {
			outChan <- &DNSResultMessage{
				Success:   false,
				Host:      host,
				Addresses: nil,
			}
		} else {
			outChan <- &DNSResultMessage{
				Success:   true,
				Host:      host,
				Addresses: addrs,
			}
		}
	}
}

func (b *DNSBruteforcer) Runner(inChannel chan string, outChannel chan message.ResultMessage) func() {
	return func() {
		defer close(outChannel)

		workersDone := new(sync.WaitGroup)
		workersDone.Add(b.numWorkers)
		for i := 1; i <= b.numWorkers; i++ {
			go worker(b.domain, inChannel, outChannel, workersDone)
		}

		workersDone.Wait()
	}
}
