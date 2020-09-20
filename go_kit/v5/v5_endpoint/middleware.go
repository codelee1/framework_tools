package v5_endpoint

import (
	"context"
	"errors"
	"fmt"
	"framework_tools/go_kit/v5/v5_service"
	"github.com/go-kit/kit/endpoint"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

func LoggingMiddleware(logger *zap.Logger, method string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Debug(fmt.Sprint(ctx.Value(v5_service.ContextReqUUid)), zap.Any("调用 v5_endpoint LoggingMiddleware :"+method, "完成请求处理"), zap.Any("耗时毫秒", time.Since(begin).Milliseconds()))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func NewGolangRateAllowMiddleware(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, errors.New("limit request not Allow")
			}
			return next(ctx, request)
		}
	}
}

func NewTracerMiddleware(method string, tracer opentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			span, ctxContext := opentracing.StartSpanFromContextWithTracer(ctx, tracer, "endpoint", opentracing.Tag{
				Key:   string(ext.Component),
				Value: "NewTracerEndpointMiddleware." + method,
			})
			defer span.Finish()
			return next(ctxContext, request)
		}
	}
}
