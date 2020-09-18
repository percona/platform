import * as jspb from 'google-protobuf'



export class CheckDoubleSHA1Request extends jspb.Message {
  getHashPrefix(): Uint8Array | string;
  getHashPrefix_asU8(): Uint8Array;
  getHashPrefix_asB64(): string;
  setHashPrefix(value: Uint8Array | string): CheckDoubleSHA1Request;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CheckDoubleSHA1Request.AsObject;
  static toObject(includeInstance: boolean, msg: CheckDoubleSHA1Request): CheckDoubleSHA1Request.AsObject;
  static serializeBinaryToWriter(message: CheckDoubleSHA1Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CheckDoubleSHA1Request;
  static deserializeBinaryFromReader(message: CheckDoubleSHA1Request, reader: jspb.BinaryReader): CheckDoubleSHA1Request;
}

export namespace CheckDoubleSHA1Request {
  export type AsObject = {
    hashPrefix: Uint8Array | string,
  }
}

export class CheckDoubleSHA1Response extends jspb.Message {
  getResultsList(): Array<CheckDoubleSHA1Response.Result>;
  setResultsList(value: Array<CheckDoubleSHA1Response.Result>): CheckDoubleSHA1Response;
  clearResultsList(): CheckDoubleSHA1Response;
  addResults(value?: CheckDoubleSHA1Response.Result, index?: number): CheckDoubleSHA1Response.Result;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CheckDoubleSHA1Response.AsObject;
  static toObject(includeInstance: boolean, msg: CheckDoubleSHA1Response): CheckDoubleSHA1Response.AsObject;
  static serializeBinaryToWriter(message: CheckDoubleSHA1Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CheckDoubleSHA1Response;
  static deserializeBinaryFromReader(message: CheckDoubleSHA1Response, reader: jspb.BinaryReader): CheckDoubleSHA1Response;
}

export namespace CheckDoubleSHA1Response {
  export type AsObject = {
    resultsList: Array<CheckDoubleSHA1Response.Result.AsObject>,
  }

  export class Result extends jspb.Message {
    getHash(): Uint8Array | string;
    getHash_asU8(): Uint8Array;
    getHash_asB64(): string;
    setHash(value: Uint8Array | string): Result;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Result.AsObject;
    static toObject(includeInstance: boolean, msg: Result): Result.AsObject;
    static serializeBinaryToWriter(message: Result, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Result;
    static deserializeBinaryFromReader(message: Result, reader: jspb.BinaryReader): Result;
  }

  export namespace Result {
    export type AsObject = {
      hash: Uint8Array | string,
    }
  }

}

