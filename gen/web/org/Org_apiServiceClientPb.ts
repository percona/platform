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

  methodInfoSearchOrganizations = new grpcWeb.AbstractClientBase.MethodInfo(
    org_org_api_pb.SearchOrganizationsResponse,
    (request: org_org_api_pb.SearchOrganizationsRequest) => {
      return request.serializeBinary();
    },
    org_org_api_pb.SearchOrganizationsResponse.deserializeBinary
  );

  searchOrganizations(
    request: org_org_api_pb.SearchOrganizationsRequest,
    metadata: grpcWeb.Metadata | null): Promise<org_org_api_pb.SearchOrganizationsResponse>;

  searchOrganizations(
    request: org_org_api_pb.SearchOrganizationsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: org_org_api_pb.SearchOrganizationsResponse) => void): grpcWeb.ClientReadableStream<org_org_api_pb.SearchOrganizationsResponse>;

  searchOrganizations(
    request: org_org_api_pb.SearchOrganizationsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: org_org_api_pb.SearchOrganizationsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.org.v1.OrgAPI/SearchOrganizations',
        request,
        metadata || {},
        this.methodInfoSearchOrganizations,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.org.v1.OrgAPI/SearchOrganizations',
    request,
    metadata || {},
    this.methodInfoSearchOrganizations);
  }

  methodInfoDeleteOrganization = new grpcWeb.AbstractClientBase.MethodInfo(
    org_org_api_pb.DeleteOrganizationResponse,
    (request: org_org_api_pb.DeleteOrganizationRequest) => {
      return request.serializeBinary();
    },
    org_org_api_pb.DeleteOrganizationResponse.deserializeBinary
  );

  deleteOrganization(
    request: org_org_api_pb.DeleteOrganizationRequest,
    metadata: grpcWeb.Metadata | null): Promise<org_org_api_pb.DeleteOrganizationResponse>;

  deleteOrganization(
    request: org_org_api_pb.DeleteOrganizationRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: org_org_api_pb.DeleteOrganizationResponse) => void): grpcWeb.ClientReadableStream<org_org_api_pb.DeleteOrganizationResponse>;

  deleteOrganization(
    request: org_org_api_pb.DeleteOrganizationRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: org_org_api_pb.DeleteOrganizationResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.org.v1.OrgAPI/DeleteOrganization',
        request,
        metadata || {},
        this.methodInfoDeleteOrganization,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.org.v1.OrgAPI/DeleteOrganization',
    request,
    metadata || {},
    this.methodInfoDeleteOrganization);
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

  methodInfoSearchMembers = new grpcWeb.AbstractClientBase.MethodInfo(
    org_org_api_pb.SearchMembersResponse,
    (request: org_org_api_pb.SearchMembersRequest) => {
      return request.serializeBinary();
    },
    org_org_api_pb.SearchMembersResponse.deserializeBinary
  );

  searchMembers(
    request: org_org_api_pb.SearchMembersRequest,
    metadata: grpcWeb.Metadata | null): Promise<org_org_api_pb.SearchMembersResponse>;

  searchMembers(
    request: org_org_api_pb.SearchMembersRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: org_org_api_pb.SearchMembersResponse) => void): grpcWeb.ClientReadableStream<org_org_api_pb.SearchMembersResponse>;

  searchMembers(
    request: org_org_api_pb.SearchMembersRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: org_org_api_pb.SearchMembersResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.org.v1.OrgAPI/SearchMembers',
        request,
        metadata || {},
        this.methodInfoSearchMembers,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.org.v1.OrgAPI/SearchMembers',
    request,
    metadata || {},
    this.methodInfoSearchMembers);
  }

}

