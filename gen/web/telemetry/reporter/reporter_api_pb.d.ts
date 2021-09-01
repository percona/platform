import * as jspb from 'google-protobuf'

import * as github_com_mwitkow_go$proto$validators_validator_pb from '../../github.com/mwitkow/go-proto-validators/validator_pb';
import * as telemetry_reporter_event_pb from '../../telemetry/reporter/event_pb';
import * as google_api_annotations_pb from '../../google/api/annotations_pb';


export class ReportRequest extends jspb.Message {
  getEventsList(): Array<telemetry_reporter_event_pb.Event>;
  setEventsList(value: Array<telemetry_reporter_event_pb.Event>): ReportRequest;
  clearEventsList(): ReportRequest;
  addEvents(value?: telemetry_reporter_event_pb.Event, index?: number): telemetry_reporter_event_pb.Event;

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

