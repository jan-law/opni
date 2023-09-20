// @generated by protoc-gen-es v1.3.0 with parameter "target=ts,import_extension=none,ts_nocheck=false"
// @generated from file github.com/rancher/opni/pkg/apis/management/v1/management.proto (package management, syntax proto3)
/* eslint-disable */

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Duration, Message, MethodDescriptorProto, proto3, protoInt64, ServiceDescriptorProto } from "@bufbuild/protobuf";
import { CertInfo, Cluster, LabelSelector, MatchOptions, Reference, ReferenceList, TokenCapability } from "../../core/v1/core_pb";
import { HttpRule } from "../../../../../../../google/api/http_pb";
import { Details, InstallRequest, UninstallRequest } from "../../capability/v1/capability_pb";

/**
 * @generated from enum management.WatchEventType
 */
export enum WatchEventType {
  /**
   * @generated from enum value: Put = 0;
   */
  Put = 0,

  /**
   * @generated from enum value: Delete = 2;
   */
  Delete = 2,
}
// Retrieve enum metadata with: proto3.getEnumType(WatchEventType)
proto3.util.setEnumType(WatchEventType, "management.WatchEventType", [
  { no: 0, name: "Put" },
  { no: 2, name: "Delete" },
]);

/**
 * @generated from message management.CreateBootstrapTokenRequest
 */
export class CreateBootstrapTokenRequest extends Message<CreateBootstrapTokenRequest> {
  /**
   * @generated from field: google.protobuf.Duration ttl = 1;
   */
  ttl?: Duration;

  /**
   * @generated from field: map<string, string> labels = 2;
   */
  labels: { [key: string]: string } = {};

  /**
   * @generated from field: repeated core.TokenCapability capabilities = 3;
   */
  capabilities: TokenCapability[] = [];

  /**
   * @generated from field: int64 maxUsages = 4;
   */
  maxUsages = protoInt64.zero;

  constructor(data?: PartialMessage<CreateBootstrapTokenRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.CreateBootstrapTokenRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "ttl", kind: "message", T: Duration },
    { no: 2, name: "labels", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
    { no: 3, name: "capabilities", kind: "message", T: TokenCapability, repeated: true },
    { no: 4, name: "maxUsages", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateBootstrapTokenRequest {
    return new CreateBootstrapTokenRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateBootstrapTokenRequest {
    return new CreateBootstrapTokenRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateBootstrapTokenRequest {
    return new CreateBootstrapTokenRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CreateBootstrapTokenRequest | PlainMessage<CreateBootstrapTokenRequest> | undefined, b: CreateBootstrapTokenRequest | PlainMessage<CreateBootstrapTokenRequest> | undefined): boolean {
    return proto3.util.equals(CreateBootstrapTokenRequest, a, b);
  }
}

/**
 * @generated from message management.CertsInfoResponse
 */
export class CertsInfoResponse extends Message<CertsInfoResponse> {
  /**
   * @generated from field: repeated core.CertInfo chain = 1;
   */
  chain: CertInfo[] = [];

  constructor(data?: PartialMessage<CertsInfoResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.CertsInfoResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "chain", kind: "message", T: CertInfo, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CertsInfoResponse {
    return new CertsInfoResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CertsInfoResponse {
    return new CertsInfoResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CertsInfoResponse {
    return new CertsInfoResponse().fromJsonString(jsonString, options);
  }

  static equals(a: CertsInfoResponse | PlainMessage<CertsInfoResponse> | undefined, b: CertsInfoResponse | PlainMessage<CertsInfoResponse> | undefined): boolean {
    return proto3.util.equals(CertsInfoResponse, a, b);
  }
}

/**
 * @generated from message management.ListClustersRequest
 */
export class ListClustersRequest extends Message<ListClustersRequest> {
  /**
   * @generated from field: core.LabelSelector matchLabels = 1;
   */
  matchLabels?: LabelSelector;

  /**
   * @generated from field: core.MatchOptions matchOptions = 2;
   */
  matchOptions = MatchOptions.Default;

  constructor(data?: PartialMessage<ListClustersRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.ListClustersRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "matchLabels", kind: "message", T: LabelSelector },
    { no: 2, name: "matchOptions", kind: "enum", T: proto3.getEnumType(MatchOptions) },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ListClustersRequest {
    return new ListClustersRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ListClustersRequest {
    return new ListClustersRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ListClustersRequest {
    return new ListClustersRequest().fromJsonString(jsonString, options);
  }

  static equals(a: ListClustersRequest | PlainMessage<ListClustersRequest> | undefined, b: ListClustersRequest | PlainMessage<ListClustersRequest> | undefined): boolean {
    return proto3.util.equals(ListClustersRequest, a, b);
  }
}

/**
 * @generated from message management.EditClusterRequest
 */
export class EditClusterRequest extends Message<EditClusterRequest> {
  /**
   * @generated from field: core.Reference cluster = 1;
   */
  cluster?: Reference;

  /**
   * @generated from field: map<string, string> labels = 2;
   */
  labels: { [key: string]: string } = {};

  constructor(data?: PartialMessage<EditClusterRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.EditClusterRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "cluster", kind: "message", T: Reference },
    { no: 2, name: "labels", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): EditClusterRequest {
    return new EditClusterRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): EditClusterRequest {
    return new EditClusterRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): EditClusterRequest {
    return new EditClusterRequest().fromJsonString(jsonString, options);
  }

  static equals(a: EditClusterRequest | PlainMessage<EditClusterRequest> | undefined, b: EditClusterRequest | PlainMessage<EditClusterRequest> | undefined): boolean {
    return proto3.util.equals(EditClusterRequest, a, b);
  }
}

/**
 * @generated from message management.WatchClustersRequest
 */
export class WatchClustersRequest extends Message<WatchClustersRequest> {
  /**
   * @generated from field: core.ReferenceList knownClusters = 1;
   */
  knownClusters?: ReferenceList;

  constructor(data?: PartialMessage<WatchClustersRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.WatchClustersRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "knownClusters", kind: "message", T: ReferenceList },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): WatchClustersRequest {
    return new WatchClustersRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): WatchClustersRequest {
    return new WatchClustersRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): WatchClustersRequest {
    return new WatchClustersRequest().fromJsonString(jsonString, options);
  }

  static equals(a: WatchClustersRequest | PlainMessage<WatchClustersRequest> | undefined, b: WatchClustersRequest | PlainMessage<WatchClustersRequest> | undefined): boolean {
    return proto3.util.equals(WatchClustersRequest, a, b);
  }
}

/**
 * @generated from message management.WatchEvent
 */
export class WatchEvent extends Message<WatchEvent> {
  /**
   * @generated from field: core.Cluster cluster = 1;
   */
  cluster?: Cluster;

  /**
   * @generated from field: management.WatchEventType type = 2;
   */
  type = WatchEventType.Put;

  /**
   * @generated from field: core.Cluster previous = 3;
   */
  previous?: Cluster;

  constructor(data?: PartialMessage<WatchEvent>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.WatchEvent";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "cluster", kind: "message", T: Cluster },
    { no: 2, name: "type", kind: "enum", T: proto3.getEnumType(WatchEventType) },
    { no: 3, name: "previous", kind: "message", T: Cluster },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): WatchEvent {
    return new WatchEvent().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): WatchEvent {
    return new WatchEvent().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): WatchEvent {
    return new WatchEvent().fromJsonString(jsonString, options);
  }

  static equals(a: WatchEvent | PlainMessage<WatchEvent> | undefined, b: WatchEvent | PlainMessage<WatchEvent> | undefined): boolean {
    return proto3.util.equals(WatchEvent, a, b);
  }
}

/**
 * @generated from message management.APIExtensionInfoList
 */
export class APIExtensionInfoList extends Message<APIExtensionInfoList> {
  /**
   * @generated from field: repeated management.APIExtensionInfo items = 1;
   */
  items: APIExtensionInfo[] = [];

  constructor(data?: PartialMessage<APIExtensionInfoList>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.APIExtensionInfoList";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "items", kind: "message", T: APIExtensionInfo, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): APIExtensionInfoList {
    return new APIExtensionInfoList().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): APIExtensionInfoList {
    return new APIExtensionInfoList().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): APIExtensionInfoList {
    return new APIExtensionInfoList().fromJsonString(jsonString, options);
  }

  static equals(a: APIExtensionInfoList | PlainMessage<APIExtensionInfoList> | undefined, b: APIExtensionInfoList | PlainMessage<APIExtensionInfoList> | undefined): boolean {
    return proto3.util.equals(APIExtensionInfoList, a, b);
  }
}

/**
 * @generated from message management.APIExtensionInfo
 */
export class APIExtensionInfo extends Message<APIExtensionInfo> {
  /**
   * @generated from field: google.protobuf.ServiceDescriptorProto serviceDesc = 1;
   */
  serviceDesc?: ServiceDescriptorProto;

  /**
   * @generated from field: repeated management.HTTPRuleDescriptor rules = 2;
   */
  rules: HTTPRuleDescriptor[] = [];

  constructor(data?: PartialMessage<APIExtensionInfo>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.APIExtensionInfo";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "serviceDesc", kind: "message", T: ServiceDescriptorProto },
    { no: 2, name: "rules", kind: "message", T: HTTPRuleDescriptor, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): APIExtensionInfo {
    return new APIExtensionInfo().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): APIExtensionInfo {
    return new APIExtensionInfo().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): APIExtensionInfo {
    return new APIExtensionInfo().fromJsonString(jsonString, options);
  }

  static equals(a: APIExtensionInfo | PlainMessage<APIExtensionInfo> | undefined, b: APIExtensionInfo | PlainMessage<APIExtensionInfo> | undefined): boolean {
    return proto3.util.equals(APIExtensionInfo, a, b);
  }
}

/**
 * @generated from message management.HTTPRuleDescriptor
 */
export class HTTPRuleDescriptor extends Message<HTTPRuleDescriptor> {
  /**
   * @generated from field: google.api.HttpRule http = 1;
   */
  http?: HttpRule;

  /**
   * @generated from field: google.protobuf.MethodDescriptorProto method = 2;
   */
  method?: MethodDescriptorProto;

  constructor(data?: PartialMessage<HTTPRuleDescriptor>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.HTTPRuleDescriptor";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "http", kind: "message", T: HttpRule },
    { no: 2, name: "method", kind: "message", T: MethodDescriptorProto },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): HTTPRuleDescriptor {
    return new HTTPRuleDescriptor().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): HTTPRuleDescriptor {
    return new HTTPRuleDescriptor().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): HTTPRuleDescriptor {
    return new HTTPRuleDescriptor().fromJsonString(jsonString, options);
  }

  static equals(a: HTTPRuleDescriptor | PlainMessage<HTTPRuleDescriptor> | undefined, b: HTTPRuleDescriptor | PlainMessage<HTTPRuleDescriptor> | undefined): boolean {
    return proto3.util.equals(HTTPRuleDescriptor, a, b);
  }
}

/**
 * @generated from message management.GatewayConfig
 */
export class GatewayConfig extends Message<GatewayConfig> {
  /**
   * @generated from field: repeated management.ConfigDocumentWithSchema documents = 1;
   */
  documents: ConfigDocumentWithSchema[] = [];

  constructor(data?: PartialMessage<GatewayConfig>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.GatewayConfig";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "documents", kind: "message", T: ConfigDocumentWithSchema, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GatewayConfig {
    return new GatewayConfig().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GatewayConfig {
    return new GatewayConfig().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GatewayConfig {
    return new GatewayConfig().fromJsonString(jsonString, options);
  }

  static equals(a: GatewayConfig | PlainMessage<GatewayConfig> | undefined, b: GatewayConfig | PlainMessage<GatewayConfig> | undefined): boolean {
    return proto3.util.equals(GatewayConfig, a, b);
  }
}

/**
 * @generated from message management.ConfigDocumentWithSchema
 */
export class ConfigDocumentWithSchema extends Message<ConfigDocumentWithSchema> {
  /**
   * @generated from field: bytes json = 1;
   */
  json = new Uint8Array(0);

  /**
   * @generated from field: bytes yaml = 2;
   */
  yaml = new Uint8Array(0);

  /**
   * @generated from field: bytes schema = 3;
   */
  schema = new Uint8Array(0);

  constructor(data?: PartialMessage<ConfigDocumentWithSchema>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.ConfigDocumentWithSchema";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "json", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 2, name: "yaml", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 3, name: "schema", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ConfigDocumentWithSchema {
    return new ConfigDocumentWithSchema().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ConfigDocumentWithSchema {
    return new ConfigDocumentWithSchema().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ConfigDocumentWithSchema {
    return new ConfigDocumentWithSchema().fromJsonString(jsonString, options);
  }

  static equals(a: ConfigDocumentWithSchema | PlainMessage<ConfigDocumentWithSchema> | undefined, b: ConfigDocumentWithSchema | PlainMessage<ConfigDocumentWithSchema> | undefined): boolean {
    return proto3.util.equals(ConfigDocumentWithSchema, a, b);
  }
}

/**
 * @generated from message management.ConfigDocument
 */
export class ConfigDocument extends Message<ConfigDocument> {
  /**
   * @generated from field: bytes json = 1;
   */
  json = new Uint8Array(0);

  constructor(data?: PartialMessage<ConfigDocument>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.ConfigDocument";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "json", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ConfigDocument {
    return new ConfigDocument().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ConfigDocument {
    return new ConfigDocument().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ConfigDocument {
    return new ConfigDocument().fromJsonString(jsonString, options);
  }

  static equals(a: ConfigDocument | PlainMessage<ConfigDocument> | undefined, b: ConfigDocument | PlainMessage<ConfigDocument> | undefined): boolean {
    return proto3.util.equals(ConfigDocument, a, b);
  }
}

/**
 * @generated from message management.UpdateConfigRequest
 */
export class UpdateConfigRequest extends Message<UpdateConfigRequest> {
  /**
   * @generated from field: repeated management.ConfigDocument documents = 1;
   */
  documents: ConfigDocument[] = [];

  constructor(data?: PartialMessage<UpdateConfigRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.UpdateConfigRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "documents", kind: "message", T: ConfigDocument, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UpdateConfigRequest {
    return new UpdateConfigRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UpdateConfigRequest {
    return new UpdateConfigRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UpdateConfigRequest {
    return new UpdateConfigRequest().fromJsonString(jsonString, options);
  }

  static equals(a: UpdateConfigRequest | PlainMessage<UpdateConfigRequest> | undefined, b: UpdateConfigRequest | PlainMessage<UpdateConfigRequest> | undefined): boolean {
    return proto3.util.equals(UpdateConfigRequest, a, b);
  }
}

/**
 * @generated from message management.CapabilityList
 */
export class CapabilityList extends Message<CapabilityList> {
  /**
   * @generated from field: repeated management.CapabilityInfo items = 1;
   */
  items: CapabilityInfo[] = [];

  constructor(data?: PartialMessage<CapabilityList>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.CapabilityList";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "items", kind: "message", T: CapabilityInfo, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CapabilityList {
    return new CapabilityList().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CapabilityList {
    return new CapabilityList().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CapabilityList {
    return new CapabilityList().fromJsonString(jsonString, options);
  }

  static equals(a: CapabilityList | PlainMessage<CapabilityList> | undefined, b: CapabilityList | PlainMessage<CapabilityList> | undefined): boolean {
    return proto3.util.equals(CapabilityList, a, b);
  }
}

/**
 * @generated from message management.CapabilityInfo
 */
export class CapabilityInfo extends Message<CapabilityInfo> {
  /**
   * @generated from field: capability.Details details = 1;
   */
  details?: Details;

  /**
   * @generated from field: int32 nodeCount = 2;
   */
  nodeCount = 0;

  constructor(data?: PartialMessage<CapabilityInfo>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.CapabilityInfo";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "details", kind: "message", T: Details },
    { no: 2, name: "nodeCount", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CapabilityInfo {
    return new CapabilityInfo().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CapabilityInfo {
    return new CapabilityInfo().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CapabilityInfo {
    return new CapabilityInfo().fromJsonString(jsonString, options);
  }

  static equals(a: CapabilityInfo | PlainMessage<CapabilityInfo> | undefined, b: CapabilityInfo | PlainMessage<CapabilityInfo> | undefined): boolean {
    return proto3.util.equals(CapabilityInfo, a, b);
  }
}

/**
 * @generated from message management.CapabilityInstallerRequest
 */
export class CapabilityInstallerRequest extends Message<CapabilityInstallerRequest> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: string token = 2;
   */
  token = "";

  /**
   * @generated from field: string pin = 3;
   */
  pin = "";

  constructor(data?: PartialMessage<CapabilityInstallerRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.CapabilityInstallerRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "pin", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CapabilityInstallerRequest {
    return new CapabilityInstallerRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CapabilityInstallerRequest {
    return new CapabilityInstallerRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CapabilityInstallerRequest {
    return new CapabilityInstallerRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CapabilityInstallerRequest | PlainMessage<CapabilityInstallerRequest> | undefined, b: CapabilityInstallerRequest | PlainMessage<CapabilityInstallerRequest> | undefined): boolean {
    return proto3.util.equals(CapabilityInstallerRequest, a, b);
  }
}

/**
 * @generated from message management.CapabilityInstallRequest
 */
export class CapabilityInstallRequest extends Message<CapabilityInstallRequest> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: capability.InstallRequest target = 2;
   */
  target?: InstallRequest;

  constructor(data?: PartialMessage<CapabilityInstallRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.CapabilityInstallRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "target", kind: "message", T: InstallRequest },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CapabilityInstallRequest {
    return new CapabilityInstallRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CapabilityInstallRequest {
    return new CapabilityInstallRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CapabilityInstallRequest {
    return new CapabilityInstallRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CapabilityInstallRequest | PlainMessage<CapabilityInstallRequest> | undefined, b: CapabilityInstallRequest | PlainMessage<CapabilityInstallRequest> | undefined): boolean {
    return proto3.util.equals(CapabilityInstallRequest, a, b);
  }
}

/**
 * @generated from message management.CapabilityInstallerResponse
 */
export class CapabilityInstallerResponse extends Message<CapabilityInstallerResponse> {
  /**
   * @generated from field: string command = 1;
   */
  command = "";

  constructor(data?: PartialMessage<CapabilityInstallerResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.CapabilityInstallerResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "command", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CapabilityInstallerResponse {
    return new CapabilityInstallerResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CapabilityInstallerResponse {
    return new CapabilityInstallerResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CapabilityInstallerResponse {
    return new CapabilityInstallerResponse().fromJsonString(jsonString, options);
  }

  static equals(a: CapabilityInstallerResponse | PlainMessage<CapabilityInstallerResponse> | undefined, b: CapabilityInstallerResponse | PlainMessage<CapabilityInstallerResponse> | undefined): boolean {
    return proto3.util.equals(CapabilityInstallerResponse, a, b);
  }
}

/**
 * @generated from message management.CapabilityUninstallRequest
 */
export class CapabilityUninstallRequest extends Message<CapabilityUninstallRequest> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: capability.UninstallRequest target = 2;
   */
  target?: UninstallRequest;

  constructor(data?: PartialMessage<CapabilityUninstallRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.CapabilityUninstallRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "target", kind: "message", T: UninstallRequest },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CapabilityUninstallRequest {
    return new CapabilityUninstallRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CapabilityUninstallRequest {
    return new CapabilityUninstallRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CapabilityUninstallRequest {
    return new CapabilityUninstallRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CapabilityUninstallRequest | PlainMessage<CapabilityUninstallRequest> | undefined, b: CapabilityUninstallRequest | PlainMessage<CapabilityUninstallRequest> | undefined): boolean {
    return proto3.util.equals(CapabilityUninstallRequest, a, b);
  }
}

/**
 * @generated from message management.CapabilityStatusRequest
 */
export class CapabilityStatusRequest extends Message<CapabilityStatusRequest> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: core.Reference cluster = 2;
   */
  cluster?: Reference;

  constructor(data?: PartialMessage<CapabilityStatusRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.CapabilityStatusRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "cluster", kind: "message", T: Reference },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CapabilityStatusRequest {
    return new CapabilityStatusRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CapabilityStatusRequest {
    return new CapabilityStatusRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CapabilityStatusRequest {
    return new CapabilityStatusRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CapabilityStatusRequest | PlainMessage<CapabilityStatusRequest> | undefined, b: CapabilityStatusRequest | PlainMessage<CapabilityStatusRequest> | undefined): boolean {
    return proto3.util.equals(CapabilityStatusRequest, a, b);
  }
}

/**
 * @generated from message management.CapabilityUninstallCancelRequest
 */
export class CapabilityUninstallCancelRequest extends Message<CapabilityUninstallCancelRequest> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: core.Reference cluster = 2;
   */
  cluster?: Reference;

  constructor(data?: PartialMessage<CapabilityUninstallCancelRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.CapabilityUninstallCancelRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "cluster", kind: "message", T: Reference },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CapabilityUninstallCancelRequest {
    return new CapabilityUninstallCancelRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CapabilityUninstallCancelRequest {
    return new CapabilityUninstallCancelRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CapabilityUninstallCancelRequest {
    return new CapabilityUninstallCancelRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CapabilityUninstallCancelRequest | PlainMessage<CapabilityUninstallCancelRequest> | undefined, b: CapabilityUninstallCancelRequest | PlainMessage<CapabilityUninstallCancelRequest> | undefined): boolean {
    return proto3.util.equals(CapabilityUninstallCancelRequest, a, b);
  }
}

/**
 * @generated from message management.DashboardSettings
 */
export class DashboardSettings extends Message<DashboardSettings> {
  /**
   * @generated from field: optional management.DashboardGlobalSettings global = 1;
   */
  global?: DashboardGlobalSettings;

  /**
   * @generated from field: map<string, string> user = 2;
   */
  user: { [key: string]: string } = {};

  constructor(data?: PartialMessage<DashboardSettings>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.DashboardSettings";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "global", kind: "message", T: DashboardGlobalSettings, opt: true },
    { no: 2, name: "user", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DashboardSettings {
    return new DashboardSettings().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DashboardSettings {
    return new DashboardSettings().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DashboardSettings {
    return new DashboardSettings().fromJsonString(jsonString, options);
  }

  static equals(a: DashboardSettings | PlainMessage<DashboardSettings> | undefined, b: DashboardSettings | PlainMessage<DashboardSettings> | undefined): boolean {
    return proto3.util.equals(DashboardSettings, a, b);
  }
}

/**
 * @generated from message management.DashboardGlobalSettings
 */
export class DashboardGlobalSettings extends Message<DashboardGlobalSettings> {
  /**
   * @generated from field: string defaultImageRepository = 1;
   */
  defaultImageRepository = "";

  /**
   * @generated from field: google.protobuf.Duration defaultTokenTtl = 2;
   */
  defaultTokenTtl?: Duration;

  /**
   * @generated from field: map<string, string> defaultTokenLabels = 3;
   */
  defaultTokenLabels: { [key: string]: string } = {};

  constructor(data?: PartialMessage<DashboardGlobalSettings>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "management.DashboardGlobalSettings";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "defaultImageRepository", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "defaultTokenTtl", kind: "message", T: Duration },
    { no: 3, name: "defaultTokenLabels", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DashboardGlobalSettings {
    return new DashboardGlobalSettings().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DashboardGlobalSettings {
    return new DashboardGlobalSettings().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DashboardGlobalSettings {
    return new DashboardGlobalSettings().fromJsonString(jsonString, options);
  }

  static equals(a: DashboardGlobalSettings | PlainMessage<DashboardGlobalSettings> | undefined, b: DashboardGlobalSettings | PlainMessage<DashboardGlobalSettings> | undefined): boolean {
    return proto3.util.equals(DashboardGlobalSettings, a, b);
  }
}

