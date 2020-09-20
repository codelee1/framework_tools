package v4_endpoint

import (
	"context"
	"errors"
	"fmt"
	"framework_tools/go_kit/v4/v4_service"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

func LoggingMiddleware(logger *zap.Logger, method string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Debug(fmt.Sprint(ctx.Value(v4_service.ContextReqUUid)), zap.Any("调用 v4_endpoint LoggingMiddleware :"+method, "完成请求处理"), zap.Any("耗时毫秒", time.Since(begin).Milliseconds()))
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
