// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: timecraft/server/v1/timecraft.proto

package serverv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/stealthrocket/timecraft/gen/proto/go/timecraft/server/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// TimecraftServiceName is the fully-qualified name of the TimecraftService service.
	TimecraftServiceName = "timecraft.server.v1.TimecraftService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// TimecraftServiceSubmitTaskProcedure is the fully-qualified name of the TimecraftService's
	// SubmitTask RPC.
	TimecraftServiceSubmitTaskProcedure = "/timecraft.server.v1.TimecraftService/SubmitTask"
)

// TimecraftServiceClient is a client for the timecraft.server.v1.TimecraftService service.
type TimecraftServiceClient interface {
	SubmitTask(context.Context, *connect_go.Request[v1.SubmitTaskRequest]) (*connect_go.Response[v1.SubmitTaskResponse], error)
}

// NewTimecraftServiceClient constructs a client for the timecraft.server.v1.TimecraftService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTimecraftServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) TimecraftServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &timecraftServiceClient{
		submitTask: connect_go.NewClient[v1.SubmitTaskRequest, v1.SubmitTaskResponse](
			httpClient,
			baseURL+TimecraftServiceSubmitTaskProcedure,
			opts...,
		),
	}
}

// timecraftServiceClient implements TimecraftServiceClient.
type timecraftServiceClient struct {
	submitTask *connect_go.Client[v1.SubmitTaskRequest, v1.SubmitTaskResponse]
}

// SubmitTask calls timecraft.server.v1.TimecraftService.SubmitTask.
func (c *timecraftServiceClient) SubmitTask(ctx context.Context, req *connect_go.Request[v1.SubmitTaskRequest]) (*connect_go.Response[v1.SubmitTaskResponse], error) {
	return c.submitTask.CallUnary(ctx, req)
}

// TimecraftServiceHandler is an implementation of the timecraft.server.v1.TimecraftService service.
type TimecraftServiceHandler interface {
	SubmitTask(context.Context, *connect_go.Request[v1.SubmitTaskRequest]) (*connect_go.Response[v1.SubmitTaskResponse], error)
}

// NewTimecraftServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTimecraftServiceHandler(svc TimecraftServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(TimecraftServiceSubmitTaskProcedure, connect_go.NewUnaryHandler(
		TimecraftServiceSubmitTaskProcedure,
		svc.SubmitTask,
		opts...,
	))
	return "/timecraft.server.v1.TimecraftService/", mux
}

// UnimplementedTimecraftServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTimecraftServiceHandler struct{}

func (UnimplementedTimecraftServiceHandler) SubmitTask(context.Context, *connect_go.Request[v1.SubmitTaskRequest]) (*connect_go.Response[v1.SubmitTaskResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("timecraft.server.v1.TimecraftService.SubmitTask is not implemented"))
}
