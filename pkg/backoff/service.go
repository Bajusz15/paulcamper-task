package backoff

import (
	"math"
	"time"
)

type Backoff struct {
	attempt int
	Max     time.Duration
}

func NewService() *Backoff {
	return &Backoff{attempt: 0}
}

func (b *Backoff) Try(f func() error) {
	// Perform some delay based on the current b.attempt, likely using time.Sleep().
	d := b.Duration()
	time.Sleep(d)

	if err := f(); err == nil {
		// Reset the backoff counter on success.
		b.Reset()
	} else {
		b.attempt++
	}
}

func (b *Backoff) Reset() {
	b.attempt = 0
}

func (b *Backoff) Duration() time.Duration {
	if b.attempt == 0 {
		return 0
	}
	durf := math.Pow(2, float64(b.attempt))
	dur := time.Duration(durf)
	if dur > b.Max {
		return b.Max
	}
	return dur

}
