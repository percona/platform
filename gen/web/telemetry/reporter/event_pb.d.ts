import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class AnyEvent extends jspb.Message {
  getTypeUrl(): string;
  setTypeUrl(value: string): AnyEvent;

  getBinary(): Uint8Array | string;
  getBinary_asU8(): Uint8Array;
  getBinary_asB64(): string;
  setBinary(value: Uint8Array | string): AnyEvent;

  getJson(): string;
  setJson(value: string): AnyEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AnyEvent.AsObject;
  static toObject(includeInstance: boolean, msg: AnyEvent): AnyEvent.AsObject;
  static serializeBinaryToWriter(message: AnyEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AnyEvent;
  static deserializeBinaryFromReader(message: AnyEvent, reader: jspb.BinaryReader): AnyEvent;
}

export namespace AnyEvent {
  export type AsObject = {
    typeUrl: string,
    binary: Uint8Array | string,
    json: string,
  }
}

export class Event extends jspb.Message {
  getId(): Uint8Array | string;
  getId_asU8(): Uint8Array;
  getId_asB64(): string;
  setId(value: Uint8Array | string): Event;

  getTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setTime(value?: google_protobuf_timestamp_pb.Timestamp): Event;
  hasTime(): boolean;
  clearTime(): Event;

  getEvent(): AnyEvent | undefined;
  setEvent(value?: AnyEvent): Event;
  hasEvent(): boolean;
  clearEvent(): Event;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Event.AsObject;
  static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
  static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Event;
  static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
}

export namespace Event {
  export type AsObject = {
    id: Uint8Array | string,
    time?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    event?: AnyEvent.AsObject,
  }
}

