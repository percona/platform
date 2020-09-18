/**
 * @fileoverview gRPC-Web generated client stub for percona.platform.telemetry.reporter.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as telemetry_reporter_reporter_api_pb from '../../telemetry/reporter/reporter_api_pb';


export class ReporterAPIClient {
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

  methodInfoReport = new grpcWeb.AbstractClientBase.MethodInfo(
    telemetry_reporter_reporter_api_pb.ReportResponse,
    (request: telemetry_reporter_reporter_api_pb.ReportRequest) => {
      return request.serializeBinary();
    },
    telemetry_reporter_reporter_api_pb.ReportResponse.deserializeBinary
  );

  report(
    request: telemetry_reporter_reporter_api_pb.ReportRequest,
    metadata: grpcWeb.Metadata | null): Promise<telemetry_reporter_reporter_api_pb.ReportResponse>;

  report(
    request: telemetry_reporter_reporter_api_pb.ReportRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: telemetry_reporter_reporter_api_pb.ReportResponse) => void): grpcWeb.ClientReadableStream<telemetry_reporter_reporter_api_pb.ReportResponse>;

  report(
    request: telemetry_reporter_reporter_api_pb.ReportRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: telemetry_reporter_reporter_api_pb.ReportResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/percona.platform.telemetry.reporter.v1.ReporterAPI/Report',
        request,
        metadata || {},
        this.methodInfoReport,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/percona.platform.telemetry.reporter.v1.ReporterAPI/Report',
    request,
    metadata || {},
    this.methodInfoReport);
  }

}

