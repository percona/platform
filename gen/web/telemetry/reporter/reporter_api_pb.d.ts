import * as jspb from 'google-protobuf'

import * as github_com_mwitkow_go$proto$validators_validator_pb from '../../github.com/mwitkow/go-proto-validators/validator_pb';
import * as telemetry_reporter_event_pb from '../../telemetry/reporter/event_pb';
import * as telemetry_events_pmm_server_uptime_event_pb from '../../telemetry/events/pmm/server_uptime_event_pb';
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb';
import * as google_api_annotations_pb from '../../google/api/annotations_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class ReportRequest extends jspb.Message {
  getEventsList(): Array<telemetry_reporter_event_pb.Event>;
  setEventsList(value: Array<telemetry_reporter_event_pb.Event>): ReportRequest;
  clearEventsList(): ReportRequest;
  addEvents(value?: telemetry_reporter_event_pb.Event, index?: number): telemetry_reporter_event_pb.Event;

  getMetricsList(): Array<ServerMetric>;
  setMetricsList(value: Array<ServerMetric>): ReportRequest;
  clearMetricsList(): ReportRequest;
  addMetrics(value?: ServerMetric, index?: number): ServerMetric;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReportRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReportRequest): ReportRequest.AsObject;
  static serializeBinaryToWriter(message: ReportRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReportRequest;
  static deserializeBinaryFromReader(message: ReportRequest, reader: jspb.BinaryReader): ReportRequest;
}

export namespace ReportRequest {
  export type AsObject = {
    eventsList: Array<telemetry_reporter_event_pb.Event.AsObject>,
    metricsList: Array<ServerMetric.AsObject>,
  }
}

export class ReportResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReportResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReportResponse): ReportResponse.AsObject;
  static serializeBinaryToWriter(message: ReportResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReportResponse;
  static deserializeBinaryFromReader(message: ReportResponse, reader: jspb.BinaryReader): ReportResponse;
}

export namespace ReportResponse {
  export type AsObject = {
  }
}

export class ServerMetric extends jspb.Message {
  getId(): Uint8Array | string;
  getId_asU8(): Uint8Array;
  getId_asB64(): string;
  setId(value: Uint8Array | string): ServerMetric;

  getTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setTime(value?: google_protobuf_timestamp_pb.Timestamp): ServerMetric;
  hasTime(): boolean;
  clearTime(): ServerMetric;

  getPmmServerTelemetryId(): Uint8Array | string;
  getPmmServerTelemetryId_asU8(): Uint8Array;
  getPmmServerTelemetryId_asB64(): string;
  setPmmServerTelemetryId(value: Uint8Array | string): ServerMetric;

  getPmmServerVersion(): string;
  setPmmServerVersion(value: string): ServerMetric;

  getUpDuration(): google_protobuf_duration_pb.Duration | undefined;
  setUpDuration(value?: google_protobuf_duration_pb.Duration): ServerMetric;
  hasUpDuration(): boolean;
  clearUpDuration(): ServerMetric;

  getDistributionMethod(): telemetry_events_pmm_server_uptime_event_pb.DistributionMethod;
  setDistributionMethod(value: telemetry_events_pmm_server_uptime_event_pb.DistributionMethod): ServerMetric;

  getMetricsList(): Array<ServerMetric.Metric>;
  setMetricsList(value: Array<ServerMetric.Metric>): ServerMetric;
  clearMetricsList(): ServerMetric;
  addMetrics(value?: ServerMetric.Metric, index?: number): ServerMetric.Metric;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerMetric.AsObject;
  static toObject(includeInstance: boolean, msg: ServerMetric): ServerMetric.AsObject;
  static serializeBinaryToWriter(message: ServerMetric, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerMetric;
  static deserializeBinaryFromReader(message: ServerMetric, reader: jspb.BinaryReader): ServerMetric;
}

export namespace ServerMetric {
  export type AsObject = {
    id: Uint8Array | string,
    time?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    pmmServerTelemetryId: Uint8Array | string,
    pmmServerVersion: string,
    upDuration?: google_protobuf_duration_pb.Duration.AsObject,
    distributionMethod: telemetry_events_pmm_server_uptime_event_pb.DistributionMethod,
    metricsList: Array<ServerMetric.Metric.AsObject>,
  }

  export class Metric extends jspb.Message {
    getKey(): string;
    setKey(value: string): Metric;

    getValue(): string;
    setValue(value: string): Metric;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Metric.AsObject;
    static toObject(includeInstance: boolean, msg: Metric): Metric.AsObject;
    static serializeBinaryToWriter(message: Metric, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Metric;
    static deserializeBinaryFromReader(message: Metric, reader: jspb.BinaryReader): Metric;
  }

  export namespace Metric {
    export type AsObject = {
      key: string,
      value: string,
    }
  }

}

