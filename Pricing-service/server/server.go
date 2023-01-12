package server

import (
	"context"
	"github.com/aviscode/Coins/Pricing-service/collector"
	"github.com/aviscode/Coins/Pricing-service/data"
	"github.com/aviscode/Coins/Pricing-service/fetcher"
	pb "github.com/aviscode/Coins/Pricing-service/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"syscall"
)

type server struct {
	pb.UnimplementedServiceServer
	lis         net.Listener
	server      *grpc.Server
	CoinsData   *data.Coins
	DataFetcher *fetcher.ClientApi
}

func (s *server) GetPrice(ctx context.Context, request *pb.GetCoinPrice) (*pb.Response, error) {
	price, err := s.CoinsData.GetSymbolPrice(request.Symbol)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Price: price}, nil
}

func (s *server) Stop() {
	s.server.GracefulStop()
}

func createServerInterceptors() []grpc.ServerOption {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    false,
		QuoteEmptyFields: true,
	})
	ui := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(
				logrus.NewEntry(l),
			),
		),
	)
	return []grpc.ServerOption{ui}
}

// StartGrpcServer start listening on the initiated port and blocks to receive new connections
func (s *server) StartGrpcServer() {
	if err := s.server.Serve(s.lis); err != nil {
		logrus.Errorf("server terminated with error: %v", err)
		logrus.Info("self destructing")
		//sending
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}
}

func NewGrpcServer(addr string, coinsData *data.Coins, dataFetcher *fetcher.ClientApi) (*server, error) {
	s := &server{}
	var err error
	s.CoinsData = coinsData
	s.DataFetcher = dataFetcher
	logrus.Infof("Listening on port:%s\n ", addr)
	if s.lis, err = net.Listen("tcp", addr); err != nil {
		return nil, err
	}

	go collector.NewCollector(coinsData, dataFetcher)

	s.server = grpc.NewServer(createServerInterceptors()...)
	pb.RegisterServiceServer(s.server, s)
	return s, nil
}
