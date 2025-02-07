package dir

import "fmt"

type DirResultMessage struct {
	StatusCode int
	URL        string
}

func (m DirResultMessage) String() string {
	return fmt.Sprintf("%03d\t%s", m.StatusCode, m.URL)
}

func (m DirResultMessage) CSV() string {
	return fmt.Sprintf("%s,%d", m.URL, m.StatusCode)
}
