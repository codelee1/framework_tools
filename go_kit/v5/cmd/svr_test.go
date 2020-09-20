package main

import (
	"context"
	"framework_tools/go_kit/v5/utils"
	"framework_tools/go_kit/v5/v5_endpoint"
	"framework_tools/go_kit/v5/v5_service"
	"framework_tools/go_kit/v5/v5_transport"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"math/rand"
	"testing"
	"time"
)

func TestHttpClient(t *testing.T) {

	logger, _ := zap.NewDevelopment()
	svc, err := v5_transport.NewHttpClient("127.0.0.1:8081", logger)
	if err != nil {
		t.Error(err)
		return
	}
	loginAck, err := svc.Login(context.Background(), "codelee1", "123456")
	if err != nil {
		t.Error(err)
		return
	}
	ack := loginAck.(*v5_endpoint.LoginResp)
	t.Logf("loginAck: %+v", ack)

	userInfoAck, err := svc.UserInfo(context.Background(), ack.Token)
	if err != nil {
		t.Error(err)
		return
	}
	userAck := userInfoAck.(*v5_endpoint.UserInfoResp)
	t.Logf("userInfoAck:%+v", userAck)
}

func TestGRPCClient(t *testing.T) {
	var (
		etcdAddrs    = []string{"127.0.0.1:2379"}
		serName      = "svc.service.grpc"
		ttl          = 5 * time.Second
		retryMax     = 3
		retryTimeout = 500 * time.Millisecond
	)
	options := etcdv3.ClientOptions{
		DialTimeout:   ttl,
		DialKeepAlive: ttl,
	}
	etcdClient, err := etcdv3.NewClient(context.Background(), etcdAddrs, options)
	if err != nil {
		t.Error(err)
		return
	}
	instancerm, err := etcdv3.NewInstancer(etcdClient, serName, log.NewNopLogger())
	if err != nil {
		t.Error(err)
		return
	}
	tracer, closer, err := utils.NewJaegerTracer("client.user")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		closer.Close()
	}()
	var (
		endpoints v5_endpoint.Set
		kitlog    = log.NewNopLogger()
	)
	{
		loginFactory := factoryForGrpcService(tracer, v5_endpoint.MakeLoginClientEndPoint)

		endpointer := sd.NewEndpointer(instancerm, loginFactory, kitlog)
		balancer := lb.NewRandom(endpointer, time.Now().UnixNano())
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.LoginEndPoint = retry
	}
	{
		userInfoFactory := factoryForGrpcService(tracer, v5_endpoint.MakeGetUserInfoClientEndPoint)
		endpointer := sd.NewEndpointer(instancerm, userInfoFactory, kitlog)
		balancer := lb.NewRandom(endpointer, time.Now().UnixNano())
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.UserInfoEndPoint = retry
	}

	// 服务熔断
	funcName := "login-userInfo"
	fallback := func(err error) error {
		// todo 执行备用逻辑
		return err
	}
	hy := utils.NewHystrix(fallback)
	c, _, _ := hystrix.GetCircuit(funcName)

	rand.Seed(time.Now().Unix())
	for i := 0; i < 1; i++ {
		time.Sleep(time.Millisecond * 100)
		err := hy.Run(funcName, func() error {
			loginAck, err := endpoints.Login(context.Background(), "codelee1", "123456")
			if err != nil {
				//t.Error("endpoints.Login",err)
				return err
			}

			ack := loginAck.(*v5_endpoint.LoginResp)
			t.Logf("loginAck: %+v", ack)

			userInfoAck, err := endpoints.UserInfo(context.Background(), ack.Token)
			if err != nil {
				//t.Error(err)
				return err
			}
			userAck := userInfoAck.(*v5_endpoint.UserInfoResp)
			t.Logf("userInfoAck:%+v", userAck)
			t.Logf("调用成功--------------------%d,:%+v", i, userAck)
			return err
		})
		t.Log(i, " 熔断器是否开启：", c.IsOpen(), ". 请求是否允许 :", c.AllowRequest())
		if err != nil {
			t.Logf("调用失败--------------------%d,:err:%s", i, err)
		}
	}
}

func factoryForGrpcService(tracer opentracing.Tracer, makeEndpoint func(v5_service.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		logger, _ := zap.NewDevelopment()
		chainUnaryServer := grpcmiddleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_zap.UnaryClientInterceptor(logger),
			//utils.JaegerServerMiddleware(tracer),
		)
		conn, err := grpc.Dial(instance, grpc.WithInsecure(), grpc.WithUnaryInterceptor(chainUnaryServer))
		if err != nil {
			logger.Error("err:", zap.Error(err))
			return nil, nil, err
		}
		//defer conn.Close()
		svc, err := v5_transport.NewGRPCClient(conn, logger)
		if err != nil {
			logger.Error("err:", zap.Error(err))
			return nil, nil, err
		}
		e := makeEndpoint(svc)

		return e, conn, nil
	}
}
