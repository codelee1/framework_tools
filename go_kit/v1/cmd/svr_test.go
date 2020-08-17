package main

import (
	"context"
	"framework_tools/go_kit/v1/v1_endpoint"
	"framework_tools/go_kit/v1/v1_transport"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"testing"
)


func TestHttpClient(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	svc, err := v1_transport.NewHttpClient("127.0.0.1:8081", logger)
	if err != nil {
		t.Error(err)
		return
	}
	loginAck, err := svc.Login(context.Background(), "codelee1","123456")
	if err != nil {
		t.Error(err)
		return
	}
	ack := loginAck.(*v1_endpoint.LoginResp)
	t.Logf("loginAck: %+v", ack)

	userInfoAck, err := svc.UserInfo(context.Background(), ack.Token)
	if err != nil {
		t.Error(err)
		return
	}
	userAck := userInfoAck.(*v1_endpoint.UserInfoResp)
	t.Logf("userInfoAck:%+v", userAck)
}

func TestGRPCClient(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	conn, err := grpc.Dial("127.0.0.1:8082", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	svc, err := v1_transport.NewGRPCClient(conn, logger)
	if err != nil {
		t.Error(err)
		return
	}
	loginAck, err := svc.Login(context.Background(), "codelee1","123456")
	if err != nil {
		t.Error(err)
		return
	}
	ack := loginAck.(*v1_endpoint.LoginResp)
	t.Logf("loginAck: %+v", ack)

	userInfoAck, err := svc.UserInfo(context.Background(), ack.Token)
	if err != nil {
		t.Error(err)
		return
	}
	userAck := userInfoAck.(*v1_endpoint.UserInfoResp)
	t.Logf("userInfoAck:%+v", userAck)
}
