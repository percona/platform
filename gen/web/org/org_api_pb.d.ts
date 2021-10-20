import * as jspb from 'google-protobuf'

import * as github_com_mwitkow_go$proto$validators_validator_pb from '../github.com/mwitkow/go-proto-validators/validator_pb';
import * as google_api_annotations_pb from '../google/api/annotations_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';


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

export class SearchOrganizationEntitlementsRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): SearchOrganizationEntitlementsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchOrganizationEntitlementsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SearchOrganizationEntitlementsRequest): SearchOrganizationEntitlementsRequest.AsObject;
  static serializeBinaryToWriter(message: SearchOrganizationEntitlementsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchOrganizationEntitlementsRequest;
  static deserializeBinaryFromReader(message: SearchOrganizationEntitlementsRequest, reader: jspb.BinaryReader): SearchOrganizationEntitlementsRequest;
}

export namespace SearchOrganizationEntitlementsRequest {
  export type AsObject = {
    orgId: string,
  }
}

export class SearchOrganizationEntitlementsResponse extends jspb.Message {
  getEntitlementsList(): Array<OrganizationEntitlement>;
  setEntitlementsList(value: Array<OrganizationEntitlement>): SearchOrganizationEntitlementsResponse;
  clearEntitlementsList(): SearchOrganizationEntitlementsResponse;
  addEntitlements(value?: OrganizationEntitlement, index?: number): OrganizationEntitlement;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchOrganizationEntitlementsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SearchOrganizationEntitlementsResponse): SearchOrganizationEntitlementsResponse.AsObject;
  static serializeBinaryToWriter(message: SearchOrganizationEntitlementsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchOrganizationEntitlementsResponse;
  static deserializeBinaryFromReader(message: SearchOrganizationEntitlementsResponse, reader: jspb.BinaryReader): SearchOrganizationEntitlementsResponse;
}

export namespace SearchOrganizationEntitlementsResponse {
  export type AsObject = {
    entitlementsList: Array<OrganizationEntitlement.AsObject>,
  }
}

export class OrganizationEntitlement extends jspb.Message {
  getNumber(): string;
  setNumber(value: string): OrganizationEntitlement;

  getName(): string;
  setName(value: string): OrganizationEntitlement;

  getSummary(): string;
  setSummary(value: string): OrganizationEntitlement;

  getTier(): google_protobuf_wrappers_pb.StringValue | undefined;
  setTier(value?: google_protobuf_wrappers_pb.StringValue): OrganizationEntitlement;
  hasTier(): boolean;
  clearTier(): OrganizationEntitlement;

  getTotalUnits(): google_protobuf_wrappers_pb.StringValue | undefined;
  setTotalUnits(value?: google_protobuf_wrappers_pb.StringValue): OrganizationEntitlement;
  hasTotalUnits(): boolean;
  clearTotalUnits(): OrganizationEntitlement;

  getUnlimitedUnits(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUnlimitedUnits(value?: google_protobuf_wrappers_pb.BoolValue): OrganizationEntitlement;
  hasUnlimitedUnits(): boolean;
  clearUnlimitedUnits(): OrganizationEntitlement;

  getSupportLevel(): google_protobuf_wrappers_pb.StringValue | undefined;
  setSupportLevel(value?: google_protobuf_wrappers_pb.StringValue): OrganizationEntitlement;
  hasSupportLevel(): boolean;
  clearSupportLevel(): OrganizationEntitlement;

  getSoftwareFamiliesList(): Array<string>;
  setSoftwareFamiliesList(value: Array<string>): OrganizationEntitlement;
  clearSoftwareFamiliesList(): OrganizationEntitlement;
  addSoftwareFamilies(value: string, index?: number): OrganizationEntitlement;

  getStartDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setStartDate(value?: google_protobuf_timestamp_pb.Timestamp): OrganizationEntitlement;
  hasStartDate(): boolean;
  clearStartDate(): OrganizationEntitlement;

  getEndDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setEndDate(value?: google_protobuf_timestamp_pb.Timestamp): OrganizationEntitlement;
  hasEndDate(): boolean;
  clearEndDate(): OrganizationEntitlement;

  getPlatform(): OrganizationEntitlement.Platform | undefined;
  setPlatform(value?: OrganizationEntitlement.Platform): OrganizationEntitlement;
  hasPlatform(): boolean;
  clearPlatform(): OrganizationEntitlement;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationEntitlement.AsObject;
  static toObject(includeInstance: boolean, msg: OrganizationEntitlement): OrganizationEntitlement.AsObject;
  static serializeBinaryToWriter(message: OrganizationEntitlement, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationEntitlement;
  static deserializeBinaryFromReader(message: OrganizationEntitlement, reader: jspb.BinaryReader): OrganizationEntitlement;
}

export namespace OrganizationEntitlement {
  export type AsObject = {
    number: string,
    name: string,
    summary: string,
    tier?: google_protobuf_wrappers_pb.StringValue.AsObject,
    totalUnits?: google_protobuf_wrappers_pb.StringValue.AsObject,
    unlimitedUnits?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    supportLevel?: google_protobuf_wrappers_pb.StringValue.AsObject,
    softwareFamiliesList: Array<string>,
    startDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    endDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    platform?: OrganizationEntitlement.Platform.AsObject,
  }

  export class Platform extends jspb.Message {
    getSecurityAdvisor(): google_protobuf_wrappers_pb.StringValue | undefined;
    setSecurityAdvisor(value?: google_protobuf_wrappers_pb.StringValue): Platform;
    hasSecurityAdvisor(): boolean;
    clearSecurityAdvisor(): Platform;

    getConfigAdvisor(): google_protobuf_wrappers_pb.StringValue | undefined;
    setConfigAdvisor(value?: google_protobuf_wrappers_pb.StringValue): Platform;
    hasConfigAdvisor(): boolean;
    clearConfigAdvisor(): Platform;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Platform.AsObject;
    static toObject(includeInstance: boolean, msg: Platform): Platform.AsObject;
    static serializeBinaryToWriter(message: Platform, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Platform;
    static deserializeBinaryFromReader(message: Platform, reader: jspb.BinaryReader): Platform;
  }

  export namespace Platform {
    export type AsObject = {
      securityAdvisor?: google_protobuf_wrappers_pb.StringValue.AsObject,
      configAdvisor?: google_protobuf_wrappers_pb.StringValue.AsObject,
    }
  }

}

export class SearchOrganizationTicketsRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): SearchOrganizationTicketsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchOrganizationTicketsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SearchOrganizationTicketsRequest): SearchOrganizationTicketsRequest.AsObject;
  static serializeBinaryToWriter(message: SearchOrganizationTicketsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchOrganizationTicketsRequest;
  static deserializeBinaryFromReader(message: SearchOrganizationTicketsRequest, reader: jspb.BinaryReader): SearchOrganizationTicketsRequest;
}

export namespace SearchOrganizationTicketsRequest {
  export type AsObject = {
    orgId: string,
  }
}

export class SearchOrganizationTicketsResponse extends jspb.Message {
  getTicketsList(): Array<OrganizationTicket>;
  setTicketsList(value: Array<OrganizationTicket>): SearchOrganizationTicketsResponse;
  clearTicketsList(): SearchOrganizationTicketsResponse;
  addTickets(value?: OrganizationTicket, index?: number): OrganizationTicket;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchOrganizationTicketsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SearchOrganizationTicketsResponse): SearchOrganizationTicketsResponse.AsObject;
  static serializeBinaryToWriter(message: SearchOrganizationTicketsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchOrganizationTicketsResponse;
  static deserializeBinaryFromReader(message: SearchOrganizationTicketsResponse, reader: jspb.BinaryReader): SearchOrganizationTicketsResponse;
}

export namespace SearchOrganizationTicketsResponse {
  export type AsObject = {
    ticketsList: Array<OrganizationTicket.AsObject>,
  }
}

export class OrganizationTicket extends jspb.Message {
  getNumber(): string;
  setNumber(value: string): OrganizationTicket;

  getShortDescription(): string;
  setShortDescription(value: string): OrganizationTicket;

  getPriority(): string;
  setPriority(value: string): OrganizationTicket;

  getState(): string;
  setState(value: string): OrganizationTicket;

  getCreateTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreateTime(value?: google_protobuf_timestamp_pb.Timestamp): OrganizationTicket;
  hasCreateTime(): boolean;
  clearCreateTime(): OrganizationTicket;

  getDepartment(): string;
  setDepartment(value: string): OrganizationTicket;

  getRequester(): string;
  setRequester(value: string): OrganizationTicket;

  getTaskType(): string;
  setTaskType(value: string): OrganizationTicket;

  getUrl(): string;
  setUrl(value: string): OrganizationTicket;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationTicket.AsObject;
  static toObject(includeInstance: boolean, msg: OrganizationTicket): OrganizationTicket.AsObject;
  static serializeBinaryToWriter(message: OrganizationTicket, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationTicket;
  static deserializeBinaryFromReader(message: OrganizationTicket, reader: jspb.BinaryReader): OrganizationTicket;
}

export namespace OrganizationTicket {
  export type AsObject = {
    number: string,
    shortDescription: string,
    priority: string,
    state: string,
    createTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    department: string,
    requester: string,
    taskType: string,
    url: string,
  }
}

