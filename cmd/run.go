package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"marketflow/internal/api"
	"marketflow/internal/app/mode"
	"marketflow/internal/app/pipeline"
	"marketflow/internal/domain"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	priceChan := make(chan domain.PriceUpdate, 1000)
	manager := mode.NewManager()

	if err := manager.Start(ctx, priceChan, mode.Test); err != nil {
		log.Println(err)
	}

	workerChans := pipeline.FanOut(priceChan, 3)

	for i, ch := range workerChans {
		go func(id int, ch <-chan domain.PriceUpdate) {
			for update := range ch {
				fmt.Printf("[Worker %d] %s %s %.2f\n", id, update.Exchange, update.Pair, update.Price)
			}
		}(i, ch)
	}

	go func() {
		http.HandleFunc("/mode/live", api.HandleLiveMode(manager)) 
		http.HandleFunc("/mode/test", api.HandleTestMode(manager)) 
		http.HandleFunc("/health", api.HealthCheckHandler)           
		log.Fatal(http.ListenAndServe(":8080", nil))            
	}()

	// Даем серверу время на старте
	time.Sleep(10 * time.Second)
}
