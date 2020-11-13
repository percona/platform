/**
 * @fileoverview gRPC-Web generated client stub for percona.platform.check.retrieval.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as check_retrieval_retrieval_api_pb from '../../check/retrieval/retrieval_api_pb';


export class RetrievalAPIClient {
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

  methodInfoGetAllChecks = new grpcWeb.AbstractClientBase.MethodInfo(
    check_retrieval_retrieval_api_pb.GetAllChecksResponse,
    (request: check_retrieval_retrieval_api_pb.GetAllChecksRequest) => {
      return request.serializeBinary();
    },
    check_retrieval_retrieval_api_pb.GetAllChecksResponse.deserializeBinary
  );

  getAllChecks(
    request: check_retrieval_retrieval_api_pb.GetAllChecksRequest,
    metadata: grpcWeb.Metadata | null): Promise<check_retrieval_retrieval_api_pb.GetAllChecksResponse>;

  getAllChecks(
    request: check_retrieval_retrieval_api_pb.GetAllChecksRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: check_retrieval_retrieval_api_pb.GetAllChecksResponse) => void): grpcWeb.ClientReadableStream<check_retrieval_retrieval_api_pb.GetAllChecksResponse>;

  getAllChecks(
    request: check_retrieval_retrieval_api_pb.GetAllChecksRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: check_retrieval_retrieval_api_pb.GetAllChecksResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.check.retrieval.v1.RetrievalAPI/GetAllChecks',
        request,
        metadata || {},
        this.methodInfoGetAllChecks,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.check.retrieval.v1.RetrievalAPI/GetAllChecks',
    request,
    metadata || {},
    this.methodInfoGetAllChecks);
  }

  methodInfoGetAllAlertRuleTemplates = new grpcWeb.AbstractClientBase.MethodInfo(
    check_retrieval_retrieval_api_pb.GetAllAlertRuleTemplatesResponse,
    (request: check_retrieval_retrieval_api_pb.GetAllAlertRuleTemplatesRequest) => {
      return request.serializeBinary();
    },
    check_retrieval_retrieval_api_pb.GetAllAlertRuleTemplatesResponse.deserializeBinary
  );

  getAllAlertRuleTemplates(
    request: check_retrieval_retrieval_api_pb.GetAllAlertRuleTemplatesRequest,
    metadata: grpcWeb.Metadata | null): Promise<check_retrieval_retrieval_api_pb.GetAllAlertRuleTemplatesResponse>;

  getAllAlertRuleTemplates(
    request: check_retrieval_retrieval_api_pb.GetAllAlertRuleTemplatesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: check_retrieval_retrieval_api_pb.GetAllAlertRuleTemplatesResponse) => void): grpcWeb.ClientReadableStream<check_retrieval_retrieval_api_pb.GetAllAlertRuleTemplatesResponse>;

  getAllAlertRuleTemplates(
    request: check_retrieval_retrieval_api_pb.GetAllAlertRuleTemplatesRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: check_retrieval_retrieval_api_pb.GetAllAlertRuleTemplatesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.check.retrieval.v1.RetrievalAPI/GetAllAlertRuleTemplates',
        request,
        metadata || {},
        this.methodInfoGetAllAlertRuleTemplates,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.check.retrieval.v1.RetrievalAPI/GetAllAlertRuleTemplates',
    request,
    metadata || {},
    this.methodInfoGetAllAlertRuleTemplates);
  }

}

