package v5_transport

import (
	"context"
	"framework_tools/go_kit/v5/v5_endpoint"
	"framework_tools/go_kit/v5/v5_service"
	"framework_tools/go_kit/v5/v5_transport/pb"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// grpc传输编解码

type grpcServer struct {
	loginServer    grpctransport.Handler
	userInfoServer grpctransport.Handler
}

func NewGRPCServer(endpoint v5_endpoint.Set, log *zap.Logger) pb.UserServer {
	//serverOps := grpctransport.ServerRequestFunc
	options := []grpctransport.ServerOption{
		grpctransport.ServerBefore(func(ctx context.Context, md metadata.MD) context.Context {
			ctx = context.WithValue(ctx, v5_service.ContextReqUUid, md.Get(v5_service.ContextReqUUid))
			return ctx
		}),

		grpctransport.ServerErrorHandler(NewZapLogErrorHandler(log)),
	}
	login := grpctransport.NewServer(
		endpoint.LoginEndPoint,
		decodeGRPCLoginRequest,
		encodeGRPCLoginResponse,
		options...,
	)
	getUserInfo := grpctransport.NewServer(
		endpoint.UserInfoEndPoint,
		decodeGRPCUserInfoRequest,
		encodeGRPCUserInfoResponse,
		options...,
	)

	return &grpcServer{
		loginServer:    login,
		userInfoServer: getUserInfo,
	}
}

func (s *grpcServer) Login(ctx context.Context, req *pb.LoginREQ) (*pb.LoginACK, error) {
	_, rep, err := s.loginServer.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LoginACK), nil
}

func (s *grpcServer) UserInfo(ctx context.Context, req *pb.UserInfoREQ) (*pb.UserInfoACK, error) {
	_, rep, err := s.userInfoServer.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserInfoACK), nil
}

func NewGRPCClient(conn *grpc.ClientConn, log *zap.Logger) (v5_service.Service, error) {
	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(func(ctx context.Context, md *metadata.MD) context.Context {
			UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
			//log.Debug("添加uuid", zap.Any("UUID", UUID))
			md.Set(v5_service.ContextReqUUid, UUID)
			ctx = metadata.NewOutgoingContext(context.Background(), *md)
			return ctx
		}),
	}
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"Login",
			encodeGRPCLoginRequest,
			decodeGRPCLoginResponse,
			pb.LoginACK{},
			options...).Endpoint()
	}
	var getUserInfoEndpoint endpoint.Endpoint
	{
		getUserInfoEndpoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"UserInfo",
			encodeGRPCUserInfoRequest,
			decodeGRPCUserInfoResponse,
			pb.UserInfoACK{},
			options...).Endpoint()
	}
	return v5_endpoint.Set{
		LoginEndPoint:    loginEndpoint,
		UserInfoEndPoint: getUserInfoEndpoint,
	}, nil
}

// server

func decodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.LoginREQ)
	return &v5_endpoint.LoginReq{Account: req.GetAccount(), Password: req.GetPassword()}, nil
}

func encodeGRPCLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*v5_endpoint.LoginResp)
	return &pb.LoginACK{Token: resp.Token}, resp.Err
}

func decodeGRPCUserInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserInfoREQ)
	return &v5_endpoint.UserInfoReq{Token: req.GetToken()}, nil
}

func encodeGRPCUserInfoResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*v5_endpoint.UserInfoResp)
	return &pb.UserInfoACK{Account: resp.Account}, resp.Err
}

// client

func encodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*v5_endpoint.LoginReq)
	return &pb.LoginREQ{Account: req.Account, Password: req.Password}, nil
}

func decodeGRPCLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.LoginACK)
	return &v5_endpoint.LoginResp{Token: resp.GetToken()}, nil
}

func encodeGRPCUserInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*v5_endpoint.UserInfoReq)
	return &pb.UserInfoREQ{Token: req.Token}, nil
}

func decodeGRPCUserInfoResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserInfoACK)
	return &v5_endpoint.UserInfoResp{Account: resp.GetAccount()}, nil
}
