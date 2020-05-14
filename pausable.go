package pausable

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Pausable struct {
	mu  *sync.Mutex
	sem *semaphore.Weighted
}

func New() *Pausable {
	return &Pausable{
		mu:  &sync.Mutex{},
		sem: semaphore.NewWeighted(1),
	}
}

// Toggle changes the status, and returns the new status (true if locked, false if unlocked)
func (p *Pausable) Toggle() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	acquired := p.sem.TryAcquire(1)
	if acquired {
		return true
	}
	// If failed to acquire, it means that it's locked, and thus we unlock it by releasing:
	p.sem.Release(1)
	return false
}

// IsPaused tells whether the Pausable is currently paused.
func (p *Pausable) IsPaused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	cannotAcquire := !p.sem.TryAcquire(1)
	if cannotAcquire {
		return true
	}
	// If succeeded to acquire, then need to release:
	p.sem.Release(1)
	return false
}

// Wait returns immediately if the Pausable is not paused;
// if the Pausable is paused, then it will wait for unpausing.
func (p *Pausable) Wait() {
	// if locked, then wait for unlock; otherwise return immediately
	p.sem.Acquire(context.Background(), 1)
	p.sem.Release(1)
}
