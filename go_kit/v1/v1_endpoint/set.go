package v1_endpoint

import (
	"context"
	"framework_tools/go_kit/v1/v1_service"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// set实现了service的login，可作为一个service使用 for client
type Set struct {
	LoginEndPoint    endpoint.Endpoint
	UserInfoEndPoint endpoint.Endpoint
}

func NewEndPointServer(svc v1_service.Service, log *zap.Logger, limit *rate.Limiter) Set {
	var loginEndPoint endpoint.Endpoint
	{
		loginEndPoint = MakeLoginEndPoint(svc)
		loginEndPoint = LoggingMiddleware(log, "Login")(loginEndPoint)
		loginEndPoint = NewGolangRateAllowMiddleware(limit)(loginEndPoint)
	}

	var userInfoEndPoint endpoint.Endpoint
	{
		userInfoEndPoint = MakeGetUserInfoEndPoint(svc)
		userInfoEndPoint = LoggingMiddleware(log, "UserInfo")(userInfoEndPoint)
		userInfoEndPoint = NewGolangRateAllowMiddleware(limit)(userInfoEndPoint)
	}
	return Set{
		LoginEndPoint:    loginEndPoint,
		UserInfoEndPoint: userInfoEndPoint,
	}
}

func (s Set) Login(ctx context.Context, account, password string) (interface{}, error) {
	resp, err := s.LoginEndPoint(ctx, &LoginReq{Account: account, Password: password})
	if err != nil {
		return nil, err
	}
	response := resp.(*LoginResp)
	return response, response.Err
}

func (s Set) UserInfo(ctx context.Context, token string) (interface{}, error) {
	resp, err := s.UserInfoEndPoint(ctx, &UserInfoReq{Token: token})
	if err != nil {
		return nil, err
	}
	response := resp.(*UserInfoResp)
	return response, response.Err
}

func MakeLoginEndPoint(s v1_service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*LoginReq)
		resp, err := s.Login(ctx, req.Account, req.Password)
		return &LoginResp{Token: resp.(string), Err: err}, nil
	}
}

func MakeGetUserInfoEndPoint(s v1_service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*UserInfoReq)
		resp, err := s.UserInfo(ctx, req.Token)
		return &UserInfoResp{Account: resp.(string), Err: err}, nil
	}
}

type LoginReq struct {
	Account, Password string
}

// 实现 endpoint.Failer.
type LoginResp struct {
	Token string `json:"token"`
	Err   error  `json:"-"`
}

func (r LoginResp) Failed() error { return r.Err }

type UserInfoReq struct {
	Token string
}

// 实现  endpoint.Failer.
type UserInfoResp struct {
	Account string `json:"account"`
	Err     error  `json:"-"`
}

func (r UserInfoResp) Failed() error { return r.Err }
