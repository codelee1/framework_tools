package v5_service

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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

func (m MetricsMiddlewareServer) Login(ctx context.Context, account, password string) (resp interface{}, err error) {
	defer func(start time.Time) {
		method := []string{"method", "Login"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	resp, err = m.next.Login(ctx, account, password)
	return
}

func (m MetricsMiddlewareServer) UserInfo(ctx context.Context, token string) (resp interface{}, err error) {
	defer func(start time.Time) {
		method := []string{"method", "UserInfo"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	resp, err = m.next.UserInfo(ctx, token)
	return
}

type TracerMiddlewareServer struct {
	next   Service
	tracer opentracing.Tracer
}

func NewTraceMiddlewareServer(trace opentracing.Tracer) NewMiddlewareServer {
	return func(service Service) Service {
		return TracerMiddlewareServer{
			tracer: trace,
			next:   service,
		}
	}
}

func (t TracerMiddlewareServer) Login(ctx context.Context, account, password string) (resp interface{}, err error) {
	span, ctxContext := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "service.login", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "NewTracerLoginMiddleware",
	})
	defer func() {
		span.LogKV("account", account, "password", password)
		span.Finish()
	}()
	resp, err = t.next.Login(ctxContext, account, password)
	return
}

func (t TracerMiddlewareServer) UserInfo(ctx context.Context, token string) (resp interface{}, err error) {
	span, ctxContext := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "service.user_info", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "NewTracerUserInfoMiddleware",
	})
	defer func() {
		span.LogKV("token", token)
		span.Finish()
	}()
	resp, err = t.next.UserInfo(ctxContext, token)
	return
}
