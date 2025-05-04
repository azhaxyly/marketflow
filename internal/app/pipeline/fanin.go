package pipeline

import (
	"log"
	"marketflow/internal/domain"
)

func FanIn(clients []domain.ExchangeClient) <-chan domain.PriceUpdate {
	out := make(chan domain.PriceUpdate, 100)

	for _, client := range clients {
		go func(c domain.ExchangeClient) {
			err := c.Start(out)
			if err != nil {
				log.Printf("[FAN-IN] error starting client: %v", err)
			}
		}(client)
	}

	return out
}
