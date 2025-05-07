package mode

import (
	"context"
	"errors"
	"sync"

	"marketflow/internal/adapters/exchange"
	"marketflow/internal/domain"
	"marketflow/internal/logger"
)

type Mode string

const (
	Live Mode = "live"
	Test Mode = "test"
)

type Manager struct {
	mu         sync.Mutex
	mode       Mode
	cancelFunc context.CancelFunc
}

func NewManager() *Manager {
	return &Manager{mode: Test} // можно по умолчанию запускать в Test
}

func (m *Manager) Start(ctx context.Context, out chan<- domain.PriceUpdate, mode Mode) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.mode == mode {
		logger.Info("mode already started", "mode", mode)
		return nil
	}

	if m.cancelFunc != nil {
		logger.Info("stopping previous mode")
		m.cancelFunc()
	}

	// Создаем новый контекст и cancel функцию
	ctx, cancel := context.WithCancel(ctx)
	m.cancelFunc = cancel
	m.mode = mode

	var clients []domain.ExchangeClient
	if mode == Test {
		clients = []domain.ExchangeClient{
			exchange.NewTestGenerator("ex1"),
			exchange.NewTestGenerator("ex2"),
			exchange.NewTestGenerator("ex3"),
		}
	} else if mode == Live {
		clients = []domain.ExchangeClient{
			exchange.NewTCPClient("ex1", "localhost:40101"),
			exchange.NewTCPClient("ex2", "localhost:40102"),
			exchange.NewTCPClient("ex3", "localhost:40103"),
		}
	} else {
		return errors.New("invalid mode")
	}

	go func() {
		for _, c := range clients {
			go func(client domain.ExchangeClient) {
				err := client.Start(out)
				if err != nil {
					logger.Error("failed to start client", "error", err)
				}
			}(c)
		}
	}()

	logger.Info("started mode", "mode", mode)
	return nil
}

func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cancelFunc != nil {
		logger.Info("stopping mode")
		m.cancelFunc()
	}
}

func (m *Manager) Current() Mode {
	m.mu.Lock()
	defer m.mu.Unlock()
	logger.Info("current mode", "mode", m.mode)
	return m.mode
}
