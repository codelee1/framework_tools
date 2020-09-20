package utils

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"io"
)

func NewJaegerTracer(serviceName string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	cfg := &jaegerConfig.Configuration{
		ServiceName: serviceName,
		Disabled:    false,
		RPCMetrics:  false,
		Tags:        nil,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:                     "const", //固定采样
			Param:                    1,       //1=全采样、0=不采样
			SamplingServerURL:        "",
			SamplingRefreshInterval:  0,
			MaxOperations:            0,
			OperationNameLateBinding: false,
			Options:                  nil,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			QueueSize:                  0,
			BufferFlushInterval:        0,
			LogSpans:                   true,
			LocalAgentHostPort:         "127.0.0.1:6831",
			DisableAttemptReconnecting: false,
			AttemptReconnectInterval:   0,
			CollectorEndpoint:          "",
			User:                       "",
			Password:                   "",
			HTTPHeaders:                nil,
		},
		Headers:             nil,
		BaggageRestrictions: nil,
		Throttler:           nil,
	}

	tracer, closer, err = cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		return
	}
	opentracing.SetGlobalTracer(tracer)
	return
}
