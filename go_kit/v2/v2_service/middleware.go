package v2_service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
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
