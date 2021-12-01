import * as jspb from 'google-protobuf'

import * as github_com_mwitkow_go$proto$validators_validator_pb from '../github.com/mwitkow/go-proto-validators/validator_pb';
import * as google_api_annotations_pb from '../google/api/annotations_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';


export class PMMServerSSODetails extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): PMMServerSSODetails;

  getClientSecret(): string;
  setClientSecret(value: string): PMMServerSSODetails;

  getIssuerUrl(): string;
  setIssuerUrl(value: string): PMMServerSSODetails;

  getScope(): string;
  setScope(value: string): PMMServerSSODetails;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PMMServerSSODetails.AsObject;
  static toObject(includeInstance: boolean, msg: PMMServerSSODetails): PMMServerSSODetails.AsObject;
  static serializeBinaryToWriter(message: PMMServerSSODetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PMMServerSSODetails;
  static deserializeBinaryFromReader(message: PMMServerSSODetails, reader: jspb.BinaryReader): PMMServerSSODetails;
}

export namespace PMMServerSSODetails {
  export type AsObject = {
    clientId: string,
    clientSecret: string,
    issuerUrl: string,
    scope: string,
  }
}

export class ConnectPMMRequest extends jspb.Message {
  getPmmServerId(): string;
  setPmmServerId(value: string): ConnectPMMRequest;

  getPmmServerName(): string;
  setPmmServerName(value: string): ConnectPMMRequest;

  getPmmServerUrl(): string;
  setPmmServerUrl(value: string): ConnectPMMRequest;

  getPmmServerOauthCallbackUrl(): string;
  setPmmServerOauthCallbackUrl(value: string): ConnectPMMRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConnectPMMRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConnectPMMRequest): ConnectPMMRequest.AsObject;
  static serializeBinaryToWriter(message: ConnectPMMRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConnectPMMRequest;
  static deserializeBinaryFromReader(message: ConnectPMMRequest, reader: jspb.BinaryReader): ConnectPMMRequest;
}

export namespace ConnectPMMRequest {
  export type AsObject = {
    pmmServerId: string,
    pmmServerName: string,
    pmmServerUrl: string,
    pmmServerOauthCallbackUrl: string,
  }
}

export class ConnectPMMResponse extends jspb.Message {
  getSsoDetails(): PMMServerSSODetails | undefined;
  setSsoDetails(value?: PMMServerSSODetails): ConnectPMMResponse;
  hasSsoDetails(): boolean;
  clearSsoDetails(): ConnectPMMResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConnectPMMResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ConnectPMMResponse): ConnectPMMResponse.AsObject;
  static serializeBinaryToWriter(message: ConnectPMMResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConnectPMMResponse;
  static deserializeBinaryFromReader(message: ConnectPMMResponse, reader: jspb.BinaryReader): ConnectPMMResponse;
}

export namespace ConnectPMMResponse {
  export type AsObject = {
    ssoDetails?: PMMServerSSODetails.AsObject,
  }
}

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
  getOrgId(): string;
  setOrgId(value: string): DeleteOrganizationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteOrganizationRequest): DeleteOrganizationRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteOrganizationRequest;
  static deserializeBinaryFromReader(message: DeleteOrganizationRequest, reader: jspb.BinaryReader): DeleteOrganizationRequest;
}

export namespace DeleteOrganizationRequest {
  export type AsObject = {
    orgId: string,
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

export class SearchUserCompanyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchUserCompanyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SearchUserCompanyRequest): SearchUserCompanyRequest.AsObject;
  static serializeBinaryToWriter(message: SearchUserCompanyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchUserCompanyRequest;
  static deserializeBinaryFromReader(message: SearchUserCompanyRequest, reader: jspb.BinaryReader): SearchUserCompanyRequest;
}

export namespace SearchUserCompanyRequest {
  export type AsObject = {
  }
}

export class SearchUserCompanyResponse extends jspb.Message {
  getName(): string;
  setName(value: string): SearchUserCompanyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchUserCompanyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SearchUserCompanyResponse): SearchUserCompanyResponse.AsObject;
  static serializeBinaryToWriter(message: SearchUserCompanyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchUserCompanyResponse;
  static deserializeBinaryFromReader(message: SearchUserCompanyResponse, reader: jspb.BinaryReader): SearchUserCompanyResponse;
}

export namespace SearchUserCompanyResponse {
  export type AsObject = {
    name: string,
  }
}

export class InviteMemberRequest extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): InviteMemberRequest;

  getOrgId(): string;
  setOrgId(value: string): InviteMemberRequest;

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
    orgId: string,
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
  getMemberId(): string;
  setMemberId(value: string): OrganizationMember;

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
    memberId: string,
    username: string,
    firstName: string,
    lastName: string,
    role: string,
    status: string,
  }
}

export class SearchMembersRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): SearchMembersRequest;

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
    orgId: string,
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

export class UpdateMemberRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): UpdateMemberRequest;

  getMemberId(): string;
  setMemberId(value: string): UpdateMemberRequest;

  getRole(): string;
  setRole(value: string): UpdateMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateMemberRequest): UpdateMemberRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateMemberRequest;
  static deserializeBinaryFromReader(message: UpdateMemberRequest, reader: jspb.BinaryReader): UpdateMemberRequest;
}

export namespace UpdateMemberRequest {
  export type AsObject = {
    orgId: string,
    memberId: string,
    role: string,
  }
}

export class UpdateMemberResponse extends jspb.Message {
  getMember(): OrganizationMember | undefined;
  setMember(value?: OrganizationMember): UpdateMemberResponse;
  hasMember(): boolean;
  clearMember(): UpdateMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateMemberResponse): UpdateMemberResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateMemberResponse;
  static deserializeBinaryFromReader(message: UpdateMemberResponse, reader: jspb.BinaryReader): UpdateMemberResponse;
}

export namespace UpdateMemberResponse {
  export type AsObject = {
    member?: OrganizationMember.AsObject,
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

