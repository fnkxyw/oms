// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: PuP-service/v1/pup_service.proto

package pup_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PupService_AcceptOrder_FullMethodName = "/pup_service.PupService/AcceptOrder"
	PupService_PlaceOrder_FullMethodName  = "/pup_service.PupService/PlaceOrder"
	PupService_ReturnOrder_FullMethodName = "/pup_service.PupService/ReturnOrder"
	PupService_ListOrder_FullMethodName   = "/pup_service.PupService/ListOrder"
	PupService_RefundOrder_FullMethodName = "/pup_service.PupService/RefundOrder"
	PupService_ListReturns_FullMethodName = "/pup_service.PupService/ListReturns"
	PupService_WorkerNum_FullMethodName   = "/pup_service.PupService/WorkerNum"
)

// PupServiceClient is the client API for PupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PupServiceClient interface {
	AcceptOrder(ctx context.Context, in *AcceptOrderRequest, opts ...grpc.CallOption) (*AcceptOrderResponse, error)
	PlaceOrder(ctx context.Context, in *PlaceOrderRequest, opts ...grpc.CallOption) (*PlaceOrderResponse, error)
	ReturnOrder(ctx context.Context, in *ReturnOrderRequest, opts ...grpc.CallOption) (*ReturnOrderResponse, error)
	ListOrder(ctx context.Context, in *ListOrdersRequest, opts ...grpc.CallOption) (*ListOrderResponse, error)
	RefundOrder(ctx context.Context, in *RefundOrderRequest, opts ...grpc.CallOption) (*RefundOrderResponse, error)
	ListReturns(ctx context.Context, in *ListReturnsRequest, opts ...grpc.CallOption) (*ListReturnsResponse, error)
	WorkerNum(ctx context.Context, in *WorkersNumRequest, opts ...grpc.CallOption) (*WorkersNumResponse, error)
}

type pupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPupServiceClient(cc grpc.ClientConnInterface) PupServiceClient {
	return &pupServiceClient{cc}
}

func (c *pupServiceClient) AcceptOrder(ctx context.Context, in *AcceptOrderRequest, opts ...grpc.CallOption) (*AcceptOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AcceptOrderResponse)
	err := c.cc.Invoke(ctx, PupService_AcceptOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pupServiceClient) PlaceOrder(ctx context.Context, in *PlaceOrderRequest, opts ...grpc.CallOption) (*PlaceOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PlaceOrderResponse)
	err := c.cc.Invoke(ctx, PupService_PlaceOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pupServiceClient) ReturnOrder(ctx context.Context, in *ReturnOrderRequest, opts ...grpc.CallOption) (*ReturnOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReturnOrderResponse)
	err := c.cc.Invoke(ctx, PupService_ReturnOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pupServiceClient) ListOrder(ctx context.Context, in *ListOrdersRequest, opts ...grpc.CallOption) (*ListOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListOrderResponse)
	err := c.cc.Invoke(ctx, PupService_ListOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pupServiceClient) RefundOrder(ctx context.Context, in *RefundOrderRequest, opts ...grpc.CallOption) (*RefundOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RefundOrderResponse)
	err := c.cc.Invoke(ctx, PupService_RefundOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pupServiceClient) ListReturns(ctx context.Context, in *ListReturnsRequest, opts ...grpc.CallOption) (*ListReturnsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListReturnsResponse)
	err := c.cc.Invoke(ctx, PupService_ListReturns_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pupServiceClient) WorkerNum(ctx context.Context, in *WorkersNumRequest, opts ...grpc.CallOption) (*WorkersNumResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WorkersNumResponse)
	err := c.cc.Invoke(ctx, PupService_WorkerNum_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PupServiceServer is the server API for PupService service.
// All implementations must embed UnimplementedPupServiceServer
// for forward compatibility.
type PupServiceServer interface {
	AcceptOrder(context.Context, *AcceptOrderRequest) (*AcceptOrderResponse, error)
	PlaceOrder(context.Context, *PlaceOrderRequest) (*PlaceOrderResponse, error)
	ReturnOrder(context.Context, *ReturnOrderRequest) (*ReturnOrderResponse, error)
	ListOrder(context.Context, *ListOrdersRequest) (*ListOrderResponse, error)
	RefundOrder(context.Context, *RefundOrderRequest) (*RefundOrderResponse, error)
	ListReturns(context.Context, *ListReturnsRequest) (*ListReturnsResponse, error)
	WorkerNum(context.Context, *WorkersNumRequest) (*WorkersNumResponse, error)
	mustEmbedUnimplementedPupServiceServer()
}

// UnimplementedPupServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPupServiceServer struct{}

func (UnimplementedPupServiceServer) AcceptOrder(context.Context, *AcceptOrderRequest) (*AcceptOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptOrder not implemented")
}
func (UnimplementedPupServiceServer) PlaceOrder(context.Context, *PlaceOrderRequest) (*PlaceOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceOrder not implemented")
}
func (UnimplementedPupServiceServer) ReturnOrder(context.Context, *ReturnOrderRequest) (*ReturnOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReturnOrder not implemented")
}
func (UnimplementedPupServiceServer) ListOrder(context.Context, *ListOrdersRequest) (*ListOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOrder not implemented")
}
func (UnimplementedPupServiceServer) RefundOrder(context.Context, *RefundOrderRequest) (*RefundOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefundOrder not implemented")
}
func (UnimplementedPupServiceServer) ListReturns(context.Context, *ListReturnsRequest) (*ListReturnsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListReturns not implemented")
}
func (UnimplementedPupServiceServer) WorkerNum(context.Context, *WorkersNumRequest) (*WorkersNumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WorkerNum not implemented")
}
func (UnimplementedPupServiceServer) mustEmbedUnimplementedPupServiceServer() {}
func (UnimplementedPupServiceServer) testEmbeddedByValue()                    {}

// UnsafePupServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PupServiceServer will
// result in compilation errors.
type UnsafePupServiceServer interface {
	mustEmbedUnimplementedPupServiceServer()
}

func RegisterPupServiceServer(s grpc.ServiceRegistrar, srv PupServiceServer) {
	// If the following call pancis, it indicates UnimplementedPupServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PupService_ServiceDesc, srv)
}

func _PupService_AcceptOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PupServiceServer).AcceptOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PupService_AcceptOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PupServiceServer).AcceptOrder(ctx, req.(*AcceptOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PupService_PlaceOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaceOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PupServiceServer).PlaceOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PupService_PlaceOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PupServiceServer).PlaceOrder(ctx, req.(*PlaceOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PupService_ReturnOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReturnOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PupServiceServer).ReturnOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PupService_ReturnOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PupServiceServer).ReturnOrder(ctx, req.(*ReturnOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PupService_ListOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListOrdersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PupServiceServer).ListOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PupService_ListOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PupServiceServer).ListOrder(ctx, req.(*ListOrdersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PupService_RefundOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefundOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PupServiceServer).RefundOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PupService_RefundOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PupServiceServer).RefundOrder(ctx, req.(*RefundOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PupService_ListReturns_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListReturnsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PupServiceServer).ListReturns(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PupService_ListReturns_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PupServiceServer).ListReturns(ctx, req.(*ListReturnsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PupService_WorkerNum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkersNumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PupServiceServer).WorkerNum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PupService_WorkerNum_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PupServiceServer).WorkerNum(ctx, req.(*WorkersNumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PupService_ServiceDesc is the grpc.ServiceDesc for PupService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pup_service.PupService",
	HandlerType: (*PupServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AcceptOrder",
			Handler:    _PupService_AcceptOrder_Handler,
		},
		{
			MethodName: "PlaceOrder",
			Handler:    _PupService_PlaceOrder_Handler,
		},
		{
			MethodName: "ReturnOrder",
			Handler:    _PupService_ReturnOrder_Handler,
		},
		{
			MethodName: "ListOrder",
			Handler:    _PupService_ListOrder_Handler,
		},
		{
			MethodName: "RefundOrder",
			Handler:    _PupService_RefundOrder_Handler,
		},
		{
			MethodName: "ListReturns",
			Handler:    _PupService_ListReturns_Handler,
		},
		{
			MethodName: "WorkerNum",
			Handler:    _PupService_WorkerNum_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "PuP-service/v1/pup_service.proto",
}
