// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: proto/auction/auction.proto

package auction

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
	AuctionService_StartAuction_FullMethodName    = "/auction.AuctionService/StartAuction"
	AuctionService_SubmitBids_FullMethodName      = "/auction.AuctionService/SubmitBids"
	AuctionService_GetAuctionInfo_FullMethodName  = "/auction.AuctionService/GetAuctionInfo"
	AuctionService_GetLatestTob_FullMethodName    = "/auction.AuctionService/GetLatestTob"
	AuctionService_GetAuctionState_FullMethodName = "/auction.AuctionService/GetAuctionState"
)

// AuctionServiceClient is the client API for AuctionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// AuctionService defines the RPC methods for auction operations.
type AuctionServiceClient interface {
	// Initiates a new auction.
	StartAuction(ctx context.Context, in *StartAuctionRequest, opts ...grpc.CallOption) (*StartAuctionResponse, error)
	// Submits multiple bids for a specific auction.
	SubmitBids(ctx context.Context, in *SubmitBidsRequest, opts ...grpc.CallOption) (*SubmitBidsResponse, error)
	// Retrieves detailed information about a specific auction.
	GetAuctionInfo(ctx context.Context, in *GetAuctionInfoRequest, opts ...grpc.CallOption) (*GetAuctionInfoResponse, error)
	// Retrieves the Tx list of the latest block.
	GetLatestTob(ctx context.Context, in *GetLatestTobRequest, opts ...grpc.CallOption) (*GetLatestTobResponse, error)
	// Retrieves the current state of an auction.
	GetAuctionState(ctx context.Context, in *GetAuctionStateRequest, opts ...grpc.CallOption) (*GetAuctionStateResponse, error)
}

type auctionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuctionServiceClient(cc grpc.ClientConnInterface) AuctionServiceClient {
	return &auctionServiceClient{cc}
}

func (c *auctionServiceClient) StartAuction(ctx context.Context, in *StartAuctionRequest, opts ...grpc.CallOption) (*StartAuctionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StartAuctionResponse)
	err := c.cc.Invoke(ctx, AuctionService_StartAuction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) SubmitBids(ctx context.Context, in *SubmitBidsRequest, opts ...grpc.CallOption) (*SubmitBidsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SubmitBidsResponse)
	err := c.cc.Invoke(ctx, AuctionService_SubmitBids_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) GetAuctionInfo(ctx context.Context, in *GetAuctionInfoRequest, opts ...grpc.CallOption) (*GetAuctionInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAuctionInfoResponse)
	err := c.cc.Invoke(ctx, AuctionService_GetAuctionInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) GetLatestTob(ctx context.Context, in *GetLatestTobRequest, opts ...grpc.CallOption) (*GetLatestTobResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetLatestTobResponse)
	err := c.cc.Invoke(ctx, AuctionService_GetLatestTob_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) GetAuctionState(ctx context.Context, in *GetAuctionStateRequest, opts ...grpc.CallOption) (*GetAuctionStateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAuctionStateResponse)
	err := c.cc.Invoke(ctx, AuctionService_GetAuctionState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuctionServiceServer is the server API for AuctionService service.
// All implementations must embed UnimplementedAuctionServiceServer
// for forward compatibility.
//
// AuctionService defines the RPC methods for auction operations.
type AuctionServiceServer interface {
	// Initiates a new auction.
	StartAuction(context.Context, *StartAuctionRequest) (*StartAuctionResponse, error)
	// Submits multiple bids for a specific auction.
	SubmitBids(context.Context, *SubmitBidsRequest) (*SubmitBidsResponse, error)
	// Retrieves detailed information about a specific auction.
	GetAuctionInfo(context.Context, *GetAuctionInfoRequest) (*GetAuctionInfoResponse, error)
	// Retrieves the Tx list of the latest block.
	GetLatestTob(context.Context, *GetLatestTobRequest) (*GetLatestTobResponse, error)
	// Retrieves the current state of an auction.
	GetAuctionState(context.Context, *GetAuctionStateRequest) (*GetAuctionStateResponse, error)
	mustEmbedUnimplementedAuctionServiceServer()
}

// UnimplementedAuctionServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuctionServiceServer struct{}

func (UnimplementedAuctionServiceServer) StartAuction(context.Context, *StartAuctionRequest) (*StartAuctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartAuction not implemented")
}
func (UnimplementedAuctionServiceServer) SubmitBids(context.Context, *SubmitBidsRequest) (*SubmitBidsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitBids not implemented")
}
func (UnimplementedAuctionServiceServer) GetAuctionInfo(context.Context, *GetAuctionInfoRequest) (*GetAuctionInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuctionInfo not implemented")
}
func (UnimplementedAuctionServiceServer) GetLatestTob(context.Context, *GetLatestTobRequest) (*GetLatestTobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLatestTob not implemented")
}
func (UnimplementedAuctionServiceServer) GetAuctionState(context.Context, *GetAuctionStateRequest) (*GetAuctionStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuctionState not implemented")
}
func (UnimplementedAuctionServiceServer) mustEmbedUnimplementedAuctionServiceServer() {}
func (UnimplementedAuctionServiceServer) testEmbeddedByValue()                        {}

// UnsafeAuctionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuctionServiceServer will
// result in compilation errors.
type UnsafeAuctionServiceServer interface {
	mustEmbedUnimplementedAuctionServiceServer()
}

func RegisterAuctionServiceServer(s grpc.ServiceRegistrar, srv AuctionServiceServer) {
	// If the following call pancis, it indicates UnimplementedAuctionServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AuctionService_ServiceDesc, srv)
}

func _AuctionService_StartAuction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartAuctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).StartAuction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuctionService_StartAuction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).StartAuction(ctx, req.(*StartAuctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_SubmitBids_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitBidsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).SubmitBids(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuctionService_SubmitBids_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).SubmitBids(ctx, req.(*SubmitBidsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_GetAuctionInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuctionInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).GetAuctionInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuctionService_GetAuctionInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).GetAuctionInfo(ctx, req.(*GetAuctionInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_GetLatestTob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLatestTobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).GetLatestTob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuctionService_GetLatestTob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).GetLatestTob(ctx, req.(*GetLatestTobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_GetAuctionState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuctionStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).GetAuctionState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuctionService_GetAuctionState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).GetAuctionState(ctx, req.(*GetAuctionStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuctionService_ServiceDesc is the grpc.ServiceDesc for AuctionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuctionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auction.AuctionService",
	HandlerType: (*AuctionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartAuction",
			Handler:    _AuctionService_StartAuction_Handler,
		},
		{
			MethodName: "SubmitBids",
			Handler:    _AuctionService_SubmitBids_Handler,
		},
		{
			MethodName: "GetAuctionInfo",
			Handler:    _AuctionService_GetAuctionInfo_Handler,
		},
		{
			MethodName: "GetLatestTob",
			Handler:    _AuctionService_GetLatestTob_Handler,
		},
		{
			MethodName: "GetAuctionState",
			Handler:    _AuctionService_GetAuctionState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/auction/auction.proto",
}
