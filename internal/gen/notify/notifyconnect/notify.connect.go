// Code generated by protoc-gen-connect-go.exe. DO NOT EDIT.
//
// Source: notify/notify.proto

package notifyconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	notify "github.com/grpc-buf/internal/gen/notify"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// NotificationServiceName is the fully-qualified name of the NotificationService service.
	NotificationServiceName = "notify.NotificationService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// NotificationServiceCreateNotificationProcedure is the fully-qualified name of the
	// NotificationService's CreateNotification RPC.
	NotificationServiceCreateNotificationProcedure = "/notify.NotificationService/CreateNotification"
	// NotificationServiceGetNotificationsProcedure is the fully-qualified name of the
	// NotificationService's GetNotifications RPC.
	NotificationServiceGetNotificationsProcedure = "/notify.NotificationService/GetNotifications"
)

// NotificationServiceClient is a client for the notify.NotificationService service.
type NotificationServiceClient interface {
	// CreateNotification creates a new notification.
	CreateNotification(context.Context, *connect.Request[notify.CreateNotificationRequest]) (*connect.Response[notify.CreateNotificationResponse], error)
	// GetNotifications gets a list of all notifications for a specific user.
	GetNotifications(context.Context, *connect.Request[notify.GetNotificationsRequest]) (*connect.Response[notify.GetNotificationsResponse], error)
}

// NewNotificationServiceClient constructs a client for the notify.NotificationService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewNotificationServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) NotificationServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &notificationServiceClient{
		createNotification: connect.NewClient[notify.CreateNotificationRequest, notify.CreateNotificationResponse](
			httpClient,
			baseURL+NotificationServiceCreateNotificationProcedure,
			opts...,
		),
		getNotifications: connect.NewClient[notify.GetNotificationsRequest, notify.GetNotificationsResponse](
			httpClient,
			baseURL+NotificationServiceGetNotificationsProcedure,
			opts...,
		),
	}
}

// notificationServiceClient implements NotificationServiceClient.
type notificationServiceClient struct {
	createNotification *connect.Client[notify.CreateNotificationRequest, notify.CreateNotificationResponse]
	getNotifications   *connect.Client[notify.GetNotificationsRequest, notify.GetNotificationsResponse]
}

// CreateNotification calls notify.NotificationService.CreateNotification.
func (c *notificationServiceClient) CreateNotification(ctx context.Context, req *connect.Request[notify.CreateNotificationRequest]) (*connect.Response[notify.CreateNotificationResponse], error) {
	return c.createNotification.CallUnary(ctx, req)
}

// GetNotifications calls notify.NotificationService.GetNotifications.
func (c *notificationServiceClient) GetNotifications(ctx context.Context, req *connect.Request[notify.GetNotificationsRequest]) (*connect.Response[notify.GetNotificationsResponse], error) {
	return c.getNotifications.CallUnary(ctx, req)
}

// NotificationServiceHandler is an implementation of the notify.NotificationService service.
type NotificationServiceHandler interface {
	// CreateNotification creates a new notification.
	CreateNotification(context.Context, *connect.Request[notify.CreateNotificationRequest]) (*connect.Response[notify.CreateNotificationResponse], error)
	// GetNotifications gets a list of all notifications for a specific user.
	GetNotifications(context.Context, *connect.Request[notify.GetNotificationsRequest]) (*connect.Response[notify.GetNotificationsResponse], error)
}

// NewNotificationServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewNotificationServiceHandler(svc NotificationServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	notificationServiceCreateNotificationHandler := connect.NewUnaryHandler(
		NotificationServiceCreateNotificationProcedure,
		svc.CreateNotification,
		opts...,
	)
	notificationServiceGetNotificationsHandler := connect.NewUnaryHandler(
		NotificationServiceGetNotificationsProcedure,
		svc.GetNotifications,
		opts...,
	)
	return "/notify.NotificationService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case NotificationServiceCreateNotificationProcedure:
			notificationServiceCreateNotificationHandler.ServeHTTP(w, r)
		case NotificationServiceGetNotificationsProcedure:
			notificationServiceGetNotificationsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedNotificationServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedNotificationServiceHandler struct{}

func (UnimplementedNotificationServiceHandler) CreateNotification(context.Context, *connect.Request[notify.CreateNotificationRequest]) (*connect.Response[notify.CreateNotificationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("notify.NotificationService.CreateNotification is not implemented"))
}

func (UnimplementedNotificationServiceHandler) GetNotifications(context.Context, *connect.Request[notify.GetNotificationsRequest]) (*connect.Response[notify.GetNotificationsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("notify.NotificationService.GetNotifications is not implemented"))
}
