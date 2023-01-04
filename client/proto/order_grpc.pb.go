// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.11
// source: proto/order.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OrderServiceClient is the client API for OrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderServiceClient interface {
	GreetCustomer(ctx context.Context, in *GreetRequest, opts ...grpc.CallOption) (*GreetResponse, error)
	CreateOrderByIngredients(ctx context.Context, opts ...grpc.CallOption) (OrderService_CreateOrderByIngredientsClient, error)
	GenerateRecommendedIngredientsCombo(ctx context.Context, in *EmptyParam, opts ...grpc.CallOption) (OrderService_GenerateRecommendedIngredientsComboClient, error)
	ProcessOrder(ctx context.Context, opts ...grpc.CallOption) (OrderService_ProcessOrderClient, error)
}

type orderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return &orderServiceClient{cc}
}

func (c *orderServiceClient) GreetCustomer(ctx context.Context, in *GreetRequest, opts ...grpc.CallOption) (*GreetResponse, error) {
	out := new(GreetResponse)
	err := c.cc.Invoke(ctx, "/order_proto.OrderService/GreetCustomer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) CreateOrderByIngredients(ctx context.Context, opts ...grpc.CallOption) (OrderService_CreateOrderByIngredientsClient, error) {
	stream, err := c.cc.NewStream(ctx, &OrderService_ServiceDesc.Streams[0], "/order_proto.OrderService/CreateOrderByIngredients", opts...)
	if err != nil {
		return nil, err
	}
	x := &orderServiceCreateOrderByIngredientsClient{stream}
	return x, nil
}

type OrderService_CreateOrderByIngredientsClient interface {
	Send(*Ingredient) error
	CloseAndRecv() (*OrderMessage, error)
	grpc.ClientStream
}

type orderServiceCreateOrderByIngredientsClient struct {
	grpc.ClientStream
}

func (x *orderServiceCreateOrderByIngredientsClient) Send(m *Ingredient) error {
	return x.ClientStream.SendMsg(m)
}

func (x *orderServiceCreateOrderByIngredientsClient) CloseAndRecv() (*OrderMessage, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(OrderMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *orderServiceClient) GenerateRecommendedIngredientsCombo(ctx context.Context, in *EmptyParam, opts ...grpc.CallOption) (OrderService_GenerateRecommendedIngredientsComboClient, error) {
	stream, err := c.cc.NewStream(ctx, &OrderService_ServiceDesc.Streams[1], "/order_proto.OrderService/GenerateRecommendedIngredientsCombo", opts...)
	if err != nil {
		return nil, err
	}
	x := &orderServiceGenerateRecommendedIngredientsComboClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OrderService_GenerateRecommendedIngredientsComboClient interface {
	Recv() (*Ingredient, error)
	grpc.ClientStream
}

type orderServiceGenerateRecommendedIngredientsComboClient struct {
	grpc.ClientStream
}

func (x *orderServiceGenerateRecommendedIngredientsComboClient) Recv() (*Ingredient, error) {
	m := new(Ingredient)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *orderServiceClient) ProcessOrder(ctx context.Context, opts ...grpc.CallOption) (OrderService_ProcessOrderClient, error) {
	stream, err := c.cc.NewStream(ctx, &OrderService_ServiceDesc.Streams[2], "/order_proto.OrderService/ProcessOrder", opts...)
	if err != nil {
		return nil, err
	}
	x := &orderServiceProcessOrderClient{stream}
	return x, nil
}

type OrderService_ProcessOrderClient interface {
	Send(*OrderId) error
	Recv() (*OrderId, error)
	grpc.ClientStream
}

type orderServiceProcessOrderClient struct {
	grpc.ClientStream
}

func (x *orderServiceProcessOrderClient) Send(m *OrderId) error {
	return x.ClientStream.SendMsg(m)
}

func (x *orderServiceProcessOrderClient) Recv() (*OrderId, error) {
	m := new(OrderId)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OrderServiceServer is the server API for OrderService service.
// All implementations must embed UnimplementedOrderServiceServer
// for forward compatibility
type OrderServiceServer interface {
	GreetCustomer(context.Context, *GreetRequest) (*GreetResponse, error)
	CreateOrderByIngredients(OrderService_CreateOrderByIngredientsServer) error
	GenerateRecommendedIngredientsCombo(*EmptyParam, OrderService_GenerateRecommendedIngredientsComboServer) error
	ProcessOrder(OrderService_ProcessOrderServer) error
	mustEmbedUnimplementedOrderServiceServer()
}

// UnimplementedOrderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOrderServiceServer struct {
}

func (UnimplementedOrderServiceServer) GreetCustomer(context.Context, *GreetRequest) (*GreetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GreetCustomer not implemented")
}
func (UnimplementedOrderServiceServer) CreateOrderByIngredients(OrderService_CreateOrderByIngredientsServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateOrderByIngredients not implemented")
}
func (UnimplementedOrderServiceServer) GenerateRecommendedIngredientsCombo(*EmptyParam, OrderService_GenerateRecommendedIngredientsComboServer) error {
	return status.Errorf(codes.Unimplemented, "method GenerateRecommendedIngredientsCombo not implemented")
}
func (UnimplementedOrderServiceServer) ProcessOrder(OrderService_ProcessOrderServer) error {
	return status.Errorf(codes.Unimplemented, "method ProcessOrder not implemented")
}
func (UnimplementedOrderServiceServer) mustEmbedUnimplementedOrderServiceServer() {}

// UnsafeOrderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServiceServer will
// result in compilation errors.
type UnsafeOrderServiceServer interface {
	mustEmbedUnimplementedOrderServiceServer()
}

func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
	s.RegisterService(&OrderService_ServiceDesc, srv)
}

func _OrderService_GreetCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GreetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).GreetCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order_proto.OrderService/GreetCustomer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).GreetCustomer(ctx, req.(*GreetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_CreateOrderByIngredients_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OrderServiceServer).CreateOrderByIngredients(&orderServiceCreateOrderByIngredientsServer{stream})
}

type OrderService_CreateOrderByIngredientsServer interface {
	SendAndClose(*OrderMessage) error
	Recv() (*Ingredient, error)
	grpc.ServerStream
}

type orderServiceCreateOrderByIngredientsServer struct {
	grpc.ServerStream
}

func (x *orderServiceCreateOrderByIngredientsServer) SendAndClose(m *OrderMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *orderServiceCreateOrderByIngredientsServer) Recv() (*Ingredient, error) {
	m := new(Ingredient)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _OrderService_GenerateRecommendedIngredientsCombo_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EmptyParam)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OrderServiceServer).GenerateRecommendedIngredientsCombo(m, &orderServiceGenerateRecommendedIngredientsComboServer{stream})
}

type OrderService_GenerateRecommendedIngredientsComboServer interface {
	Send(*Ingredient) error
	grpc.ServerStream
}

type orderServiceGenerateRecommendedIngredientsComboServer struct {
	grpc.ServerStream
}

func (x *orderServiceGenerateRecommendedIngredientsComboServer) Send(m *Ingredient) error {
	return x.ServerStream.SendMsg(m)
}

func _OrderService_ProcessOrder_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OrderServiceServer).ProcessOrder(&orderServiceProcessOrderServer{stream})
}

type OrderService_ProcessOrderServer interface {
	Send(*OrderId) error
	Recv() (*OrderId, error)
	grpc.ServerStream
}

type orderServiceProcessOrderServer struct {
	grpc.ServerStream
}

func (x *orderServiceProcessOrderServer) Send(m *OrderId) error {
	return x.ServerStream.SendMsg(m)
}

func (x *orderServiceProcessOrderServer) Recv() (*OrderId, error) {
	m := new(OrderId)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OrderService_ServiceDesc is the grpc.ServiceDesc for OrderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order_proto.OrderService",
	HandlerType: (*OrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GreetCustomer",
			Handler:    _OrderService_GreetCustomer_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateOrderByIngredients",
			Handler:       _OrderService_CreateOrderByIngredients_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GenerateRecommendedIngredientsCombo",
			Handler:       _OrderService_GenerateRecommendedIngredientsCombo_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ProcessOrder",
			Handler:       _OrderService_ProcessOrder_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/order.proto",
}
