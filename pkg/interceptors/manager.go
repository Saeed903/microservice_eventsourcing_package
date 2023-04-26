package interceptors

import (
	"context"
	"time"

	"github.com/saeed903/microservice_eventsourcing_package/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GrpcMetricCb func(err error)

type InterceptorManager interface {
	Logger(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error)
	ClientRequestLoggerInterceptor() func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error
}

// InterceptorManager struct
type interceptorManager struct {
	log      logger.Logger
	metricCb GrpcMetricCb
}

// NewInterceptorManager  InterceptorManager constuctor
func NewInterceptorManager(logger logger.Logger, metricCb GrpcMetricCb) *interceptorManager {
	return &interceptorManager{
		log:      logger,
		metricCb: metricCb,
	}
}

// Logger Interceptor
func (im *interceptorManager) Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.log.GrpcMiddlewareAccessLogger(info.FullMethod, time.Since(start), md, err)

	if im.metricCb != nil {
		im.metricCb(err)
	}
	return reply, err
}

// ClientRequestLoggerInterceptor gRPC client interceptor
func (im *interceptorManager) ClientRequestLoggerInterceptor() func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		md, _ := metadata.FromIncomingContext(ctx)
		im.log.GrpcClientInterceptorLogger(method, req, reply, time.Since(start), md, err)
		return err
	}
}
