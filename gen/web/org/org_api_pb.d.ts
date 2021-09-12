import * as jspb from 'google-protobuf'

import * as github_com_mwitkow_go$proto$validators_validator_pb from '../github.com/mwitkow/go-proto-validators/validator_pb';
import * as google_api_annotations_pb from '../google/api/annotations_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class CreateOrganizationRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateOrganizationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateOrganizationRequest): CreateOrganizationRequest.AsObject;
  static serializeBinaryToWriter(message: CreateOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateOrganizationRequest;
  static deserializeBinaryFromReader(message: CreateOrganizationRequest, reader: jspb.BinaryReader): CreateOrganizationRequest;
}

export namespace CreateOrganizationRequest {
  export type AsObject = {
    name: string,
  }
}

export class CreateOrganizationResponse extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): CreateOrganizationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateOrganizationResponse): CreateOrganizationResponse.AsObject;
  static serializeBinaryToWriter(message: CreateOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateOrganizationResponse;
  static deserializeBinaryFromReader(message: CreateOrganizationResponse, reader: jspb.BinaryReader): CreateOrganizationResponse;
}

export namespace CreateOrganizationResponse {
  export type AsObject = {
    orgId: string,
  }
}

export class GetOrganizationRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): GetOrganizationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationRequest): GetOrganizationRequest.AsObject;
  static serializeBinaryToWriter(message: GetOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationRequest;
  static deserializeBinaryFromReader(message: GetOrganizationRequest, reader: jspb.BinaryReader): GetOrganizationRequest;
}

export namespace GetOrganizationRequest {
  export type AsObject = {
    orgId: string,
  }
}

export class GetOrganizationResponse extends jspb.Message {
  getName(): string;
  setName(value: string): GetOrganizationResponse;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): GetOrganizationResponse;
  hasCreatedAt(): boolean;
  clearCreatedAt(): GetOrganizationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationResponse): GetOrganizationResponse.AsObject;
  static serializeBinaryToWriter(message: GetOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationResponse;
  static deserializeBinaryFromReader(message: GetOrganizationResponse, reader: jspb.BinaryReader): GetOrganizationResponse;
}

export namespace GetOrganizationResponse {
  export type AsObject = {
    name: string,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class GetOrganizationByUserRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationByUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationByUserRequest): GetOrganizationByUserRequest.AsObject;
  static serializeBinaryToWriter(message: GetOrganizationByUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationByUserRequest;
  static deserializeBinaryFromReader(message: GetOrganizationByUserRequest, reader: jspb.BinaryReader): GetOrganizationByUserRequest;
}

export namespace GetOrganizationByUserRequest {
  export type AsObject = {
  }
}

export class GetOrganizationByUserResponse extends jspb.Message {
  getOrgIdsList(): Array<string>;
  setOrgIdsList(value: Array<string>): GetOrganizationByUserResponse;
  clearOrgIdsList(): GetOrganizationByUserResponse;
  addOrgIds(value: string, index?: number): GetOrganizationByUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationByUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationByUserResponse): GetOrganizationByUserResponse.AsObject;
  static serializeBinaryToWriter(message: GetOrganizationByUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationByUserResponse;
  static deserializeBinaryFromReader(message: GetOrganizationByUserResponse, reader: jspb.BinaryReader): GetOrganizationByUserResponse;
}

export namespace GetOrganizationByUserResponse {
  export type AsObject = {
    orgIdsList: Array<string>,
  }
}

export class InviteMemberRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): InviteMemberRequest;

  getRole(): OrganizationMemberRole;
  setRole(value: OrganizationMemberRole): InviteMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InviteMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: InviteMemberRequest): InviteMemberRequest.AsObject;
  static serializeBinaryToWriter(message: InviteMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InviteMemberRequest;
  static deserializeBinaryFromReader(message: InviteMemberRequest, reader: jspb.BinaryReader): InviteMemberRequest;
}

export namespace InviteMemberRequest {
  export type AsObject = {
    email: string,
    role: OrganizationMemberRole,
  }
}

export class InviteMemberResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InviteMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: InviteMemberResponse): InviteMemberResponse.AsObject;
  static serializeBinaryToWriter(message: InviteMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InviteMemberResponse;
  static deserializeBinaryFromReader(message: InviteMemberResponse, reader: jspb.BinaryReader): InviteMemberResponse;
}

export namespace InviteMemberResponse {
  export type AsObject = {
  }
}

export class OrganizationMember extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): OrganizationMember;

  getName(): string;
  setName(value: string): OrganizationMember;

  getRole(): OrganizationMemberRole;
  setRole(value: OrganizationMemberRole): OrganizationMember;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationMember.AsObject;
  static toObject(includeInstance: boolean, msg: OrganizationMember): OrganizationMember.AsObject;
  static serializeBinaryToWriter(message: OrganizationMember, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationMember;
  static deserializeBinaryFromReader(message: OrganizationMember, reader: jspb.BinaryReader): OrganizationMember;
}

export namespace OrganizationMember {
  export type AsObject = {
    email: string,
    name: string,
    role: OrganizationMemberRole,
  }
}

export class ListMembersRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMembersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMembersRequest): ListMembersRequest.AsObject;
  static serializeBinaryToWriter(message: ListMembersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMembersRequest;
  static deserializeBinaryFromReader(message: ListMembersRequest, reader: jspb.BinaryReader): ListMembersRequest;
}

export namespace ListMembersRequest {
  export type AsObject = {
  }
}

export class ListMembersResponse extends jspb.Message {
  getMembersList(): Array<OrganizationMember>;
  setMembersList(value: Array<OrganizationMember>): ListMembersResponse;
  clearMembersList(): ListMembersResponse;
  addMembers(value?: OrganizationMember, index?: number): OrganizationMember;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMembersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMembersResponse): ListMembersResponse.AsObject;
  static serializeBinaryToWriter(message: ListMembersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMembersResponse;
  static deserializeBinaryFromReader(message: ListMembersResponse, reader: jspb.BinaryReader): ListMembersResponse;
}

export namespace ListMembersResponse {
  export type AsObject = {
    membersList: Array<OrganizationMember.AsObject>,
  }
}

export enum OrganizationMemberRole { 
  ORGANIZATION_MEMBER_ROLE_INVALID = 0,
  ADMIN = 1,
  TECHNICAL = 2,
}
