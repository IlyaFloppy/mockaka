package stringscacheapp

import (
	"context"
	"log"
	"sync"
)

// Config is a set of options for StringsCache app.
type Config struct {
	Address               string
	StringsServiceAddress string
}

// Main creates and runs StringsCache app. It returns stop function when app is ready.
func Main(ctx context.Context, c Config) (stop func()) {
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)

	s := NewService(c.Address, c.StringsServiceAddress)

	go func() {
		defer wg.Done()
		if err := s.Run(ctx); err != nil {
			log.Fatalln(err.Error())
		}
	}()

	<-s.Ready()

	return func() {
		cancel()
		wg.Wait()
	}
}
