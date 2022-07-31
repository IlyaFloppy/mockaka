package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/IlyaFloppy/mockaka/example/stringscache/stringscacheapp"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	config := stringscacheapp.Config{
		Address:               os.Getenv("STRINGS_CACHE_ADDRESS"),
		StringsServiceAddress: os.Getenv("STRINGS_ADDRESS"),
	}

	stringscacheapp.Main(ctx, config)
}
