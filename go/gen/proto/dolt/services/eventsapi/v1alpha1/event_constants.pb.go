// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dolt/services/eventsapi/v1alpha1/event_constants.proto

package eventsapi

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Platform int32

const (
	Platform_PLATFORM_UNSPECIFIED Platform = 0
	Platform_LINUX                Platform = 1
	Platform_WINDOWS              Platform = 2
	Platform_DARWIN               Platform = 3
)

var Platform_name = map[int32]string{
	0: "PLATFORM_UNSPECIFIED",
	1: "LINUX",
	2: "WINDOWS",
	3: "DARWIN",
}

var Platform_value = map[string]int32{
	"PLATFORM_UNSPECIFIED": 0,
	"LINUX":                1,
	"WINDOWS":              2,
	"DARWIN":               3,
}

func (x Platform) String() string {
	return proto.EnumName(Platform_name, int32(x))
}

func (Platform) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d970d881fa70959f, []int{0}
}

type ClientEventType int32

const (
	ClientEventType_TYPE_UNSPECIFIED                 ClientEventType = 0
	ClientEventType_INIT                             ClientEventType = 1
	ClientEventType_STATUS                           ClientEventType = 2
	ClientEventType_ADD                              ClientEventType = 3
	ClientEventType_RESET                            ClientEventType = 4
	ClientEventType_COMMIT                           ClientEventType = 5
	ClientEventType_SQL                              ClientEventType = 6
	ClientEventType_SQL_SERVER                       ClientEventType = 7
	ClientEventType_LOG                              ClientEventType = 8
	ClientEventType_DIFF                             ClientEventType = 9
	ClientEventType_MERGE                            ClientEventType = 10
	ClientEventType_BRANCH                           ClientEventType = 11
	ClientEventType_CHECKOUT                         ClientEventType = 12
	ClientEventType_REMOTE                           ClientEventType = 13
	ClientEventType_PUSH                             ClientEventType = 14
	ClientEventType_PULL                             ClientEventType = 15
	ClientEventType_FETCH                            ClientEventType = 16
	ClientEventType_CLONE                            ClientEventType = 17
	ClientEventType_LOGIN                            ClientEventType = 18
	ClientEventType_VERSION                          ClientEventType = 19
	ClientEventType_CONFIG                           ClientEventType = 20
	ClientEventType_LS                               ClientEventType = 21
	ClientEventType_SCHEMA                           ClientEventType = 22
	ClientEventType_TABLE_IMPORT                     ClientEventType = 23
	ClientEventType_TABLE_EXPORT                     ClientEventType = 24
	ClientEventType_TABLE_CREATE                     ClientEventType = 25
	ClientEventType_TABLE_RM                         ClientEventType = 26
	ClientEventType_TABLE_MV                         ClientEventType = 27
	ClientEventType_TABLE_CP                         ClientEventType = 28
	ClientEventType_TABLE_SELECT                     ClientEventType = 29
	ClientEventType_TABLE_PUT_ROW                    ClientEventType = 30
	ClientEventType_TABLE_RM_ROW                     ClientEventType = 31
	ClientEventType_CREDS_NEW                        ClientEventType = 32
	ClientEventType_CREDS_RM                         ClientEventType = 33
	ClientEventType_CREDS_LS                         ClientEventType = 34
	ClientEventType_CONF_CAT                         ClientEventType = 35
	ClientEventType_CONF_RESOLVE                     ClientEventType = 36
	ClientEventType_REMOTEAPI_GET_REPO_METADATA      ClientEventType = 37
	ClientEventType_REMOTEAPI_HAS_CHUNKS             ClientEventType = 38
	ClientEventType_REMOTEAPI_GET_DOWNLOAD_LOCATIONS ClientEventType = 39
	ClientEventType_REMOTEAPI_GET_UPLOAD_LOCATIONS   ClientEventType = 40
	ClientEventType_REMOTEAPI_REBASE                 ClientEventType = 41
	ClientEventType_REMOTEAPI_ROOT                   ClientEventType = 42
	ClientEventType_REMOTEAPI_COMMIT                 ClientEventType = 43
	ClientEventType_REMOTEAPI_LIST_TABLE_FILES       ClientEventType = 44
	ClientEventType_BLAME                            ClientEventType = 45
)

var ClientEventType_name = map[int32]string{
	0:  "TYPE_UNSPECIFIED",
	1:  "INIT",
	2:  "STATUS",
	3:  "ADD",
	4:  "RESET",
	5:  "COMMIT",
	6:  "SQL",
	7:  "SQL_SERVER",
	8:  "LOG",
	9:  "DIFF",
	10: "MERGE",
	11: "BRANCH",
	12: "CHECKOUT",
	13: "REMOTE",
	14: "PUSH",
	15: "PULL",
	16: "FETCH",
	17: "CLONE",
	18: "LOGIN",
	19: "VERSION",
	20: "CONFIG",
	21: "LS",
	22: "SCHEMA",
	23: "TABLE_IMPORT",
	24: "TABLE_EXPORT",
	25: "TABLE_CREATE",
	26: "TABLE_RM",
	27: "TABLE_MV",
	28: "TABLE_CP",
	29: "TABLE_SELECT",
	30: "TABLE_PUT_ROW",
	31: "TABLE_RM_ROW",
	32: "CREDS_NEW",
	33: "CREDS_RM",
	34: "CREDS_LS",
	35: "CONF_CAT",
	36: "CONF_RESOLVE",
	37: "REMOTEAPI_GET_REPO_METADATA",
	38: "REMOTEAPI_HAS_CHUNKS",
	39: "REMOTEAPI_GET_DOWNLOAD_LOCATIONS",
	40: "REMOTEAPI_GET_UPLOAD_LOCATIONS",
	41: "REMOTEAPI_REBASE",
	42: "REMOTEAPI_ROOT",
	43: "REMOTEAPI_COMMIT",
	44: "REMOTEAPI_LIST_TABLE_FILES",
	45: "BLAME",
}

var ClientEventType_value = map[string]int32{
	"TYPE_UNSPECIFIED":                 0,
	"INIT":                             1,
	"STATUS":                           2,
	"ADD":                              3,
	"RESET":                            4,
	"COMMIT":                           5,
	"SQL":                              6,
	"SQL_SERVER":                       7,
	"LOG":                              8,
	"DIFF":                             9,
	"MERGE":                            10,
	"BRANCH":                           11,
	"CHECKOUT":                         12,
	"REMOTE":                           13,
	"PUSH":                             14,
	"PULL":                             15,
	"FETCH":                            16,
	"CLONE":                            17,
	"LOGIN":                            18,
	"VERSION":                          19,
	"CONFIG":                           20,
	"LS":                               21,
	"SCHEMA":                           22,
	"TABLE_IMPORT":                     23,
	"TABLE_EXPORT":                     24,
	"TABLE_CREATE":                     25,
	"TABLE_RM":                         26,
	"TABLE_MV":                         27,
	"TABLE_CP":                         28,
	"TABLE_SELECT":                     29,
	"TABLE_PUT_ROW":                    30,
	"TABLE_RM_ROW":                     31,
	"CREDS_NEW":                        32,
	"CREDS_RM":                         33,
	"CREDS_LS":                         34,
	"CONF_CAT":                         35,
	"CONF_RESOLVE":                     36,
	"REMOTEAPI_GET_REPO_METADATA":      37,
	"REMOTEAPI_HAS_CHUNKS":             38,
	"REMOTEAPI_GET_DOWNLOAD_LOCATIONS": 39,
	"REMOTEAPI_GET_UPLOAD_LOCATIONS":   40,
	"REMOTEAPI_REBASE":                 41,
	"REMOTEAPI_ROOT":                   42,
	"REMOTEAPI_COMMIT":                 43,
	"REMOTEAPI_LIST_TABLE_FILES":       44,
	"BLAME":                            45,
}

func (x ClientEventType) String() string {
	return proto.EnumName(ClientEventType_name, int32(x))
}

func (ClientEventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d970d881fa70959f, []int{1}
}

type MetricID int32

const (
	MetricID_METRIC_UNSPECIFIED  MetricID = 0
	MetricID_BYTES_DOWNLOADED    MetricID = 1
	MetricID_DOWNLOAD_MS_ELAPSED MetricID = 2
	MetricID_REMOTEAPI_RPC_ERROR MetricID = 3
)

var MetricID_name = map[int32]string{
	0: "METRIC_UNSPECIFIED",
	1: "BYTES_DOWNLOADED",
	2: "DOWNLOAD_MS_ELAPSED",
	3: "REMOTEAPI_RPC_ERROR",
}

var MetricID_value = map[string]int32{
	"METRIC_UNSPECIFIED":  0,
	"BYTES_DOWNLOADED":    1,
	"DOWNLOAD_MS_ELAPSED": 2,
	"REMOTEAPI_RPC_ERROR": 3,
}

func (x MetricID) String() string {
	return proto.EnumName(MetricID_name, int32(x))
}

func (MetricID) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d970d881fa70959f, []int{2}
}

type AttributeID int32

const (
	AttributeID_ATTRIBUTE_UNSPECIFIED AttributeID = 0
	AttributeID_REMOTE_URL_SCHEME     AttributeID = 2
)

var AttributeID_name = map[int32]string{
	0: "ATTRIBUTE_UNSPECIFIED",
	2: "REMOTE_URL_SCHEME",
}

var AttributeID_value = map[string]int32{
	"ATTRIBUTE_UNSPECIFIED": 0,
	"REMOTE_URL_SCHEME":     2,
}

func (x AttributeID) String() string {
	return proto.EnumName(AttributeID_name, int32(x))
}

func (AttributeID) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d970d881fa70959f, []int{3}
}

type AppID int32

const (
	AppID_APP_ID_UNSPECIFIED AppID = 0
	AppID_APP_DOLT           AppID = 1
)

var AppID_name = map[int32]string{
	0: "APP_ID_UNSPECIFIED",
	1: "APP_DOLT",
}

var AppID_value = map[string]int32{
	"APP_ID_UNSPECIFIED": 0,
	"APP_DOLT":           1,
}

func (x AppID) String() string {
	return proto.EnumName(AppID_name, int32(x))
}

func (AppID) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d970d881fa70959f, []int{4}
}

func init() {
	proto.RegisterEnum("dolt.services.eventsapi.v1alpha1.Platform", Platform_name, Platform_value)
	proto.RegisterEnum("dolt.services.eventsapi.v1alpha1.ClientEventType", ClientEventType_name, ClientEventType_value)
	proto.RegisterEnum("dolt.services.eventsapi.v1alpha1.MetricID", MetricID_name, MetricID_value)
	proto.RegisterEnum("dolt.services.eventsapi.v1alpha1.AttributeID", AttributeID_name, AttributeID_value)
	proto.RegisterEnum("dolt.services.eventsapi.v1alpha1.AppID", AppID_name, AppID_value)
}

func init() {
	proto.RegisterFile("dolt/services/eventsapi/v1alpha1/event_constants.proto", fileDescriptor_d970d881fa70959f)
}

var fileDescriptor_d970d881fa70959f = []byte{
	// 750 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0xdb, 0x72, 0xdb, 0x36,
	0x10, 0xad, 0x7c, 0x95, 0xe1, 0x4b, 0xd6, 0x88, 0x9d, 0x6b, 0xeb, 0xb8, 0x6e, 0x7a, 0x53, 0x6b,
	0x69, 0x32, 0x9d, 0xe9, 0x4b, 0x9f, 0x20, 0x70, 0x25, 0x61, 0x02, 0x10, 0x0c, 0x00, 0x5a, 0x49,
	0x5f, 0x30, 0xb2, 0xc2, 0x3a, 0xec, 0x28, 0x92, 0x2a, 0xd1, 0x9e, 0xe9, 0x4f, 0xf7, 0x1b, 0x3a,
	0x20, 0x1b, 0x31, 0xd6, 0x4b, 0xdf, 0xb0, 0x67, 0xcf, 0x39, 0x58, 0x2c, 0x97, 0x4b, 0x7e, 0x7d,
	0x3f, 0x9b, 0x14, 0x9d, 0x65, 0xb6, 0xb8, 0xcb, 0xc7, 0xd9, 0xb2, 0x93, 0xdd, 0x65, 0xd3, 0x62,
	0x39, 0x9a, 0xe7, 0x9d, 0xbb, 0x57, 0xa3, 0xc9, 0xfc, 0xc3, 0xe8, 0x55, 0x05, 0xf9, 0xf1, 0x6c,
	0xba, 0x2c, 0x46, 0xd3, 0x62, 0xd9, 0x9e, 0x2f, 0x66, 0xc5, 0x8c, 0x9e, 0x07, 0x5d, 0xfb, 0x93,
	0xae, 0xbd, 0xd2, 0xb5, 0x3f, 0xe9, 0x5a, 0x03, 0xd2, 0x4c, 0x26, 0xa3, 0xe2, 0x8f, 0xd9, 0xe2,
	0x23, 0x7d, 0x42, 0x4e, 0x12, 0xc9, 0x5c, 0x4f, 0x1b, 0xe5, 0xd3, 0xd8, 0x26, 0xc8, 0x45, 0x4f,
	0x60, 0x04, 0x5f, 0xd0, 0x3d, 0xb2, 0x2d, 0x45, 0x9c, 0xbe, 0x85, 0x06, 0xdd, 0x27, 0xbb, 0x43,
	0x11, 0x47, 0x7a, 0x68, 0x61, 0x83, 0x12, 0xb2, 0x13, 0x31, 0x33, 0x14, 0x31, 0x6c, 0xb6, 0xfe,
	0xd9, 0x26, 0x0f, 0xf8, 0x24, 0xcf, 0xa6, 0x05, 0x86, 0x6b, 0xdc, 0xdf, 0xf3, 0x8c, 0x9e, 0x10,
	0x70, 0xef, 0x12, 0x5c, 0x73, 0x6b, 0x92, 0x2d, 0x11, 0x0b, 0x07, 0x8d, 0xa0, 0xb7, 0x8e, 0xb9,
	0x34, 0x78, 0xed, 0x92, 0x4d, 0x16, 0x45, 0xb0, 0x19, 0x2e, 0x33, 0x68, 0xd1, 0xc1, 0x56, 0xc8,
	0x73, 0xad, 0x94, 0x70, 0xb0, 0x1d, 0xf2, 0xf6, 0x8d, 0x84, 0x1d, 0x7a, 0x44, 0x88, 0x7d, 0x23,
	0xbd, 0x45, 0x73, 0x85, 0x06, 0x76, 0x43, 0x42, 0xea, 0x3e, 0x34, 0x83, 0x6f, 0x24, 0x7a, 0x3d,
	0xd8, 0x0b, 0x16, 0x0a, 0x4d, 0x1f, 0x81, 0x04, 0x8b, 0xae, 0x61, 0x31, 0x1f, 0xc0, 0x3e, 0x3d,
	0x20, 0x4d, 0x3e, 0x40, 0xfe, 0x5a, 0xa7, 0x0e, 0x0e, 0x42, 0xc6, 0xa0, 0xd2, 0x0e, 0xe1, 0x30,
	0x48, 0x93, 0xd4, 0x0e, 0xe0, 0xa8, 0x3a, 0x49, 0x09, 0x0f, 0x82, 0x49, 0x0f, 0x1d, 0x1f, 0x00,
	0x84, 0x23, 0x97, 0x3a, 0x46, 0x38, 0x2e, 0x5b, 0xa1, 0xfb, 0x22, 0x06, 0x1a, 0x5a, 0x71, 0x85,
	0xc6, 0x0a, 0x1d, 0xc3, 0xc3, 0xaa, 0xd4, 0xb8, 0x27, 0xfa, 0x70, 0x42, 0x77, 0xc8, 0x86, 0xb4,
	0x70, 0x5a, 0x3e, 0x8f, 0x0f, 0x50, 0x31, 0x78, 0x44, 0x81, 0x1c, 0x38, 0xd6, 0x95, 0xe8, 0x85,
	0x4a, 0xb4, 0x71, 0xf0, 0xb8, 0x46, 0xf0, 0x6d, 0x89, 0x3c, 0xa9, 0x11, 0x6e, 0x90, 0x39, 0x84,
	0xa7, 0xa1, 0xe2, 0x0a, 0x31, 0x0a, 0x9e, 0xd5, 0x91, 0xba, 0x82, 0xe7, 0x75, 0xc4, 0x13, 0xf8,
	0xb2, 0xd6, 0x5a, 0x94, 0xc8, 0x1d, 0x7c, 0x45, 0x8f, 0xc9, 0x61, 0x85, 0x24, 0xa9, 0xf3, 0x46,
	0x0f, 0xe1, 0xac, 0x26, 0x19, 0x55, 0x22, 0x2f, 0xe8, 0x21, 0xd9, 0xe3, 0x06, 0x23, 0xeb, 0x63,
	0x1c, 0xc2, 0x79, 0xd9, 0xa1, 0x32, 0x34, 0x0a, 0xbe, 0xae, 0x23, 0x69, 0xe1, 0xa2, 0x8c, 0x74,
	0xdc, 0xf3, 0x9c, 0x39, 0xf8, 0x26, 0x58, 0x95, 0x91, 0x41, 0xab, 0xe5, 0x15, 0xc2, 0x4b, 0xfa,
	0x82, 0x3c, 0xaf, 0xfa, 0xc9, 0x12, 0xe1, 0xfb, 0xe8, 0xbc, 0xc1, 0x44, 0x7b, 0x85, 0x8e, 0x45,
	0xcc, 0x31, 0xf8, 0x36, 0xcc, 0x57, 0x4d, 0x18, 0x30, 0xeb, 0xf9, 0x20, 0x8d, 0x5f, 0x5b, 0xf8,
	0x8e, 0xbe, 0x24, 0xe7, 0xf7, 0xa5, 0x91, 0x1e, 0xc6, 0x52, 0xb3, 0xc8, 0x4b, 0xcd, 0x99, 0x13,
	0x3a, 0xb6, 0xf0, 0x3d, 0xbd, 0x20, 0x67, 0xf7, 0x59, 0x69, 0xb2, 0xc6, 0xf9, 0x21, 0x4c, 0x5c,
	0xcd, 0x31, 0xd8, 0x65, 0x16, 0xe1, 0x47, 0x4a, 0xc9, 0xd1, 0x67, 0xa8, 0xd6, 0x0e, 0x5a, 0xf7,
	0x99, 0xff, 0x4d, 0xd9, 0x4f, 0xf4, 0x8c, 0x3c, 0xab, 0x51, 0x29, 0xac, 0xf3, 0x55, 0xc3, 0x7a,
	0x42, 0xa2, 0x85, 0x9f, 0xc3, 0xe7, 0xef, 0x4a, 0xa6, 0x10, 0x2e, 0x5b, 0x7f, 0x92, 0xa6, 0xca,
	0x8a, 0x45, 0x3e, 0x16, 0x11, 0x7d, 0x44, 0xa8, 0x42, 0x67, 0x04, 0x5f, 0x1b, 0xf5, 0x13, 0x02,
	0xdd, 0x77, 0x0e, 0xed, 0xea, 0x41, 0x18, 0x41, 0x83, 0x3e, 0x26, 0x0f, 0x57, 0x0f, 0x54, 0xd6,
	0xa3, 0x64, 0x89, 0xc5, 0x08, 0x36, 0x42, 0xe2, 0xb3, 0x3a, 0x13, 0xee, 0xd1, 0x18, 0x6d, 0x60,
	0xb3, 0x85, 0x64, 0x9f, 0x15, 0xc5, 0x22, 0xbf, 0xbe, 0x2d, 0x32, 0x11, 0xd1, 0xa7, 0xe4, 0x94,
	0x39, 0x67, 0x44, 0x37, 0x75, 0xeb, 0x3f, 0xd7, 0x29, 0x39, 0xae, 0x2c, 0x7c, 0x6a, 0xa4, 0x2f,
	0xc7, 0x0f, 0x61, 0xe3, 0x62, 0xab, 0xd9, 0x80, 0x46, 0xeb, 0x92, 0x6c, 0xb3, 0xf9, 0xbc, 0xaa,
	0x97, 0x25, 0x89, 0x17, 0xd1, 0x9a, 0xfa, 0x80, 0x34, 0x03, 0x1e, 0x69, 0xe9, 0xa0, 0xd1, 0x1d,
	0xfe, 0x9e, 0xde, 0xe4, 0xc5, 0x87, 0xdb, 0xeb, 0xf6, 0x78, 0xf6, 0xb1, 0x33, 0xc9, 0xff, 0xba,
	0xcd, 0xdf, 0x8f, 0x8a, 0xd1, 0x65, 0x3e, 0x1d, 0x77, 0xca, 0x8d, 0x74, 0x33, 0xeb, 0xdc, 0x64,
	0xd3, 0x4e, 0xb9, 0x6c, 0x3a, 0xff, 0xb7, 0xa3, 0x7e, 0x5b, 0x41, 0xd7, 0x3b, 0xa5, 0xe2, 0x97,
	0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x8e, 0xc2, 0xb6, 0x75, 0xd8, 0x04, 0x00, 0x00,
}
