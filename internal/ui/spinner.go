package ui
import (
	"fmt"
	"time"
)

type Spinner struct {
	stopChan chan struct{}
	message  string
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		stopChan: make(chan struct{}),
		message:  message,
	}
}

func (s *Spinner) Start() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		for {
			for _, frame := range frames {
				select {
				case <-s.stopChan:
					return
				default:
					fmt.Printf("\r%s %s", ColorCyan(frame), s.message)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()
}

func (s *Spinner) Stop() {
	close(s.stopChan)
	fmt.Print("\r\033[K")
}
