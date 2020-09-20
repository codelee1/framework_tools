package main

import (
	"context"
	"flag"
	"fmt"
	"framework_tools/go_kit/v5/utils"
	"framework_tools/go_kit/v5/v5_endpoint"
	"framework_tools/go_kit/v5/v5_service"
	"framework_tools/go_kit/v5/v5_transport"
	"framework_tools/go_kit/v5/v5_transport/pb"
	"github.com/go-kit/kit/log"
	metricsprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/sd/etcdv3"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	fs := flag.NewFlagSet("v5", flag.ExitOnError)
	var (
		httpAddr       = fs.String("http-addr", "127.0.0.1:8081", "HTTP listen address")
		grpcAddr       = fs.String("grpc-addr", "127.0.0.1:8082", "gRPC listen address")
		prometheusAddr = fs.String("prometheus-addr", "127.0.0.1:8000", "gRPC listen address")
	)
	tracer, closer, err := utils.NewJaegerTracer("user.service")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		closer.Close()
	}()
	var (
		count = metricsprometheus.NewCounterFrom(prometheus.CounterOpts{
			Subsystem: "user_service",
			Name:      "request_count",
			Help:      "Number of requests",
		}, []string{"method"})

		histogram = metricsprometheus.NewHistogramFrom(prometheus.HistogramOpts{
			Subsystem: "user_service",
			Name:      "request_consume",
			Help:      "Request consumes time",
		}, []string{"method"})

		logger, _   = zap.NewDevelopment()
		golangLimit = rate.NewLimiter(100, 100) //每秒产生10个令牌,令牌桶可装2个令牌（一秒2个）
		server      = v5_service.NewService(logger, count, histogram, tracer)
		endpoints   = v5_endpoint.NewEndPointServer(server, logger, golangLimit, tracer)
		grpcServer  = v5_transport.NewGRPCServer(endpoints, logger)
		httpHandler = v5_transport.NewHttpServer(endpoints, logger)
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
			httpRegistrar.Register()
			logger.Info("httpRegistrar.Register")
			defer func() {
				logger.Error("httpRegistrar.Deregister")
				httpRegistrar.Deregister()
			}()
			logger.Warn("transport HTTP Listen err", zap.Error(err))
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Info("transport HTTP addr:" + *httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	// prometheus
	{
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())

		g.Add(func() error {
			logger.Info("prometheus addr:" + *prometheusAddr)
			return http.ListenAndServe(*prometheusAddr, m)
		}, func(error) {
		})
	}
	logger.Warn("exit", zap.Error(g.Run()))
	logger.Info("run-----------")
}
