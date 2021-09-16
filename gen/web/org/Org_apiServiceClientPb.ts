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

  methodInfoGetOrganizationByID = new grpcWeb.AbstractClientBase.MethodInfo(
    org_org_api_pb.GetOrganizationByIDResponse,
    (request: org_org_api_pb.GetOrganizationByIDRequest) => {
      return request.serializeBinary();
    },
    org_org_api_pb.GetOrganizationByIDResponse.deserializeBinary
  );

  getOrganizationByID(
    request: org_org_api_pb.GetOrganizationByIDRequest,
    metadata: grpcWeb.Metadata | null): Promise<org_org_api_pb.GetOrganizationByIDResponse>;

  getOrganizationByID(
    request: org_org_api_pb.GetOrganizationByIDRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: org_org_api_pb.GetOrganizationByIDResponse) => void): grpcWeb.ClientReadableStream<org_org_api_pb.GetOrganizationByIDResponse>;

  getOrganizationByID(
    request: org_org_api_pb.GetOrganizationByIDRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: org_org_api_pb.GetOrganizationByIDResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.org.v1.OrgAPI/GetOrganizationByID',
        request,
        metadata || {},
        this.methodInfoGetOrganizationByID,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.org.v1.OrgAPI/GetOrganizationByID',
    request,
    metadata || {},
    this.methodInfoGetOrganizationByID);
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

}

