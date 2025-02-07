package dns

import "fmt"

type DNSResultMessage struct {
	Success   bool
	Host      string
	Addresses []string
}

func (m DNSResultMessage) String() string {
	return fmt.Sprintf("%s\t(%v)", m.Host, m.Addresses)
}

func (m DNSResultMessage) CSV() string {
	if !m.Success || m.Addresses == nil {
		return fmt.Sprintf("%s,0", m.Host)
	}
	return fmt.Sprintf("%s,%d", m.Host, len(m.Addresses))
}
