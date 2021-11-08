/**
 * @fileoverview gRPC-Web generated client stub for percona.platform.org.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as org_inventory_api_pb from '../org/inventory_api_pb';


export class InventoryAPIClient {
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

  methodInfoAddPMM = new grpcWeb.AbstractClientBase.MethodInfo(
    org_inventory_api_pb.AddPMMResponse,
    (request: org_inventory_api_pb.AddPMMRequest) => {
      return request.serializeBinary();
    },
    org_inventory_api_pb.AddPMMResponse.deserializeBinary
  );

  addPMM(
    request: org_inventory_api_pb.AddPMMRequest,
    metadata: grpcWeb.Metadata | null): Promise<org_inventory_api_pb.AddPMMResponse>;

  addPMM(
    request: org_inventory_api_pb.AddPMMRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: org_inventory_api_pb.AddPMMResponse) => void): grpcWeb.ClientReadableStream<org_inventory_api_pb.AddPMMResponse>;

  addPMM(
    request: org_inventory_api_pb.AddPMMRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: org_inventory_api_pb.AddPMMResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.org.v1.InventoryAPI/AddPMM',
        request,
        metadata || {},
        this.methodInfoAddPMM,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.org.v1.InventoryAPI/AddPMM',
    request,
    metadata || {},
    this.methodInfoAddPMM);
  }

}

