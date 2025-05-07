// В exchange/live_client.go
package exchange

import (
	"marketflow/internal/domain"
)

type LiveClient struct {
	exchange string
}

func NewTCPClient(exchange, addr string) *LiveClient {
	return &LiveClient{
		exchange: exchange,
	}
}

func (c *LiveClient) Start(out chan<- domain.PriceUpdate) error {
	// TODO: Реальная реализация подключения к сокету
	// Пока просто используем генератор как заглушку
	gen := NewTestGenerator(c.exchange + "_LIVE")
	return gen.Start(out)
}

func (c *LiveClient) Stop() error {
	// TODO: Завершение соединения
	return nil
}
