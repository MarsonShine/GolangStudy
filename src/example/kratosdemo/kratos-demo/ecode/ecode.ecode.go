// Code generated by protoc-gen-ecode v0.1, DO NOT EDIT.
// source: ecode/ecode.proto

package ecode

import (
	"github.com/go-kratos/kratos/pkg/ecode"
)

// to suppressed 'imported but not used warning'
var _ ecode.Codes

// UserErrCode ecode
var (
	UserNotNull  = ecode.New(404)
	UserNotLogin = ecode.New(401)
)
