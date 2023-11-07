// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: concord/github/v1/github.proto

package gh_pb

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Organization struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Defaults     *Defaults     `protobuf:"bytes,2,opt,name=defaults,proto3" json:"defaults,omitempty"`
	Teams        []string      `protobuf:"bytes,10,rep,name=teams,proto3" json:"teams,omitempty"`
	People       []*People     `protobuf:"bytes,11,rep,name=people,proto3" json:"people,omitempty"`
	Repositories []*Repository `protobuf:"bytes,12,rep,name=repositories,proto3" json:"repositories,omitempty"`
	Labels       []string      `protobuf:"bytes,13,rep,name=labels,proto3" json:"labels,omitempty"`
}

func (x *Organization) Reset() {
	*x = Organization{}
	if protoimpl.UnsafeEnabled {
		mi := &file_concord_github_v1_github_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Organization) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Organization) ProtoMessage() {}

func (x *Organization) ProtoReflect() protoreflect.Message {
	mi := &file_concord_github_v1_github_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Organization.ProtoReflect.Descriptor instead.
func (*Organization) Descriptor() ([]byte, []int) {
	return file_concord_github_v1_github_proto_rawDescGZIP(), []int{0}
}

func (x *Organization) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Organization) GetDefaults() *Defaults {
	if x != nil {
		return x.Defaults
	}
	return nil
}

func (x *Organization) GetTeams() []string {
	if x != nil {
		return x.Teams
	}
	return nil
}

func (x *Organization) GetPeople() []*People {
	if x != nil {
		return x.People
	}
	return nil
}

func (x *Organization) GetRepositories() []*Repository {
	if x != nil {
		return x.Repositories
	}
	return nil
}

func (x *Organization) GetLabels() []string {
	if x != nil {
		return x.Labels
	}
	return nil
}

// Defaults are overriden by the same settings specified in the repository
type Defaults struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Private                *bool     `protobuf:"varint,1,opt,name=private,proto3,oneof" json:"private,omitempty"`
	DefaultBranch          *string   `protobuf:"bytes,2,opt,name=default_branch,json=defaultBranch,proto3,oneof" json:"default_branch,omitempty"`
	AllowAutoMerge         *bool     `protobuf:"varint,3,opt,name=allow_auto_merge,json=allowAutoMerge,proto3,oneof" json:"allow_auto_merge,omitempty"`
	AutoDeleteHeadBranches *bool     `protobuf:"varint,4,opt,name=auto_delete_head_branches,json=autoDeleteHeadBranches,proto3,oneof" json:"auto_delete_head_branches,omitempty"`
	ProtectedBranches      []*Branch `protobuf:"bytes,5,rep,name=protected_branches,json=protectedBranches,proto3" json:"protected_branches,omitempty"`
	Files                  []*File   `protobuf:"bytes,6,rep,name=files,proto3" json:"files,omitempty"` //repeated Secret     secrets               = 7;
}

func (x *Defaults) Reset() {
	*x = Defaults{}
	if protoimpl.UnsafeEnabled {
		mi := &file_concord_github_v1_github_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Defaults) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Defaults) ProtoMessage() {}

func (x *Defaults) ProtoReflect() protoreflect.Message {
	mi := &file_concord_github_v1_github_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Defaults.ProtoReflect.Descriptor instead.
func (*Defaults) Descriptor() ([]byte, []int) {
	return file_concord_github_v1_github_proto_rawDescGZIP(), []int{1}
}

func (x *Defaults) GetPrivate() bool {
	if x != nil && x.Private != nil {
		return *x.Private
	}
	return false
}

func (x *Defaults) GetDefaultBranch() string {
	if x != nil && x.DefaultBranch != nil {
		return *x.DefaultBranch
	}
	return ""
}

func (x *Defaults) GetAllowAutoMerge() bool {
	if x != nil && x.AllowAutoMerge != nil {
		return *x.AllowAutoMerge
	}
	return false
}

func (x *Defaults) GetAutoDeleteHeadBranches() bool {
	if x != nil && x.AutoDeleteHeadBranches != nil {
		return *x.AutoDeleteHeadBranches
	}
	return false
}

func (x *Defaults) GetProtectedBranches() []*Branch {
	if x != nil {
		return x.ProtectedBranches
	}
	return nil
}

func (x *Defaults) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

type People struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Username string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Teams    []string `protobuf:"bytes,10,rep,name=teams,proto3" json:"teams,omitempty"`
}

func (x *People) Reset() {
	*x = People{}
	if protoimpl.UnsafeEnabled {
		mi := &file_concord_github_v1_github_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *People) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*People) ProtoMessage() {}

func (x *People) ProtoReflect() protoreflect.Message {
	mi := &file_concord_github_v1_github_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use People.ProtoReflect.Descriptor instead.
func (*People) Descriptor() ([]byte, []int) {
	return file_concord_github_v1_github_proto_rawDescGZIP(), []int{2}
}

func (x *People) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *People) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *People) GetTeams() []string {
	if x != nil {
		return x.Teams
	}
	return nil
}

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source       string   `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Destination  string   `protobuf:"bytes,2,opt,name=destination,proto3" json:"destination,omitempty"`
	TargetLabels []string `protobuf:"bytes,3,rep,name=target_labels,json=targetLabels,proto3" json:"target_labels,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_concord_github_v1_github_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_concord_github_v1_github_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_concord_github_v1_github_proto_rawDescGZIP(), []int{3}
}

func (x *File) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *File) GetDestination() string {
	if x != nil {
		return x.Destination
	}
	return ""
}

func (x *File) GetTargetLabels() []string {
	if x != nil {
		return x.TargetLabels
	}
	return nil
}

type Repository struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description *string  `protobuf:"bytes,2,opt,name=description,proto3,oneof" json:"description,omitempty"`
	Archived    *bool    `protobuf:"varint,3,opt,name=archived,proto3,oneof" json:"archived,omitempty"`
	Labels      []string `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty"`
	// Overrides defaults
	Private                *bool     `protobuf:"varint,10,opt,name=private,proto3,oneof" json:"private,omitempty"`
	DefaultBranch          *string   `protobuf:"bytes,11,opt,name=default_branch,json=defaultBranch,proto3,oneof" json:"default_branch,omitempty"`
	AllowAutoMerge         *bool     `protobuf:"varint,12,opt,name=allow_auto_merge,json=allowAutoMerge,proto3,oneof" json:"allow_auto_merge,omitempty"`
	AutoDeleteHeadBranches *bool     `protobuf:"varint,13,opt,name=auto_delete_head_branches,json=autoDeleteHeadBranches,proto3,oneof" json:"auto_delete_head_branches,omitempty"`
	ProtectedBranches      []*Branch `protobuf:"bytes,14,rep,name=protected_branches,json=protectedBranches,proto3" json:"protected_branches,omitempty"`
	Files                  []*File   `protobuf:"bytes,15,rep,name=files,proto3" json:"files,omitempty"` //repeated Secret secrets                   = 16;
}

func (x *Repository) Reset() {
	*x = Repository{}
	if protoimpl.UnsafeEnabled {
		mi := &file_concord_github_v1_github_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Repository) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Repository) ProtoMessage() {}

func (x *Repository) ProtoReflect() protoreflect.Message {
	mi := &file_concord_github_v1_github_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Repository.ProtoReflect.Descriptor instead.
func (*Repository) Descriptor() ([]byte, []int) {
	return file_concord_github_v1_github_proto_rawDescGZIP(), []int{4}
}

func (x *Repository) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Repository) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Repository) GetArchived() bool {
	if x != nil && x.Archived != nil {
		return *x.Archived
	}
	return false
}

func (x *Repository) GetLabels() []string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Repository) GetPrivate() bool {
	if x != nil && x.Private != nil {
		return *x.Private
	}
	return false
}

func (x *Repository) GetDefaultBranch() string {
	if x != nil && x.DefaultBranch != nil {
		return *x.DefaultBranch
	}
	return ""
}

func (x *Repository) GetAllowAutoMerge() bool {
	if x != nil && x.AllowAutoMerge != nil {
		return *x.AllowAutoMerge
	}
	return false
}

func (x *Repository) GetAutoDeleteHeadBranches() bool {
	if x != nil && x.AutoDeleteHeadBranches != nil {
		return *x.AutoDeleteHeadBranches
	}
	return false
}

func (x *Repository) GetProtectedBranches() []*Branch {
	if x != nil {
		return x.ProtectedBranches
	}
	return nil
}

func (x *Repository) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

type Branch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Protection *Protection `protobuf:"bytes,2,opt,name=protection,proto3" json:"protection,omitempty"`
}

func (x *Branch) Reset() {
	*x = Branch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_concord_github_v1_github_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Branch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Branch) ProtoMessage() {}

func (x *Branch) ProtoReflect() protoreflect.Message {
	mi := &file_concord_github_v1_github_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Branch.ProtoReflect.Descriptor instead.
func (*Branch) Descriptor() ([]byte, []int) {
	return file_concord_github_v1_github_proto_rawDescGZIP(), []int{5}
}

func (x *Branch) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Branch) GetProtection() *Protection {
	if x != nil {
		return x.Protection
	}
	return nil
}

type Protection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequirePr      *bool    `protobuf:"varint,1,opt,name=require_pr,json=requirePr,proto3,oneof" json:"require_pr,omitempty"`
	ChecksMustPass *bool    `protobuf:"varint,2,opt,name=checks_must_pass,json=checksMustPass,proto3,oneof" json:"checks_must_pass,omitempty"`
	SignedCommits  *bool    `protobuf:"varint,3,opt,name=signed_commits,json=signedCommits,proto3,oneof" json:"signed_commits,omitempty"`
	RequiredChecks []string `protobuf:"bytes,10,rep,name=required_checks,json=requiredChecks,proto3" json:"required_checks,omitempty"`
}

func (x *Protection) Reset() {
	*x = Protection{}
	if protoimpl.UnsafeEnabled {
		mi := &file_concord_github_v1_github_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Protection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Protection) ProtoMessage() {}

func (x *Protection) ProtoReflect() protoreflect.Message {
	mi := &file_concord_github_v1_github_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Protection.ProtoReflect.Descriptor instead.
func (*Protection) Descriptor() ([]byte, []int) {
	return file_concord_github_v1_github_proto_rawDescGZIP(), []int{6}
}

func (x *Protection) GetRequirePr() bool {
	if x != nil && x.RequirePr != nil {
		return *x.RequirePr
	}
	return false
}

func (x *Protection) GetChecksMustPass() bool {
	if x != nil && x.ChecksMustPass != nil {
		return *x.ChecksMustPass
	}
	return false
}

func (x *Protection) GetSignedCommits() bool {
	if x != nil && x.SignedCommits != nil {
		return *x.SignedCommits
	}
	return false
}

func (x *Protection) GetRequiredChecks() []string {
	if x != nil {
		return x.RequiredChecks
	}
	return nil
}

var File_concord_github_v1_github_proto protoreflect.FileDescriptor

var file_concord_github_v1_github_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x63, 0x6f, 0x6e, 0x63, 0x6f, 0x72, 0x64, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2f, 0x76, 0x31, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x11, 0x63, 0x6f, 0x6e, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x88, 0x02, 0x0a, 0x0c, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x37,
	0x0a, 0x08, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6e, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x08, 0x64,
	0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73,
	0x18, 0x0a, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x12, 0x31, 0x0a,
	0x06, 0x70, 0x65, 0x6f, 0x70, 0x6c, 0x65, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x63, 0x6f, 0x6e, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76,
	0x31, 0x2e, 0x50, 0x65, 0x6f, 0x70, 0x6c, 0x65, 0x52, 0x06, 0x70, 0x65, 0x6f, 0x70, 0x6c, 0x65,
	0x12, 0x41, 0x0a, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73,
	0x18, 0x0c, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x6f, 0x6e, 0x63, 0x6f, 0x72, 0x64,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x69, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x0d, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x22, 0x98, 0x03, 0x0a, 0x08,
	0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x1d, 0x0a, 0x07, 0x70, 0x72, 0x69, 0x76,
	0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x07, 0x70, 0x72, 0x69,
	0x76, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x33, 0x0a, 0x0e, 0x64, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x48, 0x01, 0x52, 0x0d, 0x64, 0x65, 0x66, 0x61,
	0x75, 0x6c, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x10,
	0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x72, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x02, 0x52, 0x0e, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x41,
	0x75, 0x74, 0x6f, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12, 0x3e, 0x0a, 0x19, 0x61,
	0x75, 0x74, 0x6f, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x5f,
	0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x48, 0x03,
	0x52, 0x16, 0x61, 0x75, 0x74, 0x6f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x48, 0x65, 0x61, 0x64,
	0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x73, 0x88, 0x01, 0x01, 0x12, 0x48, 0x0a, 0x12, 0x70,
	0x72, 0x6f, 0x74, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x63, 0x6f, 0x72,
	0x64, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x72, 0x61, 0x6e,
	0x63, 0x68, 0x52, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x65, 0x63, 0x74, 0x65, 0x64, 0x42, 0x72, 0x61,
	0x6e, 0x63, 0x68, 0x65, 0x73, 0x12, 0x2d, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x06,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x05, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x62, 0x72, 0x61,
	0x6e, 0x63, 0x68, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x61, 0x75,
	0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x42, 0x1c, 0x0a, 0x1a, 0x5f, 0x61, 0x75, 0x74,
	0x6f, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x5f, 0x62, 0x72,
	0x61, 0x6e, 0x63, 0x68, 0x65, 0x73, 0x22, 0x60, 0x0a, 0x06, 0x50, 0x65, 0x6f, 0x70, 0x6c, 0x65,
	0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x22, 0x77, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65,
	0x12, 0x1f, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x12, 0x29, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x22, 0xb4, 0x04, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76,
	0x65, 0x64, 0x88, 0x01, 0x01, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x1d, 0x0a,
	0x07, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x48, 0x02,
	0x52, 0x07, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x33, 0x0a, 0x0e,
	0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x48, 0x03, 0x52,
	0x0d, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x88, 0x01,
	0x01, 0x12, 0x2d, 0x0a, 0x10, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x61, 0x75, 0x74, 0x6f, 0x5f,
	0x6d, 0x65, 0x72, 0x67, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x48, 0x04, 0x52, 0x0e, 0x61,
	0x6c, 0x6c, 0x6f, 0x77, 0x41, 0x75, 0x74, 0x6f, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x88, 0x01, 0x01,
	0x12, 0x3e, 0x0a, 0x19, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f,
	0x68, 0x65, 0x61, 0x64, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x73, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x08, 0x48, 0x05, 0x52, 0x16, 0x61, 0x75, 0x74, 0x6f, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x48, 0x65, 0x61, 0x64, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x73, 0x88, 0x01, 0x01,
	0x12, 0x48, 0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x72,
	0x61, 0x6e, 0x63, 0x68, 0x65, 0x73, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63,
	0x6f, 0x6e, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31,
	0x2e, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x52, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x65, 0x63, 0x74,
	0x65, 0x64, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x73, 0x12, 0x2d, 0x0a, 0x05, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x63,
	0x6f, 0x72, 0x64, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x61, 0x72,
	0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x70, 0x72, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x62,
	0x72, 0x61, 0x6e, 0x63, 0x68, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f,
	0x61, 0x75, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x42, 0x1c, 0x0a, 0x1a, 0x5f, 0x61,
	0x75, 0x74, 0x6f, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x5f,
	0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x73, 0x22, 0x6c, 0x0a, 0x06, 0x42, 0x72, 0x61, 0x6e,
	0x63, 0x68, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x45, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x74, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x6f, 0x6e, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x74,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xeb, 0x01, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x74, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x5f, 0x70, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x09, 0x72, 0x65, 0x71,
	0x75, 0x69, 0x72, 0x65, 0x50, 0x72, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x10, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x73, 0x5f, 0x6d, 0x75, 0x73, 0x74, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x0e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x4d, 0x75, 0x73,
	0x74, 0x50, 0x61, 0x73, 0x73, 0x88, 0x01, 0x01, 0x12, 0x2a, 0x0a, 0x0e, 0x73, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x48, 0x02, 0x52, 0x0d, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x73, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64,
	0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x72,
	0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x42, 0x0d, 0x0a,
	0x0b, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x5f, 0x70, 0x72, 0x42, 0x13, 0x0a, 0x11,
	0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x5f, 0x6d, 0x75, 0x73, 0x74, 0x5f, 0x70, 0x61, 0x73,
	0x73, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x63, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x73, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x63, 0x6f,
	0x72, 0x64, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2f, 0x76, 0x31, 0x3b, 0x67, 0x68, 0x5f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_concord_github_v1_github_proto_rawDescOnce sync.Once
	file_concord_github_v1_github_proto_rawDescData = file_concord_github_v1_github_proto_rawDesc
)

func file_concord_github_v1_github_proto_rawDescGZIP() []byte {
	file_concord_github_v1_github_proto_rawDescOnce.Do(func() {
		file_concord_github_v1_github_proto_rawDescData = protoimpl.X.CompressGZIP(file_concord_github_v1_github_proto_rawDescData)
	})
	return file_concord_github_v1_github_proto_rawDescData
}

var file_concord_github_v1_github_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_concord_github_v1_github_proto_goTypes = []interface{}{
	(*Organization)(nil), // 0: concord.github.v1.Organization
	(*Defaults)(nil),     // 1: concord.github.v1.Defaults
	(*People)(nil),       // 2: concord.github.v1.People
	(*File)(nil),         // 3: concord.github.v1.File
	(*Repository)(nil),   // 4: concord.github.v1.Repository
	(*Branch)(nil),       // 5: concord.github.v1.Branch
	(*Protection)(nil),   // 6: concord.github.v1.Protection
}
var file_concord_github_v1_github_proto_depIdxs = []int32{
	1, // 0: concord.github.v1.Organization.defaults:type_name -> concord.github.v1.Defaults
	2, // 1: concord.github.v1.Organization.people:type_name -> concord.github.v1.People
	4, // 2: concord.github.v1.Organization.repositories:type_name -> concord.github.v1.Repository
	5, // 3: concord.github.v1.Defaults.protected_branches:type_name -> concord.github.v1.Branch
	3, // 4: concord.github.v1.Defaults.files:type_name -> concord.github.v1.File
	5, // 5: concord.github.v1.Repository.protected_branches:type_name -> concord.github.v1.Branch
	3, // 6: concord.github.v1.Repository.files:type_name -> concord.github.v1.File
	6, // 7: concord.github.v1.Branch.protection:type_name -> concord.github.v1.Protection
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_concord_github_v1_github_proto_init() }
func file_concord_github_v1_github_proto_init() {
	if File_concord_github_v1_github_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_concord_github_v1_github_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Organization); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_concord_github_v1_github_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Defaults); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_concord_github_v1_github_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*People); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_concord_github_v1_github_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*File); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_concord_github_v1_github_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Repository); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_concord_github_v1_github_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Branch); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_concord_github_v1_github_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Protection); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_concord_github_v1_github_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_concord_github_v1_github_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_concord_github_v1_github_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_concord_github_v1_github_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_concord_github_v1_github_proto_goTypes,
		DependencyIndexes: file_concord_github_v1_github_proto_depIdxs,
		MessageInfos:      file_concord_github_v1_github_proto_msgTypes,
	}.Build()
	File_concord_github_v1_github_proto = out.File
	file_concord_github_v1_github_proto_rawDesc = nil
	file_concord_github_v1_github_proto_goTypes = nil
	file_concord_github_v1_github_proto_depIdxs = nil
}
