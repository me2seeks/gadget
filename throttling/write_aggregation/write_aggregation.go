package writeaggregation

import (
	"fmt"
	"sync"

	"go.uber.org/multierr"
)

type SyncJob[T any] struct {
	*sync.Cond
	holding   int32
	err       error
	syncPoint *sync.Once
	syncFunc  func(T) error
}

func NewSyncJob[T any](fn func(T) error) *SyncJob[T] {
	return &SyncJob[T]{
		Cond:      sync.NewCond(&sync.Mutex{}),
		holding:   0,
		syncPoint: &sync.Once{},
		syncFunc:  fn,
	}
}

func (s *SyncJob[T]) Do(obj T) error {
	s.L.Lock()
	if s.holding > 0 {
		fmt.Println("sync wait")
		s.Wait()
	}
	s.holding++
	once := s.syncPoint
	s.L.Unlock()
	once.Do(func() {
		s.err = multierr.Append(s.err, s.syncFunc(obj))
		s.L.Lock()
		fmt.Printf("holding:%v\n", s.holding)
		s.holding = 0
		s.syncPoint = &sync.Once{}
		s.Broadcast()
		s.L.Unlock()
	})
	return s.err
}
