package main

import (
	"os"
	"runtime"

	"kardex-pdf-service/internal/config"
	"kardex-pdf-service/internal/server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := config.Default()
	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}

	srv := server.New(cfg)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
