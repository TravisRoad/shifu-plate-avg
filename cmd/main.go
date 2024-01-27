package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/TravisRoad/shifu-plate-avg/internal/plate"
)

const (
	EnvPlateEndpoint = "PLATE_ENDPOINT"
	EnvInterval      = "INTERVAL"
)

var (
	PlateEndpoint = "http://deviceshifu-plate-reader.deviceshifu.svc.cluster.local/get_measurement"
	Interval      = 1000
)

func init() {
	if plateEndpoint, ok := os.LookupEnv(EnvPlateEndpoint); ok {
		PlateEndpoint = plateEndpoint
	}
	if interval, ok := os.LookupEnv(EnvInterval); ok {
		x, _ := strconv.Atoi(interval)
		Interval = x
	}
}

func Exit() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	time.Sleep(2 * time.Second)
	os.Exit(0)
}

func main() {
	plate := &plate.Plate{
		URL: PlateEndpoint,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 轮询，默认时间间隔为 1000 ms
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recover", r)
			}
		}()
		plate.Poll(ctx, time.Duration(Interval)*time.Millisecond)
	}()
	go Exit()

	// health api
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Error starting health check server:", err)
	}
}
