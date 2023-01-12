package main

import (
	"flag"
	"github.com/aviscode/Coins/Pricing-service/data"
	"github.com/aviscode/Coins/Pricing-service/fetcher"
	"github.com/aviscode/Coins/Pricing-service/server"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	ApiKey, ServerAddr string
)

func initFlags() {
	flag.StringVar(&ApiKey, "ApiKey", "08b53f4d-0a5d-452e-b67f-3af8af0d035c", "The api-key for coinmarketcap services ")
	flag.StringVar(&ServerAddr, "ServerAddr", ":9091", "The the address  the service is listening on ")
	flag.Parse()
}

func verifyFlags() error {
	if ApiKey == "" {
		return fetcher.ErrEmptyApiKey
	}
	return nil
}

func main() {
	initFlags()
	if err := verifyFlags(); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Pricing service starting...")
	coinsData := data.NewCoinsDataStore()
	dataFetcher, err := fetcher.NewFetcher(ApiKey)
	if err != nil {
		logrus.Fatal(err)
	}

	s, err := server.NewGrpcServer(ServerAddr, coinsData, dataFetcher)
	if err != nil {
		logrus.Fatal(err)
	}

	//grpcErrCh := make(chan error, 0)
	go s.StartGrpcServer()
	sCh := make(chan os.Signal, 1)
	signal.Notify(sCh, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case sig := <-sCh:
		log.Println("received OS signal", sig, "Shutting down gRPC")
		s.Stop()
	}
}
