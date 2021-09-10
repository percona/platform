/**
 * @fileoverview gRPC-Web generated client stub for percona.platform.org.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as org_org_api_pb from '../org/org_api_pb';


export class OrgAPIClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname;
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodInfoCreateOrganization = new grpcWeb.AbstractClientBase.MethodInfo(
    org_org_api_pb.CreateOrganizationResponse,
    (request: org_org_api_pb.CreateOrganizationRequest) => {
      return request.serializeBinary();
    },
    org_org_api_pb.CreateOrganizationResponse.deserializeBinary
  );

  createOrganization(
    request: org_org_api_pb.CreateOrganizationRequest,
    metadata: grpcWeb.Metadata | null): Promise<org_org_api_pb.CreateOrganizationResponse>;

  createOrganization(
    request: org_org_api_pb.CreateOrganizationRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: org_org_api_pb.CreateOrganizationResponse) => void): grpcWeb.ClientReadableStream<org_org_api_pb.CreateOrganizationResponse>;

  createOrganization(
    request: org_org_api_pb.CreateOrganizationRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: org_org_api_pb.CreateOrganizationResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.org.v1.OrgAPI/CreateOrganization',
        request,
        metadata || {},
        this.methodInfoCreateOrganization,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.org.v1.OrgAPI/CreateOrganization',
    request,
    metadata || {},
    this.methodInfoCreateOrganization);
  }

  methodInfoGetOrganization = new grpcWeb.AbstractClientBase.MethodInfo(
    org_org_api_pb.GetOrganizationResponse,
    (request: org_org_api_pb.GetOrganizationRequest) => {
      return request.serializeBinary();
    },
    org_org_api_pb.GetOrganizationResponse.deserializeBinary
  );

  getOrganization(
    request: org_org_api_pb.GetOrganizationRequest,
    metadata: grpcWeb.Metadata | null): Promise<org_org_api_pb.GetOrganizationResponse>;

  getOrganization(
    request: org_org_api_pb.GetOrganizationRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: org_org_api_pb.GetOrganizationResponse) => void): grpcWeb.ClientReadableStream<org_org_api_pb.GetOrganizationResponse>;

  getOrganization(
    request: org_org_api_pb.GetOrganizationRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: org_org_api_pb.GetOrganizationResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.org.v1.OrgAPI/GetOrganization',
        request,
        metadata || {},
        this.methodInfoGetOrganization,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.org.v1.OrgAPI/GetOrganization',
    request,
    metadata || {},
    this.methodInfoGetOrganization);
  }

  methodInfoGetOrganizationByUser = new grpcWeb.AbstractClientBase.MethodInfo(
    org_org_api_pb.GetOrganizationByUserResponse,
    (request: org_org_api_pb.GetOrganizationByUserRequest) => {
      return request.serializeBinary();
    },
    org_org_api_pb.GetOrganizationByUserResponse.deserializeBinary
  );

  getOrganizationByUser(
    request: org_org_api_pb.GetOrganizationByUserRequest,
    metadata: grpcWeb.Metadata | null): Promise<org_org_api_pb.GetOrganizationByUserResponse>;

  getOrganizationByUser(
    request: org_org_api_pb.GetOrganizationByUserRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: org_org_api_pb.GetOrganizationByUserResponse) => void): grpcWeb.ClientReadableStream<org_org_api_pb.GetOrganizationByUserResponse>;

  getOrganizationByUser(
    request: org_org_api_pb.GetOrganizationByUserRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: org_org_api_pb.GetOrganizationByUserResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.org.v1.OrgAPI/GetOrganizationByUser',
        request,
        metadata || {},
        this.methodInfoGetOrganizationByUser,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.org.v1.OrgAPI/GetOrganizationByUser',
    request,
    metadata || {},
    this.methodInfoGetOrganizationByUser);
  }

  methodInfoInviteMember = new grpcWeb.AbstractClientBase.MethodInfo(
    org_org_api_pb.InviteMemberResponse,
    (request: org_org_api_pb.InviteMemberRequest) => {
      return request.serializeBinary();
    },
    org_org_api_pb.InviteMemberResponse.deserializeBinary
  );

  inviteMember(
    request: org_org_api_pb.InviteMemberRequest,
    metadata: grpcWeb.Metadata | null): Promise<org_org_api_pb.InviteMemberResponse>;

  inviteMember(
    request: org_org_api_pb.InviteMemberRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: org_org_api_pb.InviteMemberResponse) => void): grpcWeb.ClientReadableStream<org_org_api_pb.InviteMemberResponse>;

  inviteMember(
    request: org_org_api_pb.InviteMemberRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: org_org_api_pb.InviteMemberResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.org.v1.OrgAPI/InviteMember',
        request,
        metadata || {},
        this.methodInfoInviteMember,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.org.v1.OrgAPI/InviteMember',
    request,
    metadata || {},
    this.methodInfoInviteMember);
  }

  methodInfoListMembers = new grpcWeb.AbstractClientBase.MethodInfo(
    org_org_api_pb.ListMembersResponse,
    (request: org_org_api_pb.ListMembersRequest) => {
      return request.serializeBinary();
    },
    org_org_api_pb.ListMembersResponse.deserializeBinary
  );

  listMembers(
    request: org_org_api_pb.ListMembersRequest,
    metadata: grpcWeb.Metadata | null): Promise<org_org_api_pb.ListMembersResponse>;

  listMembers(
    request: org_org_api_pb.ListMembersRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: org_org_api_pb.ListMembersResponse) => void): grpcWeb.ClientReadableStream<org_org_api_pb.ListMembersResponse>;

  listMembers(
    request: org_org_api_pb.ListMembersRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: org_org_api_pb.ListMembersResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.org.v1.OrgAPI/ListMembers',
        request,
        metadata || {},
        this.methodInfoListMembers,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.org.v1.OrgAPI/ListMembers',
    request,
    metadata || {},
    this.methodInfoListMembers);
  }

}

