package v4_service

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"go.uber.org/zap"
	"time"
)

const ContextReqUUid = "req_uuid"

type NewMiddlewareServer func(Service) Service

type logMiddlewareServer struct {
	logger *zap.Logger
	next   Service
}

func NewLogMiddlewareServer(log *zap.Logger) NewMiddlewareServer {
	return func(service Service) Service {
		return logMiddlewareServer{
			logger: log,
			next:   service,
		}
	}
}

func (l logMiddlewareServer) Login(ctx context.Context, account, password string) (resp interface{}, err error) {
	defer func() {
		l.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Login logMiddlewareServer", "Login"), zap.Any("req", account+":"+password), zap.Any("res", resp), zap.Any("err", err))
	}()
	resp, err = l.next.Login(ctx, account, password)
	return
}

func (l logMiddlewareServer) UserInfo(ctx context.Context, token string) (resp interface{}, err error) {
	defer func() {
		l.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 UserInfo logMiddlewareServer", "UserInfo"), zap.Any("req", token), zap.Any("res", resp), zap.Any("err", err))
	}()
	resp, err = l.next.UserInfo(ctx, token)
	return
}

type MetricsMiddlewareServer struct {
	next      Service
	counter   metrics.Counter
	histogram metrics.Histogram
}

func NewMetricsMiddlewareServer(counter metrics.Counter, histogram metrics.Histogram) NewMiddlewareServer {
	return func(service Service) Service {
		return MetricsMiddlewareServer{
			next:      service,
			counter:   counter,
			histogram: histogram,
		}
	}
}

func (l MetricsMiddlewareServer) Login(ctx context.Context, account, password string) (resp interface{}, err error) {
	defer func(start time.Time) {
		method := []string{"method", "Login"}
		l.counter.With(method...).Add(1)
		l.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	resp, err = l.next.Login(ctx, account, password)
	return
}

func (l MetricsMiddlewareServer) UserInfo(ctx context.Context, token string) (resp interface{}, err error) {
	defer func(start time.Time) {
		method := []string{"method", "UserInfo"}
		l.counter.With(method...).Add(1)
		l.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	resp, err = l.next.UserInfo(ctx, token)
	return
}
