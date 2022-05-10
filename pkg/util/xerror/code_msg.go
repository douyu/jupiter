// Copyright 2022 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xerror

import "google.golang.org/grpc/codes"

// Msgs save the registered Err
var msgs = map[codes.Code]string{}

var (
	OK                 = New(0, "请求正常")
	Canceled           = New(1, "操作取消")
	Unknown            = New(2, "未知错误")
	InvalidArgument    = New(3, "无效参数")
	DeadlineExceeded   = New(4, "处理时间超过最后期限")
	NotFound           = New(5, "访问链接不存在")
	AlreadyExists      = New(6, "目标已存在")
	PermissionDenied   = New(7, "权限不足")
	ResourceExhausted  = New(8, "资源耗尽")
	FailedPrecondition = New(9, "前置条件出错")
	Aborted            = New(10, "操作中途失败")
	OutOfRange         = New(11, "操作超出有效范围")
	Unimplemented      = New(12, "当前服务未实现")
	Internal           = New(13, "服务内部异常")
	Unavailable        = New(14, "服务当前不可用")
	DataLoss           = New(15, "数据丢失")
	Unauthenticated    = New(16, "未授权错误")
)

// GRPCCodeFromStatus converts an Err code into the corresponding gRPC response status.
func GRPCCodeFromeErrs(code int32) codes.Code {
	switch code {
	case OK.Ecode:
		return codes.OK
	case Canceled.Ecode:
		return codes.Canceled
	case Unknown.Ecode:
		return codes.Unknown
	case InvalidArgument.Ecode:
		return codes.InvalidArgument
	case DeadlineExceeded.Ecode:
		return codes.DeadlineExceeded
	case NotFound.Ecode:
		return codes.NotFound
	case AlreadyExists.Ecode:
		return codes.AlreadyExists
	case PermissionDenied.Ecode:
		return codes.PermissionDenied
	case ResourceExhausted.Ecode:
		return codes.ResourceExhausted
	case FailedPrecondition.Ecode:
		return codes.FailedPrecondition
	case Aborted.Ecode:
		return codes.Aborted
	case OutOfRange.Ecode:
		return codes.OutOfRange
	case Unimplemented.Ecode:
		return codes.Unimplemented
	case Internal.Ecode:
		return codes.Internal
	case Unavailable.Ecode:
		return codes.Unavailable
	case DataLoss.Ecode:
		return codes.DataLoss
	case Unauthenticated.Ecode:
		return codes.Unauthenticated
	}
	return codes.Unknown
}

// StatusFromGRPCCode converts a gRPC error code into the corresponding Err code.
func ErrsFromGRPCCode(code codes.Code) int32 {
	switch code {
	case codes.OK:
		return OK.Ecode
	case codes.Canceled:
		return Canceled.Ecode
	case codes.Unknown:
		return Unknown.Ecode
	case codes.InvalidArgument:
		return InvalidArgument.Ecode
	case codes.DeadlineExceeded:
		return DeadlineExceeded.Ecode
	case codes.NotFound:
		return NotFound.Ecode
	case codes.AlreadyExists:
		return AlreadyExists.Ecode
	case codes.PermissionDenied:
		return PermissionDenied.Ecode
	case codes.ResourceExhausted:
		return ResourceExhausted.Ecode
	case codes.FailedPrecondition:
		return FailedPrecondition.Ecode
	case codes.Aborted:
		return Aborted.Ecode
	case codes.OutOfRange:
		return OutOfRange.Ecode
	case codes.Unimplemented:
		return Unimplemented.Ecode
	case codes.Internal:
		return Internal.Ecode
	case codes.Unavailable:
		return Unavailable.Ecode
	case codes.DataLoss:
		return DataLoss.Ecode
	case codes.Unauthenticated:
		return Unauthenticated.Ecode
	}
	return UnknownCode
}
