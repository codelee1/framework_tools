package main

import (
	"context"
	"flag"
	"fmt"
	"framework_tools/go_kit/v2/v2_endpoint"
	"framework_tools/go_kit/v2/v2_service"
	"framework_tools/go_kit/v2/v2_transport"
	"framework_tools/go_kit/v2/v2_transport/pb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	fs := flag.NewFlagSet("v2", flag.ExitOnError)
	var (
		httpAddr = fs.String("http-addr", "127.0.0.1:8081", "HTTP listen address")
		grpcAddr = fs.String("grpc-addr", "127.0.0.1:8082", "gRPC listen address")
	)
	var (
		logger, _   = zap.NewDevelopment()
		golangLimit = rate.NewLimiter(10, 2) //每秒产生10个令牌,令牌桶可装2个令牌（一秒2个）
		server      = v2_service.NewService(logger)
		endpoints   = v2_endpoint.NewEndPointServer(server, logger, golangLimit)
		grpcServer  = v2_transport.NewGRPCServer(endpoints, logger)
		httpHandler = v2_transport.NewHttpServer(endpoints, logger)
	)
	var (
		etcdAddrs   = []string{"127.0.0.1:2379"}
		grpcSvrName = "svc.service.grpc"
		httpSvrName = "svc.service.http"
		ttl         = 5 * time.Second
	)
	options := etcdv3.ClientOptions{
		DialTimeout:   ttl,
		DialKeepAlive: ttl,
	}
	etcdClient, err := etcdv3.NewClient(context.Background(), etcdAddrs, options)
	if err != nil {
		logger.Error("new etcd client err ", zap.Error(err))
		return
	}

	/*
		1、初始化etcd客户端，
		2、待个服务启动完成后，通过etcd客户端注册到etcd中
	*/

	var g group.Group
	// grpc server
	{
		grpcRegistrar := etcdv3.NewRegistrar(etcdClient, etcdv3.Service{
			Key:   fmt.Sprintf("%s", grpcSvrName),
			Value: *grpcAddr,
			//TTL:   nil,
		}, log.NewNopLogger())
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Warn("transport gRPC Listen err", zap.Error(err))
			os.Exit(1)
		}
		g.Add(func() error {
			// 注册服务到etcd
			grpcRegistrar.Register()
			logger.Info("grpcRegistrar.Register")
			defer func() {
				logger.Error("grpcRegistrar.Deregister")
				grpcRegistrar.Deregister()
			}()
			logger.Info("transport gRPC addr:" + *grpcAddr)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(grpctransport.Interceptor))
			pb.RegisterUserServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	// http server
	{
		httpRegistrar := etcdv3.NewRegistrar(etcdClient, etcdv3.Service{
			Key:   fmt.Sprintf("%s", httpSvrName),
			Value: *httpAddr,
			TTL:   nil,
		}, log.NewNopLogger())
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			logger.Warn("transport HTTP Listen err", zap.Error(err))
			os.Exit(1)
		}
		g.Add(func() error {
			httpRegistrar.Register()
			logger.Info("httpRegistrar.Register")
			defer func() {
				logger.Error("httpRegistrar.Deregister")
				httpRegistrar.Deregister()
			}()
			logger.Info("transport HTTP addr:" + *httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	logger.Error("run err ", zap.Error(g.Run()))

}
