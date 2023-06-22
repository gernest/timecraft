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
	// TimecraftServiceSubmitTasksProcedure is the fully-qualified name of the TimecraftService's
	// SubmitTasks RPC.
	TimecraftServiceSubmitTasksProcedure = "/timecraft.server.v1.TimecraftService/SubmitTasks"
	// TimecraftServiceLookupTasksProcedure is the fully-qualified name of the TimecraftService's
	// LookupTasks RPC.
	TimecraftServiceLookupTasksProcedure = "/timecraft.server.v1.TimecraftService/LookupTasks"
	// TimecraftServicePollTasksProcedure is the fully-qualified name of the TimecraftService's
	// PollTasks RPC.
	TimecraftServicePollTasksProcedure = "/timecraft.server.v1.TimecraftService/PollTasks"
	// TimecraftServiceDiscardTasksProcedure is the fully-qualified name of the TimecraftService's
	// DiscardTasks RPC.
	TimecraftServiceDiscardTasksProcedure = "/timecraft.server.v1.TimecraftService/DiscardTasks"
	// TimecraftServiceVersionProcedure is the fully-qualified name of the TimecraftService's Version
	// RPC.
	TimecraftServiceVersionProcedure = "/timecraft.server.v1.TimecraftService/Version"
)

// TimecraftServiceClient is a client for the timecraft.server.v1.TimecraftService service.
type TimecraftServiceClient interface {
	SubmitTasks(context.Context, *connect_go.Request[v1.SubmitTasksRequest]) (*connect_go.Response[v1.SubmitTasksResponse], error)
	LookupTasks(context.Context, *connect_go.Request[v1.LookupTasksRequest]) (*connect_go.Response[v1.LookupTasksResponse], error)
	PollTasks(context.Context, *connect_go.Request[v1.PollTasksRequest]) (*connect_go.Response[v1.PollTasksResponse], error)
	DiscardTasks(context.Context, *connect_go.Request[v1.DiscardTasksRequest]) (*connect_go.Response[v1.DiscardTasksResponse], error)
	Version(context.Context, *connect_go.Request[v1.VersionRequest]) (*connect_go.Response[v1.VersionResponse], error)
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
		submitTasks: connect_go.NewClient[v1.SubmitTasksRequest, v1.SubmitTasksResponse](
			httpClient,
			baseURL+TimecraftServiceSubmitTasksProcedure,
			opts...,
		),
		lookupTasks: connect_go.NewClient[v1.LookupTasksRequest, v1.LookupTasksResponse](
			httpClient,
			baseURL+TimecraftServiceLookupTasksProcedure,
			opts...,
		),
		pollTasks: connect_go.NewClient[v1.PollTasksRequest, v1.PollTasksResponse](
			httpClient,
			baseURL+TimecraftServicePollTasksProcedure,
			opts...,
		),
		discardTasks: connect_go.NewClient[v1.DiscardTasksRequest, v1.DiscardTasksResponse](
			httpClient,
			baseURL+TimecraftServiceDiscardTasksProcedure,
			opts...,
		),
		version: connect_go.NewClient[v1.VersionRequest, v1.VersionResponse](
			httpClient,
			baseURL+TimecraftServiceVersionProcedure,
			opts...,
		),
	}
}

// timecraftServiceClient implements TimecraftServiceClient.
type timecraftServiceClient struct {
	submitTasks  *connect_go.Client[v1.SubmitTasksRequest, v1.SubmitTasksResponse]
	lookupTasks  *connect_go.Client[v1.LookupTasksRequest, v1.LookupTasksResponse]
	pollTasks    *connect_go.Client[v1.PollTasksRequest, v1.PollTasksResponse]
	discardTasks *connect_go.Client[v1.DiscardTasksRequest, v1.DiscardTasksResponse]
	version      *connect_go.Client[v1.VersionRequest, v1.VersionResponse]
}

// SubmitTasks calls timecraft.server.v1.TimecraftService.SubmitTasks.
func (c *timecraftServiceClient) SubmitTasks(ctx context.Context, req *connect_go.Request[v1.SubmitTasksRequest]) (*connect_go.Response[v1.SubmitTasksResponse], error) {
	return c.submitTasks.CallUnary(ctx, req)
}

// LookupTasks calls timecraft.server.v1.TimecraftService.LookupTasks.
func (c *timecraftServiceClient) LookupTasks(ctx context.Context, req *connect_go.Request[v1.LookupTasksRequest]) (*connect_go.Response[v1.LookupTasksResponse], error) {
	return c.lookupTasks.CallUnary(ctx, req)
}

// PollTasks calls timecraft.server.v1.TimecraftService.PollTasks.
func (c *timecraftServiceClient) PollTasks(ctx context.Context, req *connect_go.Request[v1.PollTasksRequest]) (*connect_go.Response[v1.PollTasksResponse], error) {
	return c.pollTasks.CallUnary(ctx, req)
}

// DiscardTasks calls timecraft.server.v1.TimecraftService.DiscardTasks.
func (c *timecraftServiceClient) DiscardTasks(ctx context.Context, req *connect_go.Request[v1.DiscardTasksRequest]) (*connect_go.Response[v1.DiscardTasksResponse], error) {
	return c.discardTasks.CallUnary(ctx, req)
}

// Version calls timecraft.server.v1.TimecraftService.Version.
func (c *timecraftServiceClient) Version(ctx context.Context, req *connect_go.Request[v1.VersionRequest]) (*connect_go.Response[v1.VersionResponse], error) {
	return c.version.CallUnary(ctx, req)
}

// TimecraftServiceHandler is an implementation of the timecraft.server.v1.TimecraftService service.
type TimecraftServiceHandler interface {
	SubmitTasks(context.Context, *connect_go.Request[v1.SubmitTasksRequest]) (*connect_go.Response[v1.SubmitTasksResponse], error)
	LookupTasks(context.Context, *connect_go.Request[v1.LookupTasksRequest]) (*connect_go.Response[v1.LookupTasksResponse], error)
	PollTasks(context.Context, *connect_go.Request[v1.PollTasksRequest]) (*connect_go.Response[v1.PollTasksResponse], error)
	DiscardTasks(context.Context, *connect_go.Request[v1.DiscardTasksRequest]) (*connect_go.Response[v1.DiscardTasksResponse], error)
	Version(context.Context, *connect_go.Request[v1.VersionRequest]) (*connect_go.Response[v1.VersionResponse], error)
}

// NewTimecraftServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTimecraftServiceHandler(svc TimecraftServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(TimecraftServiceSubmitTasksProcedure, connect_go.NewUnaryHandler(
		TimecraftServiceSubmitTasksProcedure,
		svc.SubmitTasks,
		opts...,
	))
	mux.Handle(TimecraftServiceLookupTasksProcedure, connect_go.NewUnaryHandler(
		TimecraftServiceLookupTasksProcedure,
		svc.LookupTasks,
		opts...,
	))
	mux.Handle(TimecraftServicePollTasksProcedure, connect_go.NewUnaryHandler(
		TimecraftServicePollTasksProcedure,
		svc.PollTasks,
		opts...,
	))
	mux.Handle(TimecraftServiceDiscardTasksProcedure, connect_go.NewUnaryHandler(
		TimecraftServiceDiscardTasksProcedure,
		svc.DiscardTasks,
		opts...,
	))
	mux.Handle(TimecraftServiceVersionProcedure, connect_go.NewUnaryHandler(
		TimecraftServiceVersionProcedure,
		svc.Version,
		opts...,
	))
	return "/timecraft.server.v1.TimecraftService/", mux
}

// UnimplementedTimecraftServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTimecraftServiceHandler struct{}

func (UnimplementedTimecraftServiceHandler) SubmitTasks(context.Context, *connect_go.Request[v1.SubmitTasksRequest]) (*connect_go.Response[v1.SubmitTasksResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("timecraft.server.v1.TimecraftService.SubmitTasks is not implemented"))
}

func (UnimplementedTimecraftServiceHandler) LookupTasks(context.Context, *connect_go.Request[v1.LookupTasksRequest]) (*connect_go.Response[v1.LookupTasksResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("timecraft.server.v1.TimecraftService.LookupTasks is not implemented"))
}

func (UnimplementedTimecraftServiceHandler) PollTasks(context.Context, *connect_go.Request[v1.PollTasksRequest]) (*connect_go.Response[v1.PollTasksResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("timecraft.server.v1.TimecraftService.PollTasks is not implemented"))
}

func (UnimplementedTimecraftServiceHandler) DiscardTasks(context.Context, *connect_go.Request[v1.DiscardTasksRequest]) (*connect_go.Response[v1.DiscardTasksResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("timecraft.server.v1.TimecraftService.DiscardTasks is not implemented"))
}

func (UnimplementedTimecraftServiceHandler) Version(context.Context, *connect_go.Request[v1.VersionRequest]) (*connect_go.Response[v1.VersionResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("timecraft.server.v1.TimecraftService.Version is not implemented"))
}