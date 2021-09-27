import * as jspb from 'google-protobuf'

import * as github_com_mwitkow_go$proto$validators_validator_pb from '../../../github.com/mwitkow/go-proto-validators/validator_pb';
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb';
import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';


export class ServerUptimeEvent extends jspb.Message {
  getId(): Uint8Array | string;
  getId_asU8(): Uint8Array;
  getId_asB64(): string;
  setId(value: Uint8Array | string): ServerUptimeEvent;

  getVersion(): string;
  setVersion(value: string): ServerUptimeEvent;

  getUpDuration(): google_protobuf_duration_pb.Duration | undefined;
  setUpDuration(value?: google_protobuf_duration_pb.Duration): ServerUptimeEvent;
  hasUpDuration(): boolean;
  clearUpDuration(): ServerUptimeEvent;

  getDistributionMethod(): DistributionMethod;
  setDistributionMethod(value: DistributionMethod): ServerUptimeEvent;

  getSttEnabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setSttEnabled(value?: google_protobuf_wrappers_pb.BoolValue): ServerUptimeEvent;
  hasSttEnabled(): boolean;
  clearSttEnabled(): ServerUptimeEvent;

  getIaEnabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setIaEnabled(value?: google_protobuf_wrappers_pb.BoolValue): ServerUptimeEvent;
  hasIaEnabled(): boolean;
  clearIaEnabled(): ServerUptimeEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerUptimeEvent.AsObject;
  static toObject(includeInstance: boolean, msg: ServerUptimeEvent): ServerUptimeEvent.AsObject;
  static serializeBinaryToWriter(message: ServerUptimeEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerUptimeEvent;
  static deserializeBinaryFromReader(message: ServerUptimeEvent, reader: jspb.BinaryReader): ServerUptimeEvent;
}

export namespace ServerUptimeEvent {
  export type AsObject = {
    id: Uint8Array | string,
    version: string,
    upDuration?: google_protobuf_duration_pb.Duration.AsObject,
    distributionMethod: DistributionMethod,
    sttEnabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    iaEnabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export enum DistributionMethod { 
  DISTRIBUTION_METHOD_INVALID = 0,
  DOCKER = 1,
  OVF = 2,
  AMI = 3,
  AZURE = 4,
  DO = 5,
}
