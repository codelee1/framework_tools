package main

import (
	"flag"
	"framework_tools/go_kit/v1/v1_endpoint"
	"framework_tools/go_kit/v1/v1_service"
	"framework_tools/go_kit/v1/v1_transport"
	"framework_tools/go_kit/v1/v1_transport/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	fs := flag.NewFlagSet("v1", flag.ExitOnError)
	var (
		httpAddr = fs.String("http-addr", "127.0.0.1:8081", "HTTP listen address")
		grpcAddr = fs.String("grpc-addr", "127.0.0.1:8082", "gRPC listen address")
	)
	var (
		logger, _   = zap.NewDevelopment()
		golangLimit = rate.NewLimiter(10, 2) //每秒产生10个令牌,令牌桶可装2个令牌
		server      = v1_service.NewService(logger)
		endpoints   = v1_endpoint.NewEndPointServer(server, logger, golangLimit)
		grpcServer  = v1_transport.NewGRPCServer(endpoints, logger)
		httpHandler = v1_transport.NewHttpServer(endpoints, logger)
	)

	var g group.Group
	// grpc server
	{
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Warn("transport gRPC Listen err", zap.Error(err))
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Info("transport gRPC addr" + *grpcAddr)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(grpc.UnaryServerInterceptor(grpctransport.Interceptor)))
			pb.RegisterUserServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	// http server
	{
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			logger.Warn("transport HTTP Listen err", zap.Error(err))
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Info("transport HTTP addr" + *httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	logger.Warn("exit", zap.Error(g.Run()))
}
