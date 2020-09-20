package v2_endpoint

import (
	"context"
	"fmt"
	"framework_tools/go_kit/v2/v2_service"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Set struct {
	LoginEndPoint    endpoint.Endpoint
	UserInfoEndPoint endpoint.Endpoint
}

func NewEndPointServer(svc v2_service.Service, log *zap.Logger, limit *rate.Limiter) Set {
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

// set实现了service的login，可作为一个service使用 for client
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

func MakeLoginEndPoint(s v2_service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*LoginReq)
		resp, err := s.Login(ctx, req.Account, req.Password)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return &LoginResp{Token: resp.(string), Err: err}, nil
	}
}

// MakeLoginClientEndPoint为客户端创建endpoint
func MakeLoginClientEndPoint(s v2_service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*LoginReq)
		fmt.Println("service", s)
		fmt.Printf("service %+v\n", s)
		resp, err := s.Login(ctx, req.Account, req.Password) // 直接调用endpoint实现的方法，
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func MakeGetUserInfoEndPoint(s v2_service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*UserInfoReq)
		resp, err := s.UserInfo(ctx, req.Token)
		if err != nil {
			fmt.Println(err)
		}
		return &UserInfoResp{Account: resp.(string), Err: err}, nil
	}
}

// MakeGetUserInfoClientEndPoint 为客户端创建endpoint
func MakeGetUserInfoClientEndPoint(s v2_service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*UserInfoReq)
		resp, err := s.UserInfo(ctx, req.Token)
		if err != nil {
			fmt.Println(err)
		}
		return resp, nil
	}
}

type LoginReq struct {
	Account, Password string
}

type LoginResp struct {
	Token string `json:"token"`
	Err   error  `json:"-"`
}

// Failed implements endpoint.Failer.
func (r LoginResp) Failed() error { return r.Err }

type UserInfoReq struct {
	Token string
}

type UserInfoResp struct {
	Account string `json:"account"`
	Err     error  `json:"-"`
}

// Failed implements endpoint.Failer.
func (r UserInfoResp) Failed() error { return r.Err }
