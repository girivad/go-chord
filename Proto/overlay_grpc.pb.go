// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: Proto/overlay.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PredecessorClient is the client API for Predecessor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PredecessorClient interface {
	GetPredecessor(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*wrapperspb.StringValue, error)
	UpdatePredecessor(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type predecessorClient struct {
	cc grpc.ClientConnInterface
}

func NewPredecessorClient(cc grpc.ClientConnInterface) PredecessorClient {
	return &predecessorClient{cc}
}

func (c *predecessorClient) GetPredecessor(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	out := new(wrapperspb.StringValue)
	err := c.cc.Invoke(ctx, "/overlay.Predecessor/getPredecessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *predecessorClient) UpdatePredecessor(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/overlay.Predecessor/updatePredecessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PredecessorServer is the server API for Predecessor service.
// All implementations must embed UnimplementedPredecessorServer
// for forward compatibility
type PredecessorServer interface {
	GetPredecessor(context.Context, *emptypb.Empty) (*wrapperspb.StringValue, error)
	UpdatePredecessor(context.Context, *wrapperspb.StringValue) (*emptypb.Empty, error)
	mustEmbedUnimplementedPredecessorServer()
}

// UnimplementedPredecessorServer must be embedded to have forward compatible implementations.
type UnimplementedPredecessorServer struct {
}

func (UnimplementedPredecessorServer) GetPredecessor(context.Context, *emptypb.Empty) (*wrapperspb.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPredecessor not implemented")
}
func (UnimplementedPredecessorServer) UpdatePredecessor(context.Context, *wrapperspb.StringValue) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePredecessor not implemented")
}
func (UnimplementedPredecessorServer) mustEmbedUnimplementedPredecessorServer() {}

// UnsafePredecessorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PredecessorServer will
// result in compilation errors.
type UnsafePredecessorServer interface {
	mustEmbedUnimplementedPredecessorServer()
}

func RegisterPredecessorServer(s grpc.ServiceRegistrar, srv PredecessorServer) {
	s.RegisterService(&Predecessor_ServiceDesc, srv)
}

func _Predecessor_GetPredecessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PredecessorServer).GetPredecessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/overlay.Predecessor/getPredecessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PredecessorServer).GetPredecessor(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Predecessor_UpdatePredecessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PredecessorServer).UpdatePredecessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/overlay.Predecessor/updatePredecessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PredecessorServer).UpdatePredecessor(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

// Predecessor_ServiceDesc is the grpc.ServiceDesc for Predecessor service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Predecessor_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "overlay.Predecessor",
	HandlerType: (*PredecessorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getPredecessor",
			Handler:    _Predecessor_GetPredecessor_Handler,
		},
		{
			MethodName: "updatePredecessor",
			Handler:    _Predecessor_UpdatePredecessor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/overlay.proto",
}

// LookupClient is the client API for Lookup service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LookupClient interface {
	FindSuccessor(ctx context.Context, in *wrapperspb.Int64Value, opts ...grpc.CallOption) (*wrapperspb.StringValue, error)
}

type lookupClient struct {
	cc grpc.ClientConnInterface
}

func NewLookupClient(cc grpc.ClientConnInterface) LookupClient {
	return &lookupClient{cc}
}

func (c *lookupClient) FindSuccessor(ctx context.Context, in *wrapperspb.Int64Value, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	out := new(wrapperspb.StringValue)
	err := c.cc.Invoke(ctx, "/overlay.Lookup/findSuccessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LookupServer is the server API for Lookup service.
// All implementations must embed UnimplementedLookupServer
// for forward compatibility
type LookupServer interface {
	FindSuccessor(context.Context, *wrapperspb.Int64Value) (*wrapperspb.StringValue, error)
	mustEmbedUnimplementedLookupServer()
}

// UnimplementedLookupServer must be embedded to have forward compatible implementations.
type UnimplementedLookupServer struct {
}

func (UnimplementedLookupServer) FindSuccessor(context.Context, *wrapperspb.Int64Value) (*wrapperspb.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindSuccessor not implemented")
}
func (UnimplementedLookupServer) mustEmbedUnimplementedLookupServer() {}

// UnsafeLookupServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LookupServer will
// result in compilation errors.
type UnsafeLookupServer interface {
	mustEmbedUnimplementedLookupServer()
}

func RegisterLookupServer(s grpc.ServiceRegistrar, srv LookupServer) {
	s.RegisterService(&Lookup_ServiceDesc, srv)
}

func _Lookup_FindSuccessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.Int64Value)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LookupServer).FindSuccessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/overlay.Lookup/findSuccessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LookupServer).FindSuccessor(ctx, req.(*wrapperspb.Int64Value))
	}
	return interceptor(ctx, in, info, handler)
}

// Lookup_ServiceDesc is the grpc.ServiceDesc for Lookup service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Lookup_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "overlay.Lookup",
	HandlerType: (*LookupServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "findSuccessor",
			Handler:    _Lookup_FindSuccessor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/overlay.proto",
}

// CheckClient is the client API for Check service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CheckClient interface {
	LiveCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type checkClient struct {
	cc grpc.ClientConnInterface
}

func NewCheckClient(cc grpc.ClientConnInterface) CheckClient {
	return &checkClient{cc}
}

func (c *checkClient) LiveCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/overlay.Check/liveCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CheckServer is the server API for Check service.
// All implementations must embed UnimplementedCheckServer
// for forward compatibility
type CheckServer interface {
	LiveCheck(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedCheckServer()
}

// UnimplementedCheckServer must be embedded to have forward compatible implementations.
type UnimplementedCheckServer struct {
}

func (UnimplementedCheckServer) LiveCheck(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LiveCheck not implemented")
}
func (UnimplementedCheckServer) mustEmbedUnimplementedCheckServer() {}

// UnsafeCheckServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CheckServer will
// result in compilation errors.
type UnsafeCheckServer interface {
	mustEmbedUnimplementedCheckServer()
}

func RegisterCheckServer(s grpc.ServiceRegistrar, srv CheckServer) {
	s.RegisterService(&Check_ServiceDesc, srv)
}

func _Check_LiveCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckServer).LiveCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/overlay.Check/liveCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckServer).LiveCheck(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Check_ServiceDesc is the grpc.ServiceDesc for Check service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Check_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "overlay.Check",
	HandlerType: (*CheckServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "liveCheck",
			Handler:    _Check_LiveCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/overlay.proto",
}

// DataClient is the client API for Data service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataClient interface {
	TransferData(ctx context.Context, in *wrapperspb.Int64Value, opts ...grpc.CallOption) (*KVMap, error)
}

type dataClient struct {
	cc grpc.ClientConnInterface
}

func NewDataClient(cc grpc.ClientConnInterface) DataClient {
	return &dataClient{cc}
}

func (c *dataClient) TransferData(ctx context.Context, in *wrapperspb.Int64Value, opts ...grpc.CallOption) (*KVMap, error) {
	out := new(KVMap)
	err := c.cc.Invoke(ctx, "/overlay.Data/transferData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataServer is the server API for Data service.
// All implementations must embed UnimplementedDataServer
// for forward compatibility
type DataServer interface {
	TransferData(context.Context, *wrapperspb.Int64Value) (*KVMap, error)
	mustEmbedUnimplementedDataServer()
}

// UnimplementedDataServer must be embedded to have forward compatible implementations.
type UnimplementedDataServer struct {
}

func (UnimplementedDataServer) TransferData(context.Context, *wrapperspb.Int64Value) (*KVMap, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransferData not implemented")
}
func (UnimplementedDataServer) mustEmbedUnimplementedDataServer() {}

// UnsafeDataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataServer will
// result in compilation errors.
type UnsafeDataServer interface {
	mustEmbedUnimplementedDataServer()
}

func RegisterDataServer(s grpc.ServiceRegistrar, srv DataServer) {
	s.RegisterService(&Data_ServiceDesc, srv)
}

func _Data_TransferData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.Int64Value)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).TransferData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/overlay.Data/transferData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).TransferData(ctx, req.(*wrapperspb.Int64Value))
	}
	return interceptor(ctx, in, info, handler)
}

// Data_ServiceDesc is the grpc.ServiceDesc for Data service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Data_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "overlay.Data",
	HandlerType: (*DataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "transferData",
			Handler:    _Data_TransferData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/overlay.proto",
}