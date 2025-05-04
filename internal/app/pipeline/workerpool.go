package pipeline

import (
	"fmt"
	"marketflow/internal/domain"
)

type Worker struct {
	ID     int
	Input  <-chan domain.PriceUpdate
	Cache  domain.Cache
	Output chan<- domain.PriceUpdate
}

func (w *Worker) Start() {
	go func() {
		for update := range w.Input {
			err := w.Cache.SetLatest(update)
			if err != nil {
				fmt.Printf("[Worker %d] Redis error: %v\n", w.ID, err)
			}

			if w.Output != nil {
				w.Output <- update
			}
		}
	}()
}
