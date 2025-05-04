package cmd

import (
	"fmt"
	"marketflow/internal/adapters/exchange"
	"marketflow/internal/app/pipeline"
	"marketflow/internal/domain"
	"time"
)

func Run() {
	clients := []domain.ExchangeClient{
		exchange.NewTestGenerator("ex1"),
		exchange.NewTestGenerator("ex2"),
	}

	in := pipeline.FanIn(clients)
	workerChans := pipeline.FanOut(in, 3) // 3 worker'а

	for i, ch := range workerChans {
		go func(id int, ch <-chan domain.PriceUpdate) {
			for update := range ch {
				fmt.Printf("[Worker %d] %s %s %.2f\n", id, update.Exchange, update.Pair, update.Price)
			}
		}(i, ch)
	}
	time.Sleep(5 * time.Second)
}
