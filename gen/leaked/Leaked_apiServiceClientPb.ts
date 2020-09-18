/**
 * @fileoverview gRPC-Web generated client stub for percona.platform.leaked.v1beta1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as leaked_leaked_api_pb from '../leaked/leaked_api_pb';


export class LeakedAPIClient {
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

  methodInfoCheckDoubleSHA1 = new grpcWeb.AbstractClientBase.MethodInfo(
    leaked_leaked_api_pb.CheckDoubleSHA1Response,
    (request: leaked_leaked_api_pb.CheckDoubleSHA1Request) => {
      return request.serializeBinary();
    },
    leaked_leaked_api_pb.CheckDoubleSHA1Response.deserializeBinary
  );

  checkDoubleSHA1(
    request: leaked_leaked_api_pb.CheckDoubleSHA1Request,
    metadata: grpcWeb.Metadata | null): Promise<leaked_leaked_api_pb.CheckDoubleSHA1Response>;

  checkDoubleSHA1(
    request: leaked_leaked_api_pb.CheckDoubleSHA1Request,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: leaked_leaked_api_pb.CheckDoubleSHA1Response) => void): grpcWeb.ClientReadableStream<leaked_leaked_api_pb.CheckDoubleSHA1Response>;

  checkDoubleSHA1(
    request: leaked_leaked_api_pb.CheckDoubleSHA1Request,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: leaked_leaked_api_pb.CheckDoubleSHA1Response) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.leaked.v1beta1.LeakedAPI/CheckDoubleSHA1',
        request,
        metadata || {},
        this.methodInfoCheckDoubleSHA1,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.leaked.v1beta1.LeakedAPI/CheckDoubleSHA1',
    request,
    metadata || {},
    this.methodInfoCheckDoubleSHA1);
  }

}

