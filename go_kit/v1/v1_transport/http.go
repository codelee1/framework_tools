package v1_transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"framework_tools/go_kit/v1/utils"
	"framework_tools/go_kit/v1/v1_endpoint"
	"framework_tools/go_kit/v1/v1_service"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// http传输编解码

func NewHttpServer(endpoint v1_endpoint.Set, log *zap.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			ctx = context.WithValue(ctx, v1_service.ContextReqUUid, request.Header.Get(v1_service.ContextReqUUid))
			ctx = context.WithValue(ctx, utils.JWT_CONTEXT_KEY, request.Header.Get("token"))

			return ctx
		}),
		httptransport.ServerErrorHandler(NewZapLogErrorHandler(log)),
	}
	m := http.DefaultServeMux
	m.Handle("/login", httptransport.NewServer(
		endpoint.LoginEndPoint,
		decodeHttpLoginRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Handle("/info", httptransport.NewServer(
		endpoint.UserInfoEndPoint,
		decodeHttpUserInfoRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	return m
}

func NewHttpClient(addr string, log *zap.Logger) (v1_service.Service, error) {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}
	u, err := url.Parse(addr)
	if err != nil {
		return v1_endpoint.Set{}, err
	}
	options := []httptransport.ClientOption{
		httptransport.ClientBefore(func(ctx context.Context, r *http.Request) context.Context {
			UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
			log.Debug("添加uuid", zap.Any("UUID", UUID))
			r.Header.Set(v1_service.ContextReqUUid, UUID)
			return context.Background()
		}),
	}
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = httptransport.NewClient(
			"Post",
			copyURL(u, "/login"),
			encodeHTTPGenericRequest,
			decodeHttpLoginResponse,
			options...,
		).Endpoint()
	}
	var userInfoEndpoint endpoint.Endpoint
	{
		userInfoEndpoint = httptransport.NewClient(
			"Post",
			copyURL(u, "/info"),
			encodeHTTPGenericRequest,
			decodeHttpUserInfoResponse,
			options...,
		).Endpoint()
	}

	return v1_endpoint.Set{
		LoginEndPoint:    loginEndpoint,
		UserInfoEndPoint: userInfoEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

func decodeHttpLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req v1_endpoint.LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeHttpUserInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req v1_endpoint.UserInfoReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeHttpLoginResponse(_ context.Context, r *http.Response) (response interface{}, err error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp v1_endpoint.LoginResp
	err = json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

func decodeHttpUserInfoResponse(_ context.Context, r *http.Response) (response interface{}, err error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp v1_endpoint.UserInfoResp
	err = json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

// encodeHTTPGenericRequest http请求json编码 for client
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		return err
	}
	if req, ok := request.(v1_endpoint.UserInfoReq);ok {
		r.Header.Set("token",req.Token)
	}

	r.Body = ioutil.NopCloser(&buffer)
	return nil
}

// encodeHTTPGenericResponse 通用http回包处理 for server
func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}
