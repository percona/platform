import * as jspb from 'google-protobuf'

import * as github_com_mwitkow_go$proto$validators_validator_pb from '../github.com/mwitkow/go-proto-validators/validator_pb';
import * as google_api_annotations_pb from '../google/api/annotations_pb';


export class PerconaSSODetails extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): PerconaSSODetails;

  getClientSecret(): string;
  setClientSecret(value: string): PerconaSSODetails;

  getIssuerUrl(): string;
  setIssuerUrl(value: string): PerconaSSODetails;

  getScope(): string;
  setScope(value: string): PerconaSSODetails;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PerconaSSODetails.AsObject;
  static toObject(includeInstance: boolean, msg: PerconaSSODetails): PerconaSSODetails.AsObject;
  static serializeBinaryToWriter(message: PerconaSSODetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PerconaSSODetails;
  static deserializeBinaryFromReader(message: PerconaSSODetails, reader: jspb.BinaryReader): PerconaSSODetails;
}

export namespace PerconaSSODetails {
  export type AsObject = {
    clientId: string,
    clientSecret: string,
    issuerUrl: string,
    scope: string,
  }
}

export class AddPMMRequest extends jspb.Message {
  getPmmServerId(): string;
  setPmmServerId(value: string): AddPMMRequest;

  getPmmServerName(): string;
  setPmmServerName(value: string): AddPMMRequest;

  getPmmServerUrl(): string;
  setPmmServerUrl(value: string): AddPMMRequest;

  getPmmServerOauthCallbackUrl(): string;
  setPmmServerOauthCallbackUrl(value: string): AddPMMRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddPMMRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddPMMRequest): AddPMMRequest.AsObject;
  static serializeBinaryToWriter(message: AddPMMRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddPMMRequest;
  static deserializeBinaryFromReader(message: AddPMMRequest, reader: jspb.BinaryReader): AddPMMRequest;
}

export namespace AddPMMRequest {
  export type AsObject = {
    pmmServerId: string,
    pmmServerName: string,
    pmmServerUrl: string,
    pmmServerOauthCallbackUrl: string,
  }
}

export class AddPMMResponse extends jspb.Message {
  getSsoDetails(): PerconaSSODetails | undefined;
  setSsoDetails(value?: PerconaSSODetails): AddPMMResponse;
  hasSsoDetails(): boolean;
  clearSsoDetails(): AddPMMResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddPMMResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddPMMResponse): AddPMMResponse.AsObject;
  static serializeBinaryToWriter(message: AddPMMResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddPMMResponse;
  static deserializeBinaryFromReader(message: AddPMMResponse, reader: jspb.BinaryReader): AddPMMResponse;
}

export namespace AddPMMResponse {
  export type AsObject = {
    ssoDetails?: PerconaSSODetails.AsObject,
  }
}

