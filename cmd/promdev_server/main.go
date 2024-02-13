package main

import (
	"context"
	"log"
	"time"

	"github.com/barakmich/promdev_sd"
	"github.com/spf13/pflag"
)

var (
	hostport      = pflag.String("hostport", "0.0.0.0:9111", "Host and port to listen on")
	timeoutPeriod = pflag.Duration("timeout-period", 2*time.Minute, "Timeout until server drops targets that haven't posted in a while.")
)

func main() {
	pflag.Parse()
	config := promdev_sd.Config{
		TimeoutPeriod: *timeoutPeriod,
		GCPeriod:      30 * time.Second,
	}
	srv, err := promdev_sd.NewServer(config)
	if err != nil {
		log.Fatalln("Couldn't create server:", err)
	}
	srv.RunGCThread(context.Background())
	err = srv.ListenAndServe(*hostport)
	if err != nil {
		log.Fatalln("Couldn't listen with server:", err)
	}
}
