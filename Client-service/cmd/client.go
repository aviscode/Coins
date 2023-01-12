package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials/insecure"
	"runtime"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pb "github.com/aviscode/Coins/Pricing-service/grpc"
)

var (
	ServerAddr  string
	coinsPrices = sync.Map{}
)

func initFlags() {
	flag.StringVar(&ServerAddr, "ServerAddr", ":9091", "The the address the service to connect to")
	flag.Parse()
}

func initSymbles() {
	coinsPrices.Store("BTC", 0.0)
	coinsPrices.Store("ETH", 0.0)
	coinsPrices.Store("USDT", 0.0)
}

type Client struct {
	c pb.ServiceClient
}

func (c *Client) Close() error {
	return c.Close()
}

// NewGrpcClient create new gRPC client that connects to Pricing-service requests the coins prices.
func NewGrpcClient(addr string) (*Client, error) {
	keepaliveParams := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	})
	cert := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(addr, cert, keepaliveParams, grpc.WithBlock())
	return &Client{pb.NewServiceClient(conn)}, err
}

func (c *Client) fetchSymbolPrice(symbol string) (float64, error) {
	response, err := c.c.GetPrice(context.Background(), &pb.GetCoinPrice{Symbol: symbol})
	if err != nil {
		return 0, err
	}
	return response.Price, err
}

func (c *Client) fetchAllSymbols() {
	for {
		coinsPrices.Range(func(k interface{}, v interface{}) bool {
			go func(kk string, vv float64) {
				newPrice, err := c.fetchSymbolPrice(kk)
				if err != nil {
					logrus.Errorf("failed to fetch symbol %s with err %v", kk, err)
				}
				coinsPrices.Store(kk, newPrice)
				logrus.Infof("%s: %f, %f, %.2f%%", kk, vv, newPrice, (1-(vv/newPrice))*100)
			}(k.(string), v.(float64))
			return true
		})
		time.Sleep(time.Minute)
	}
}

func main() {
	initFlags()
	initSymbles()
	c, err := NewGrpcClient(ServerAddr)
	if err != nil {
		logrus.Fatal(err)
	}
	c.fetchAllSymbols()
	runtime.Goexit()
}
