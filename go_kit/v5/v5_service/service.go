package v5_service

import (
	"context"
	"errors"
	"fmt"
	"framework_tools/go_kit/v5/utils"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type Service interface {
	Login(ctx context.Context, account, password string) (resp interface{}, err error)
	UserInfo(ctx context.Context, token string) (resp interface{}, err error)
}

type baseServer struct {
	logger *zap.Logger
}

func NewService(log *zap.Logger, counter *prometheus.Counter, histogram *prometheus.Histogram, tracer opentracing.Tracer) Service {
	var server Service
	server = &baseServer{log}
	server = NewTraceMiddlewareServer(tracer)(server)
	server = NewLogMiddlewareServer(log)(server)
	server = NewMetricsMiddlewareServer(counter, histogram)(server)
	return server
}

func (s baseServer) Login(ctx context.Context, account, password string) (resp interface{}, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 v5_service Service", "Login 处理请求"))
	if account != "codelee1" || password != "123456" {
		err = errors.New("用户信息错误")
		return nil, err
	}
	token, err := utils.CreateJwtToken(account, 1)
	if err != nil {
		err = errors.New("创建token失败")
		return
	}
	// 测试超时
	rand.Seed(time.Now().UnixNano())
	//sl := rand.Int31n(10-1) + 1
	//time.Sleep(time.Millisecond * 100 * time.Duration(sl))
	// 测试错误率
	//r := rand.Int31n(10)
	//if r < 4 {
	//	return nil,errors.New("模拟服务端错误:Login")
	//}
	resp = token
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 v5_service Service", "Login 处理请求"), zap.Any("处理返回值", resp))
	return
}

func (s baseServer) UserInfo(ctx context.Context, token string) (resp interface{}, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 v5_service Service", "UserInfo 处理请求"))
	if token == "" {
		err = errors.New("token为空")
		return
	}
	jwtInfo, err := utils.ParseToken(token)
	if err != nil {
		err = errors.New("验证用户信息错误")
		return
	}
	// 测试超时
	rand.Seed(time.Now().UnixNano())
	//sl := rand.Int31n(10-1) + 1
	//time.Sleep(time.Millisecond * 100 * time.Duration(sl))
	// 测试错误率
	//r := rand.Int31n(10)
	//if r < 4 {
	//	return nil,errors.New("模拟服务端错误:UserInfo")
	//}
	resp = fmt.Sprintf("%s", jwtInfo["Name"])
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 v5_service Service", "UserInfo 处理请求"), zap.Any("处理返回值", resp))
	return
}
