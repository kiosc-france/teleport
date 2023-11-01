// Copyright 2023 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: teleport/accesslist/v1/accesslist_service.proto

package accesslistv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AccessListService_GetAccessLists_FullMethodName                          = "/teleport.accesslist.v1.AccessListService/GetAccessLists"
	AccessListService_ListAccessLists_FullMethodName                         = "/teleport.accesslist.v1.AccessListService/ListAccessLists"
	AccessListService_GetAccessList_FullMethodName                           = "/teleport.accesslist.v1.AccessListService/GetAccessList"
	AccessListService_UpsertAccessList_FullMethodName                        = "/teleport.accesslist.v1.AccessListService/UpsertAccessList"
	AccessListService_DeleteAccessList_FullMethodName                        = "/teleport.accesslist.v1.AccessListService/DeleteAccessList"
	AccessListService_DeleteAllAccessLists_FullMethodName                    = "/teleport.accesslist.v1.AccessListService/DeleteAllAccessLists"
	AccessListService_GetAccessListsToReview_FullMethodName                  = "/teleport.accesslist.v1.AccessListService/GetAccessListsToReview"
	AccessListService_ListAccessListMembers_FullMethodName                   = "/teleport.accesslist.v1.AccessListService/ListAccessListMembers"
	AccessListService_GetAccessListMember_FullMethodName                     = "/teleport.accesslist.v1.AccessListService/GetAccessListMember"
	AccessListService_UpsertAccessListMember_FullMethodName                  = "/teleport.accesslist.v1.AccessListService/UpsertAccessListMember"
	AccessListService_DeleteAccessListMember_FullMethodName                  = "/teleport.accesslist.v1.AccessListService/DeleteAccessListMember"
	AccessListService_DeleteAllAccessListMembersForAccessList_FullMethodName = "/teleport.accesslist.v1.AccessListService/DeleteAllAccessListMembersForAccessList"
	AccessListService_DeleteAllAccessListMembers_FullMethodName              = "/teleport.accesslist.v1.AccessListService/DeleteAllAccessListMembers"
	AccessListService_UpsertAccessListWithMembers_FullMethodName             = "/teleport.accesslist.v1.AccessListService/UpsertAccessListWithMembers"
	AccessListService_ListAccessListReviews_FullMethodName                   = "/teleport.accesslist.v1.AccessListService/ListAccessListReviews"
	AccessListService_CreateAccessListReview_FullMethodName                  = "/teleport.accesslist.v1.AccessListService/CreateAccessListReview"
	AccessListService_DeleteAccessListReview_FullMethodName                  = "/teleport.accesslist.v1.AccessListService/DeleteAccessListReview"
	AccessListService_AccessRequestPromote_FullMethodName                    = "/teleport.accesslist.v1.AccessListService/AccessRequestPromote"
)

// AccessListServiceClient is the client API for AccessListService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccessListServiceClient interface {
	// GetAccessLists returns a list of all access lists.
	GetAccessLists(ctx context.Context, in *GetAccessListsRequest, opts ...grpc.CallOption) (*GetAccessListsResponse, error)
	// ListAccessLists returns a paginated list of all access lists.
	ListAccessLists(ctx context.Context, in *ListAccessListsRequest, opts ...grpc.CallOption) (*ListAccessListsResponse, error)
	// GetAccessList returns the specified access list resource.
	GetAccessList(ctx context.Context, in *GetAccessListRequest, opts ...grpc.CallOption) (*AccessList, error)
	// UpsertAccessList creates or updates an access list resource.
	UpsertAccessList(ctx context.Context, in *UpsertAccessListRequest, opts ...grpc.CallOption) (*AccessList, error)
	// DeleteAccessList hard deletes the specified access list resource.
	DeleteAccessList(ctx context.Context, in *DeleteAccessListRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteAllAccessLists hard deletes all access lists.
	DeleteAllAccessLists(ctx context.Context, in *DeleteAllAccessListsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// GetAccessListsToReview will return access lists that need to be reviewed by the current user.
	GetAccessListsToReview(ctx context.Context, in *GetAccessListsToReviewRequest, opts ...grpc.CallOption) (*GetAccessListsToReviewResponse, error)
	// ListAccessListMembers returns a paginated list of all access list members.
	ListAccessListMembers(ctx context.Context, in *ListAccessListMembersRequest, opts ...grpc.CallOption) (*ListAccessListMembersResponse, error)
	// GetAccessListMember returns the specified access list member resource.
	GetAccessListMember(ctx context.Context, in *GetAccessListMemberRequest, opts ...grpc.CallOption) (*Member, error)
	// UpsertAccessListMember creates or updates an access list member resource.
	UpsertAccessListMember(ctx context.Context, in *UpsertAccessListMemberRequest, opts ...grpc.CallOption) (*Member, error)
	// DeleteAccessListMember hard deletes the specified access list member resource.
	DeleteAccessListMember(ctx context.Context, in *DeleteAccessListMemberRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteAllAccessListMembers hard deletes all access list members for an access list.
	DeleteAllAccessListMembersForAccessList(ctx context.Context, in *DeleteAllAccessListMembersForAccessListRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteAllAccessListMembers hard deletes all access list members for an access list.
	DeleteAllAccessListMembers(ctx context.Context, in *DeleteAllAccessListMembersRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// UpsertAccessListWithMembers creates or updates an access list with members.
	UpsertAccessListWithMembers(ctx context.Context, in *UpsertAccessListWithMembersRequest, opts ...grpc.CallOption) (*UpsertAccessListWithMembersResponse, error)
	// ListAccessListReviews will list access list reviews for a particular access list.
	ListAccessListReviews(ctx context.Context, in *ListAccessListReviewsRequest, opts ...grpc.CallOption) (*ListAccessListReviewsResponse, error)
	// CreateAccessListReview will create a new review for an access list. It will also modify the original access list
	// and its members depending on the details of the review.
	CreateAccessListReview(ctx context.Context, in *CreateAccessListReviewRequest, opts ...grpc.CallOption) (*CreateAccessListReviewResponse, error)
	// DeleteAccessListReview will delete an access list review from the backend.
	DeleteAccessListReview(ctx context.Context, in *DeleteAccessListReviewRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// AccessRequestPromote promotes an access request to an access list.
	AccessRequestPromote(ctx context.Context, in *AccessRequestPromoteRequest, opts ...grpc.CallOption) (*AccessRequestPromoteResponse, error)
}

type accessListServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccessListServiceClient(cc grpc.ClientConnInterface) AccessListServiceClient {
	return &accessListServiceClient{cc}
}

func (c *accessListServiceClient) GetAccessLists(ctx context.Context, in *GetAccessListsRequest, opts ...grpc.CallOption) (*GetAccessListsResponse, error) {
	out := new(GetAccessListsResponse)
	err := c.cc.Invoke(ctx, AccessListService_GetAccessLists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) ListAccessLists(ctx context.Context, in *ListAccessListsRequest, opts ...grpc.CallOption) (*ListAccessListsResponse, error) {
	out := new(ListAccessListsResponse)
	err := c.cc.Invoke(ctx, AccessListService_ListAccessLists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) GetAccessList(ctx context.Context, in *GetAccessListRequest, opts ...grpc.CallOption) (*AccessList, error) {
	out := new(AccessList)
	err := c.cc.Invoke(ctx, AccessListService_GetAccessList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) UpsertAccessList(ctx context.Context, in *UpsertAccessListRequest, opts ...grpc.CallOption) (*AccessList, error) {
	out := new(AccessList)
	err := c.cc.Invoke(ctx, AccessListService_UpsertAccessList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) DeleteAccessList(ctx context.Context, in *DeleteAccessListRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccessListService_DeleteAccessList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) DeleteAllAccessLists(ctx context.Context, in *DeleteAllAccessListsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccessListService_DeleteAllAccessLists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) GetAccessListsToReview(ctx context.Context, in *GetAccessListsToReviewRequest, opts ...grpc.CallOption) (*GetAccessListsToReviewResponse, error) {
	out := new(GetAccessListsToReviewResponse)
	err := c.cc.Invoke(ctx, AccessListService_GetAccessListsToReview_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) ListAccessListMembers(ctx context.Context, in *ListAccessListMembersRequest, opts ...grpc.CallOption) (*ListAccessListMembersResponse, error) {
	out := new(ListAccessListMembersResponse)
	err := c.cc.Invoke(ctx, AccessListService_ListAccessListMembers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) GetAccessListMember(ctx context.Context, in *GetAccessListMemberRequest, opts ...grpc.CallOption) (*Member, error) {
	out := new(Member)
	err := c.cc.Invoke(ctx, AccessListService_GetAccessListMember_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) UpsertAccessListMember(ctx context.Context, in *UpsertAccessListMemberRequest, opts ...grpc.CallOption) (*Member, error) {
	out := new(Member)
	err := c.cc.Invoke(ctx, AccessListService_UpsertAccessListMember_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) DeleteAccessListMember(ctx context.Context, in *DeleteAccessListMemberRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccessListService_DeleteAccessListMember_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) DeleteAllAccessListMembersForAccessList(ctx context.Context, in *DeleteAllAccessListMembersForAccessListRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccessListService_DeleteAllAccessListMembersForAccessList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) DeleteAllAccessListMembers(ctx context.Context, in *DeleteAllAccessListMembersRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccessListService_DeleteAllAccessListMembers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) UpsertAccessListWithMembers(ctx context.Context, in *UpsertAccessListWithMembersRequest, opts ...grpc.CallOption) (*UpsertAccessListWithMembersResponse, error) {
	out := new(UpsertAccessListWithMembersResponse)
	err := c.cc.Invoke(ctx, AccessListService_UpsertAccessListWithMembers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) ListAccessListReviews(ctx context.Context, in *ListAccessListReviewsRequest, opts ...grpc.CallOption) (*ListAccessListReviewsResponse, error) {
	out := new(ListAccessListReviewsResponse)
	err := c.cc.Invoke(ctx, AccessListService_ListAccessListReviews_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) CreateAccessListReview(ctx context.Context, in *CreateAccessListReviewRequest, opts ...grpc.CallOption) (*CreateAccessListReviewResponse, error) {
	out := new(CreateAccessListReviewResponse)
	err := c.cc.Invoke(ctx, AccessListService_CreateAccessListReview_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) DeleteAccessListReview(ctx context.Context, in *DeleteAccessListReviewRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccessListService_DeleteAccessListReview_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessListServiceClient) AccessRequestPromote(ctx context.Context, in *AccessRequestPromoteRequest, opts ...grpc.CallOption) (*AccessRequestPromoteResponse, error) {
	out := new(AccessRequestPromoteResponse)
	err := c.cc.Invoke(ctx, AccessListService_AccessRequestPromote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccessListServiceServer is the server API for AccessListService service.
// All implementations must embed UnimplementedAccessListServiceServer
// for forward compatibility
type AccessListServiceServer interface {
	// GetAccessLists returns a list of all access lists.
	GetAccessLists(context.Context, *GetAccessListsRequest) (*GetAccessListsResponse, error)
	// ListAccessLists returns a paginated list of all access lists.
	ListAccessLists(context.Context, *ListAccessListsRequest) (*ListAccessListsResponse, error)
	// GetAccessList returns the specified access list resource.
	GetAccessList(context.Context, *GetAccessListRequest) (*AccessList, error)
	// UpsertAccessList creates or updates an access list resource.
	UpsertAccessList(context.Context, *UpsertAccessListRequest) (*AccessList, error)
	// DeleteAccessList hard deletes the specified access list resource.
	DeleteAccessList(context.Context, *DeleteAccessListRequest) (*emptypb.Empty, error)
	// DeleteAllAccessLists hard deletes all access lists.
	DeleteAllAccessLists(context.Context, *DeleteAllAccessListsRequest) (*emptypb.Empty, error)
	// GetAccessListsToReview will return access lists that need to be reviewed by the current user.
	GetAccessListsToReview(context.Context, *GetAccessListsToReviewRequest) (*GetAccessListsToReviewResponse, error)
	// ListAccessListMembers returns a paginated list of all access list members.
	ListAccessListMembers(context.Context, *ListAccessListMembersRequest) (*ListAccessListMembersResponse, error)
	// GetAccessListMember returns the specified access list member resource.
	GetAccessListMember(context.Context, *GetAccessListMemberRequest) (*Member, error)
	// UpsertAccessListMember creates or updates an access list member resource.
	UpsertAccessListMember(context.Context, *UpsertAccessListMemberRequest) (*Member, error)
	// DeleteAccessListMember hard deletes the specified access list member resource.
	DeleteAccessListMember(context.Context, *DeleteAccessListMemberRequest) (*emptypb.Empty, error)
	// DeleteAllAccessListMembers hard deletes all access list members for an access list.
	DeleteAllAccessListMembersForAccessList(context.Context, *DeleteAllAccessListMembersForAccessListRequest) (*emptypb.Empty, error)
	// DeleteAllAccessListMembers hard deletes all access list members for an access list.
	DeleteAllAccessListMembers(context.Context, *DeleteAllAccessListMembersRequest) (*emptypb.Empty, error)
	// UpsertAccessListWithMembers creates or updates an access list with members.
	UpsertAccessListWithMembers(context.Context, *UpsertAccessListWithMembersRequest) (*UpsertAccessListWithMembersResponse, error)
	// ListAccessListReviews will list access list reviews for a particular access list.
	ListAccessListReviews(context.Context, *ListAccessListReviewsRequest) (*ListAccessListReviewsResponse, error)
	// CreateAccessListReview will create a new review for an access list. It will also modify the original access list
	// and its members depending on the details of the review.
	CreateAccessListReview(context.Context, *CreateAccessListReviewRequest) (*CreateAccessListReviewResponse, error)
	// DeleteAccessListReview will delete an access list review from the backend.
	DeleteAccessListReview(context.Context, *DeleteAccessListReviewRequest) (*emptypb.Empty, error)
	// AccessRequestPromote promotes an access request to an access list.
	AccessRequestPromote(context.Context, *AccessRequestPromoteRequest) (*AccessRequestPromoteResponse, error)
	mustEmbedUnimplementedAccessListServiceServer()
}

// UnimplementedAccessListServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccessListServiceServer struct {
}

func (UnimplementedAccessListServiceServer) GetAccessLists(context.Context, *GetAccessListsRequest) (*GetAccessListsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccessLists not implemented")
}
func (UnimplementedAccessListServiceServer) ListAccessLists(context.Context, *ListAccessListsRequest) (*ListAccessListsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccessLists not implemented")
}
func (UnimplementedAccessListServiceServer) GetAccessList(context.Context, *GetAccessListRequest) (*AccessList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccessList not implemented")
}
func (UnimplementedAccessListServiceServer) UpsertAccessList(context.Context, *UpsertAccessListRequest) (*AccessList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertAccessList not implemented")
}
func (UnimplementedAccessListServiceServer) DeleteAccessList(context.Context, *DeleteAccessListRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccessList not implemented")
}
func (UnimplementedAccessListServiceServer) DeleteAllAccessLists(context.Context, *DeleteAllAccessListsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAllAccessLists not implemented")
}
func (UnimplementedAccessListServiceServer) GetAccessListsToReview(context.Context, *GetAccessListsToReviewRequest) (*GetAccessListsToReviewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccessListsToReview not implemented")
}
func (UnimplementedAccessListServiceServer) ListAccessListMembers(context.Context, *ListAccessListMembersRequest) (*ListAccessListMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccessListMembers not implemented")
}
func (UnimplementedAccessListServiceServer) GetAccessListMember(context.Context, *GetAccessListMemberRequest) (*Member, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccessListMember not implemented")
}
func (UnimplementedAccessListServiceServer) UpsertAccessListMember(context.Context, *UpsertAccessListMemberRequest) (*Member, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertAccessListMember not implemented")
}
func (UnimplementedAccessListServiceServer) DeleteAccessListMember(context.Context, *DeleteAccessListMemberRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccessListMember not implemented")
}
func (UnimplementedAccessListServiceServer) DeleteAllAccessListMembersForAccessList(context.Context, *DeleteAllAccessListMembersForAccessListRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAllAccessListMembersForAccessList not implemented")
}
func (UnimplementedAccessListServiceServer) DeleteAllAccessListMembers(context.Context, *DeleteAllAccessListMembersRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAllAccessListMembers not implemented")
}
func (UnimplementedAccessListServiceServer) UpsertAccessListWithMembers(context.Context, *UpsertAccessListWithMembersRequest) (*UpsertAccessListWithMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertAccessListWithMembers not implemented")
}
func (UnimplementedAccessListServiceServer) ListAccessListReviews(context.Context, *ListAccessListReviewsRequest) (*ListAccessListReviewsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccessListReviews not implemented")
}
func (UnimplementedAccessListServiceServer) CreateAccessListReview(context.Context, *CreateAccessListReviewRequest) (*CreateAccessListReviewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccessListReview not implemented")
}
func (UnimplementedAccessListServiceServer) DeleteAccessListReview(context.Context, *DeleteAccessListReviewRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccessListReview not implemented")
}
func (UnimplementedAccessListServiceServer) AccessRequestPromote(context.Context, *AccessRequestPromoteRequest) (*AccessRequestPromoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AccessRequestPromote not implemented")
}
func (UnimplementedAccessListServiceServer) mustEmbedUnimplementedAccessListServiceServer() {}

// UnsafeAccessListServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccessListServiceServer will
// result in compilation errors.
type UnsafeAccessListServiceServer interface {
	mustEmbedUnimplementedAccessListServiceServer()
}

func RegisterAccessListServiceServer(s grpc.ServiceRegistrar, srv AccessListServiceServer) {
	s.RegisterService(&AccessListService_ServiceDesc, srv)
}

func _AccessListService_GetAccessLists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccessListsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).GetAccessLists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_GetAccessLists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).GetAccessLists(ctx, req.(*GetAccessListsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_ListAccessLists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAccessListsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).ListAccessLists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_ListAccessLists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).ListAccessLists(ctx, req.(*ListAccessListsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_GetAccessList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccessListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).GetAccessList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_GetAccessList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).GetAccessList(ctx, req.(*GetAccessListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_UpsertAccessList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpsertAccessListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).UpsertAccessList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_UpsertAccessList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).UpsertAccessList(ctx, req.(*UpsertAccessListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_DeleteAccessList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAccessListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).DeleteAccessList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_DeleteAccessList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).DeleteAccessList(ctx, req.(*DeleteAccessListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_DeleteAllAccessLists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAllAccessListsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).DeleteAllAccessLists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_DeleteAllAccessLists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).DeleteAllAccessLists(ctx, req.(*DeleteAllAccessListsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_GetAccessListsToReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccessListsToReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).GetAccessListsToReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_GetAccessListsToReview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).GetAccessListsToReview(ctx, req.(*GetAccessListsToReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_ListAccessListMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAccessListMembersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).ListAccessListMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_ListAccessListMembers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).ListAccessListMembers(ctx, req.(*ListAccessListMembersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_GetAccessListMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccessListMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).GetAccessListMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_GetAccessListMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).GetAccessListMember(ctx, req.(*GetAccessListMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_UpsertAccessListMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpsertAccessListMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).UpsertAccessListMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_UpsertAccessListMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).UpsertAccessListMember(ctx, req.(*UpsertAccessListMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_DeleteAccessListMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAccessListMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).DeleteAccessListMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_DeleteAccessListMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).DeleteAccessListMember(ctx, req.(*DeleteAccessListMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_DeleteAllAccessListMembersForAccessList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAllAccessListMembersForAccessListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).DeleteAllAccessListMembersForAccessList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_DeleteAllAccessListMembersForAccessList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).DeleteAllAccessListMembersForAccessList(ctx, req.(*DeleteAllAccessListMembersForAccessListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_DeleteAllAccessListMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAllAccessListMembersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).DeleteAllAccessListMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_DeleteAllAccessListMembers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).DeleteAllAccessListMembers(ctx, req.(*DeleteAllAccessListMembersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_UpsertAccessListWithMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpsertAccessListWithMembersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).UpsertAccessListWithMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_UpsertAccessListWithMembers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).UpsertAccessListWithMembers(ctx, req.(*UpsertAccessListWithMembersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_ListAccessListReviews_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAccessListReviewsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).ListAccessListReviews(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_ListAccessListReviews_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).ListAccessListReviews(ctx, req.(*ListAccessListReviewsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_CreateAccessListReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccessListReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).CreateAccessListReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_CreateAccessListReview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).CreateAccessListReview(ctx, req.(*CreateAccessListReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_DeleteAccessListReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAccessListReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).DeleteAccessListReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_DeleteAccessListReview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).DeleteAccessListReview(ctx, req.(*DeleteAccessListReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessListService_AccessRequestPromote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessRequestPromoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessListServiceServer).AccessRequestPromote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessListService_AccessRequestPromote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessListServiceServer).AccessRequestPromote(ctx, req.(*AccessRequestPromoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccessListService_ServiceDesc is the grpc.ServiceDesc for AccessListService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccessListService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "teleport.accesslist.v1.AccessListService",
	HandlerType: (*AccessListServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAccessLists",
			Handler:    _AccessListService_GetAccessLists_Handler,
		},
		{
			MethodName: "ListAccessLists",
			Handler:    _AccessListService_ListAccessLists_Handler,
		},
		{
			MethodName: "GetAccessList",
			Handler:    _AccessListService_GetAccessList_Handler,
		},
		{
			MethodName: "UpsertAccessList",
			Handler:    _AccessListService_UpsertAccessList_Handler,
		},
		{
			MethodName: "DeleteAccessList",
			Handler:    _AccessListService_DeleteAccessList_Handler,
		},
		{
			MethodName: "DeleteAllAccessLists",
			Handler:    _AccessListService_DeleteAllAccessLists_Handler,
		},
		{
			MethodName: "GetAccessListsToReview",
			Handler:    _AccessListService_GetAccessListsToReview_Handler,
		},
		{
			MethodName: "ListAccessListMembers",
			Handler:    _AccessListService_ListAccessListMembers_Handler,
		},
		{
			MethodName: "GetAccessListMember",
			Handler:    _AccessListService_GetAccessListMember_Handler,
		},
		{
			MethodName: "UpsertAccessListMember",
			Handler:    _AccessListService_UpsertAccessListMember_Handler,
		},
		{
			MethodName: "DeleteAccessListMember",
			Handler:    _AccessListService_DeleteAccessListMember_Handler,
		},
		{
			MethodName: "DeleteAllAccessListMembersForAccessList",
			Handler:    _AccessListService_DeleteAllAccessListMembersForAccessList_Handler,
		},
		{
			MethodName: "DeleteAllAccessListMembers",
			Handler:    _AccessListService_DeleteAllAccessListMembers_Handler,
		},
		{
			MethodName: "UpsertAccessListWithMembers",
			Handler:    _AccessListService_UpsertAccessListWithMembers_Handler,
		},
		{
			MethodName: "ListAccessListReviews",
			Handler:    _AccessListService_ListAccessListReviews_Handler,
		},
		{
			MethodName: "CreateAccessListReview",
			Handler:    _AccessListService_CreateAccessListReview_Handler,
		},
		{
			MethodName: "DeleteAccessListReview",
			Handler:    _AccessListService_DeleteAccessListReview_Handler,
		},
		{
			MethodName: "AccessRequestPromote",
			Handler:    _AccessListService_AccessRequestPromote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "teleport/accesslist/v1/accesslist_service.proto",
}
