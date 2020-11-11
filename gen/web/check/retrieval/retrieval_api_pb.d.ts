import * as jspb from 'google-protobuf'



export class GetAllChecksRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAllChecksRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAllChecksRequest): GetAllChecksRequest.AsObject;
  static serializeBinaryToWriter(message: GetAllChecksRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAllChecksRequest;
  static deserializeBinaryFromReader(message: GetAllChecksRequest, reader: jspb.BinaryReader): GetAllChecksRequest;
}

export namespace GetAllChecksRequest {
  export type AsObject = {
  }
}

export class GetAllChecksResponse extends jspb.Message {
  getFile(): string;
  setFile(value: string): GetAllChecksResponse;

  getSignaturesList(): Array<string>;
  setSignaturesList(value: Array<string>): GetAllChecksResponse;
  clearSignaturesList(): GetAllChecksResponse;
  addSignatures(value: string, index?: number): GetAllChecksResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAllChecksResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAllChecksResponse): GetAllChecksResponse.AsObject;
  static serializeBinaryToWriter(message: GetAllChecksResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAllChecksResponse;
  static deserializeBinaryFromReader(message: GetAllChecksResponse, reader: jspb.BinaryReader): GetAllChecksResponse;
}

export namespace GetAllChecksResponse {
  export type AsObject = {
    file: string,
    signaturesList: Array<string>,
  }
}

export class GetAllRulesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAllRulesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAllRulesRequest): GetAllRulesRequest.AsObject;
  static serializeBinaryToWriter(message: GetAllRulesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAllRulesRequest;
  static deserializeBinaryFromReader(message: GetAllRulesRequest, reader: jspb.BinaryReader): GetAllRulesRequest;
}

export namespace GetAllRulesRequest {
  export type AsObject = {
  }
}

export class GetAllRulesResponse extends jspb.Message {
  getFile(): string;
  setFile(value: string): GetAllRulesResponse;

  getSignaturesList(): Array<string>;
  setSignaturesList(value: Array<string>): GetAllRulesResponse;
  clearSignaturesList(): GetAllRulesResponse;
  addSignatures(value: string, index?: number): GetAllRulesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAllRulesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAllRulesResponse): GetAllRulesResponse.AsObject;
  static serializeBinaryToWriter(message: GetAllRulesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAllRulesResponse;
  static deserializeBinaryFromReader(message: GetAllRulesResponse, reader: jspb.BinaryReader): GetAllRulesResponse;
}

export namespace GetAllRulesResponse {
  export type AsObject = {
    file: string,
    signaturesList: Array<string>,
  }
}

