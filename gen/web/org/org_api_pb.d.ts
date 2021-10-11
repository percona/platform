import * as jspb from 'google-protobuf'

import * as github_com_mwitkow_go$proto$validators_validator_pb from '../github.com/mwitkow/go-proto-validators/validator_pb';
import * as google_api_annotations_pb from '../google/api/annotations_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class Organization extends jspb.Message {
  getId(): string;
  setId(value: string): Organization;

  getName(): string;
  setName(value: string): Organization;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Organization;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Organization;

  getUpdatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Organization;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Organization;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Organization.AsObject;
  static toObject(includeInstance: boolean, msg: Organization): Organization.AsObject;
  static serializeBinaryToWriter(message: Organization, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Organization;
  static deserializeBinaryFromReader(message: Organization, reader: jspb.BinaryReader): Organization;
}

export namespace Organization {
  export type AsObject = {
    id: string,
    name: string,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    updatedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

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
  getOrg(): Organization | undefined;
  setOrg(value?: Organization): CreateOrganizationResponse;
  hasOrg(): boolean;
  clearOrg(): CreateOrganizationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateOrganizationResponse): CreateOrganizationResponse.AsObject;
  static serializeBinaryToWriter(message: CreateOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateOrganizationResponse;
  static deserializeBinaryFromReader(message: CreateOrganizationResponse, reader: jspb.BinaryReader): CreateOrganizationResponse;
}

export namespace CreateOrganizationResponse {
  export type AsObject = {
    org?: Organization.AsObject,
  }
}

export class GetOrganizationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetOrganizationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationRequest): GetOrganizationRequest.AsObject;
  static serializeBinaryToWriter(message: GetOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationRequest;
  static deserializeBinaryFromReader(message: GetOrganizationRequest, reader: jspb.BinaryReader): GetOrganizationRequest;
}

export namespace GetOrganizationRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetOrganizationResponse extends jspb.Message {
  getOrg(): Organization | undefined;
  setOrg(value?: Organization): GetOrganizationResponse;
  hasOrg(): boolean;
  clearOrg(): GetOrganizationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationResponse): GetOrganizationResponse.AsObject;
  static serializeBinaryToWriter(message: GetOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationResponse;
  static deserializeBinaryFromReader(message: GetOrganizationResponse, reader: jspb.BinaryReader): GetOrganizationResponse;
}

export namespace GetOrganizationResponse {
  export type AsObject = {
    org?: Organization.AsObject,
  }
}

export class SearchOrganizationsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchOrganizationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SearchOrganizationsRequest): SearchOrganizationsRequest.AsObject;
  static serializeBinaryToWriter(message: SearchOrganizationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchOrganizationsRequest;
  static deserializeBinaryFromReader(message: SearchOrganizationsRequest, reader: jspb.BinaryReader): SearchOrganizationsRequest;
}

export namespace SearchOrganizationsRequest {
  export type AsObject = {
  }
}

export class SearchOrganizationsResponse extends jspb.Message {
  getOrgsList(): Array<Organization>;
  setOrgsList(value: Array<Organization>): SearchOrganizationsResponse;
  clearOrgsList(): SearchOrganizationsResponse;
  addOrgs(value?: Organization, index?: number): Organization;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchOrganizationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SearchOrganizationsResponse): SearchOrganizationsResponse.AsObject;
  static serializeBinaryToWriter(message: SearchOrganizationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchOrganizationsResponse;
  static deserializeBinaryFromReader(message: SearchOrganizationsResponse, reader: jspb.BinaryReader): SearchOrganizationsResponse;
}

export namespace SearchOrganizationsResponse {
  export type AsObject = {
    orgsList: Array<Organization.AsObject>,
  }
}

export class DeleteOrganizationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteOrganizationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteOrganizationRequest): DeleteOrganizationRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteOrganizationRequest;
  static deserializeBinaryFromReader(message: DeleteOrganizationRequest, reader: jspb.BinaryReader): DeleteOrganizationRequest;
}

export namespace DeleteOrganizationRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeleteOrganizationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteOrganizationResponse): DeleteOrganizationResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteOrganizationResponse;
  static deserializeBinaryFromReader(message: DeleteOrganizationResponse, reader: jspb.BinaryReader): DeleteOrganizationResponse;
}

export namespace DeleteOrganizationResponse {
  export type AsObject = {
  }
}

export class InviteMemberRequest extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): InviteMemberRequest;

  getId(): string;
  setId(value: string): InviteMemberRequest;

  getRole(): string;
  setRole(value: string): InviteMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InviteMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: InviteMemberRequest): InviteMemberRequest.AsObject;
  static serializeBinaryToWriter(message: InviteMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InviteMemberRequest;
  static deserializeBinaryFromReader(message: InviteMemberRequest, reader: jspb.BinaryReader): InviteMemberRequest;
}

export namespace InviteMemberRequest {
  export type AsObject = {
    username: string,
    id: string,
    role: string,
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
  getUsername(): string;
  setUsername(value: string): OrganizationMember;

  getFirstName(): string;
  setFirstName(value: string): OrganizationMember;

  getLastName(): string;
  setLastName(value: string): OrganizationMember;

  getRole(): string;
  setRole(value: string): OrganizationMember;

  getStatus(): string;
  setStatus(value: string): OrganizationMember;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationMember.AsObject;
  static toObject(includeInstance: boolean, msg: OrganizationMember): OrganizationMember.AsObject;
  static serializeBinaryToWriter(message: OrganizationMember, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationMember;
  static deserializeBinaryFromReader(message: OrganizationMember, reader: jspb.BinaryReader): OrganizationMember;
}

export namespace OrganizationMember {
  export type AsObject = {
    username: string,
    firstName: string,
    lastName: string,
    role: string,
    status: string,
  }
}

export class SearchMembersRequest extends jspb.Message {
  getId(): string;
  setId(value: string): SearchMembersRequest;

  getUser(): SearchMembersRequest.UserFilter | undefined;
  setUser(value?: SearchMembersRequest.UserFilter): SearchMembersRequest;
  hasUser(): boolean;
  clearUser(): SearchMembersRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchMembersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SearchMembersRequest): SearchMembersRequest.AsObject;
  static serializeBinaryToWriter(message: SearchMembersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchMembersRequest;
  static deserializeBinaryFromReader(message: SearchMembersRequest, reader: jspb.BinaryReader): SearchMembersRequest;
}

export namespace SearchMembersRequest {
  export type AsObject = {
    id: string,
    user?: SearchMembersRequest.UserFilter.AsObject,
  }

  export class UserFilter extends jspb.Message {
    getUsername(): string;
    setUsername(value: string): UserFilter;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UserFilter.AsObject;
    static toObject(includeInstance: boolean, msg: UserFilter): UserFilter.AsObject;
    static serializeBinaryToWriter(message: UserFilter, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UserFilter;
    static deserializeBinaryFromReader(message: UserFilter, reader: jspb.BinaryReader): UserFilter;
  }

  export namespace UserFilter {
    export type AsObject = {
      username: string,
    }
  }

}

export class SearchMembersResponse extends jspb.Message {
  getMembersList(): Array<OrganizationMember>;
  setMembersList(value: Array<OrganizationMember>): SearchMembersResponse;
  clearMembersList(): SearchMembersResponse;
  addMembers(value?: OrganizationMember, index?: number): OrganizationMember;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchMembersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SearchMembersResponse): SearchMembersResponse.AsObject;
  static serializeBinaryToWriter(message: SearchMembersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchMembersResponse;
  static deserializeBinaryFromReader(message: SearchMembersResponse, reader: jspb.BinaryReader): SearchMembersResponse;
}

export namespace SearchMembersResponse {
  export type AsObject = {
    membersList: Array<OrganizationMember.AsObject>,
  }
}

