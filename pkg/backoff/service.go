package backoff

import (
	"math"
	"time"
)

type Service struct {
	attempt    int
	maxBackoff time.Duration
	maxRetries int
}

func NewService(mb time.Duration, mr int) *Service {
	return &Service{attempt: 0, maxBackoff: mb, maxRetries: mr}
}

func (b *Service) Try(f func() error) {
	// Perform some delay based on the current b.attempt, likely using time.Sleep().
	d := b.Duration()
	time.Sleep(d)

	if err := f(); err == nil {
		// Reset the backoff counter on success.
		b.Reset()
	} else {
		b.attempt++
		if b.attempt > b.maxRetries {
			return
		}
	}
}

func (b *Service) Reset() {
	b.attempt = 0
}

func (b *Service) Duration() time.Duration {
	if b.attempt == 0 {
		return 0
	}
	dur := time.Duration(math.Pow(2, float64(b.attempt))) * 100 * time.Millisecond
	if dur > b.maxBackoff {
		return b.maxBackoff
	}
	return dur

}
