package v2_service

import (
	"context"
	"errors"
	"fmt"
	"framework_tools/go_kit/v2/utils"
	"go.uber.org/zap"
)

type Service interface {
	Login(ctx context.Context, account, password string) (resp interface{}, err error)
	UserInfo(ctx context.Context, token string) (resp interface{}, err error)
}

type baseServer struct {
	logger *zap.Logger
}

func NewService(log *zap.Logger) Service {
	var server Service
	server = &baseServer{log}
	server = NewLogMiddlewareServer(log)(server)
	return server
}

func (s baseServer) Login(ctx context.Context, account, password string) (resp interface{}, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 v2_service Service", "Login 处理请求"))
	if account != "codelee1" || password != "123456" {
		err = errors.New("用户信息错误")
		return nil, err
	}
	token, err := utils.CreateJwtToken(account, 1)
	if err != nil {
		err = errors.New("创建token失败")
		return
	}
	resp = token
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 v2_service Service", "Login 处理请求"), zap.Any("处理返回值", resp))
	return
}

func (s baseServer) UserInfo(ctx context.Context, token string) (resp interface{}, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 v2_service Service", "UserInfo 处理请求"))
	if token == "" {
		err = errors.New("token为空")
		return
	}
	jwtInfo, err := utils.ParseToken(token)
	if err != nil {
		err = errors.New("验证用户信息错误")
		return
	}
	resp = fmt.Sprintf("%s", jwtInfo["Name"])
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 v2_service Service", "UserInfo 处理请求"), zap.Any("处理返回值", resp))
	return
}
