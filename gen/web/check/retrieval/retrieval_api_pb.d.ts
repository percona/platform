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

export class GetAllAlertRuleTemplatesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAllAlertRuleTemplatesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAllAlertRuleTemplatesRequest): GetAllAlertRuleTemplatesRequest.AsObject;
  static serializeBinaryToWriter(message: GetAllAlertRuleTemplatesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAllAlertRuleTemplatesRequest;
  static deserializeBinaryFromReader(message: GetAllAlertRuleTemplatesRequest, reader: jspb.BinaryReader): GetAllAlertRuleTemplatesRequest;
}

export namespace GetAllAlertRuleTemplatesRequest {
  export type AsObject = {
  }
}

export class GetAllAlertRuleTemplatesResponse extends jspb.Message {
  getFile(): string;
  setFile(value: string): GetAllAlertRuleTemplatesResponse;

  getSignaturesList(): Array<string>;
  setSignaturesList(value: Array<string>): GetAllAlertRuleTemplatesResponse;
  clearSignaturesList(): GetAllAlertRuleTemplatesResponse;
  addSignatures(value: string, index?: number): GetAllAlertRuleTemplatesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAllAlertRuleTemplatesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAllAlertRuleTemplatesResponse): GetAllAlertRuleTemplatesResponse.AsObject;
  static serializeBinaryToWriter(message: GetAllAlertRuleTemplatesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAllAlertRuleTemplatesResponse;
  static deserializeBinaryFromReader(message: GetAllAlertRuleTemplatesResponse, reader: jspb.BinaryReader): GetAllAlertRuleTemplatesResponse;
}

export namespace GetAllAlertRuleTemplatesResponse {
  export type AsObject = {
    file: string,
    signaturesList: Array<string>,
  }
}

