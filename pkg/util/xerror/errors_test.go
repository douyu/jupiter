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

import (
	"errors"
	"testing"

	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestErr(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		error := Unknown.WithMsg("test error")
		assert.Equal(t, error.Error(), "xerror: ecode = 2 msg = test error")
	})
	t.Run("nil", func(t *testing.T) {
		var err *Err
		assert.Equal(t, err.Error(), "nil")
	})
	t.Run("New", func(t *testing.T) {
		_ = New(100, "test error")
		assert.Panics(t, func() { New(1, "test error") })
		assert.Panics(t, func() { New(100, "test error") })
	})
	t.Run("GetEcode", func(t *testing.T) {
		error1 := Unknown.WithMsg("test error")
		assert.Equal(t, int32(UnknownCode), error1.GetEcode())

		var error2 *Err
		assert.Equal(t, error2.GetEcode(), int32(0))
	})

	t.Run("GetMsg", func(t *testing.T) {
		error1 := Unknown.WithMsg("test error")
		assert.Equal(t, error1.GetMsg(), "test error")

		var error2 *Err
		assert.Equal(t, error2.GetMsg(), "")
	})
	t.Run("GetData", func(t *testing.T) {
		error1 := Unknown.WithMsg("test error")
		assert.Equal(t, error1.GetData(), struct{}{})

		var error2 *Err
		assert.Equal(t, error2.GetData(), nil)
	})
	t.Run("WithMsg1", func(t *testing.T) {
		error1 := Unknown.WithMsg("test error")
		error2 := error1.WithMsg("new error")
		assert.Equal(t, error1.Msg, "test error")
		assert.Equal(t, error2.Msg, "new error")
	})
	t.Run("WithData", func(t *testing.T) {
		error1 := Unknown.WithMsg("test error").WithData("test data")
		error2 := error1.WithData("new data")
		assert.Equal(t, error1.Data, "test data")
		assert.Equal(t, error2.Data, "new data")
	})
	t.Run("Convert", func(t *testing.T) {
		var err error
		var Err *Err
		assert.Equal(t, Err, Convert(err))

		error1 := Unknown.WithMsg("test error")
		assert.Equal(t, error1, Convert(error1))

		s := status.New(codes.Unknown, error1.Msg)
		assert.Equal(t, error1, Convert(s.Err()))

		error2 := errors.New("test error")
		assert.Equal(t, error1, Convert(error2))
	})
	t.Run("GRPCStatus", func(t *testing.T) {
		error1 := Unknown.WithMsg("test error")
		s := status.New(codes.Unknown, error1.Msg)
		assert.Equal(t, s, error1.GRPCStatus())
	})
	t.Run("CodeConvert", func(t *testing.T) {
		assert.Equal(t, GRPCCodeFromeErrs(OK.Ecode), codes.OK)
		assert.Equal(t, ErrsFromGRPCCode(codes.OK), OK.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(Canceled.Ecode), codes.Canceled)
		assert.Equal(t, ErrsFromGRPCCode(codes.Canceled), Canceled.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(Unknown.Ecode), codes.Unknown)
		assert.Equal(t, ErrsFromGRPCCode(codes.Unknown), Unknown.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(InvalidArgument.Ecode), codes.InvalidArgument)
		assert.Equal(t, ErrsFromGRPCCode(codes.InvalidArgument), InvalidArgument.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(DeadlineExceeded.Ecode), codes.DeadlineExceeded)
		assert.Equal(t, ErrsFromGRPCCode(codes.DeadlineExceeded), DeadlineExceeded.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(NotFound.Ecode), codes.NotFound)
		assert.Equal(t, ErrsFromGRPCCode(codes.NotFound), NotFound.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(AlreadyExists.Ecode), codes.AlreadyExists)
		assert.Equal(t, ErrsFromGRPCCode(codes.AlreadyExists), AlreadyExists.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(PermissionDenied.Ecode), codes.PermissionDenied)
		assert.Equal(t, ErrsFromGRPCCode(codes.PermissionDenied), PermissionDenied.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(ResourceExhausted.Ecode), codes.ResourceExhausted)
		assert.Equal(t, ErrsFromGRPCCode(codes.ResourceExhausted), ResourceExhausted.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(FailedPrecondition.Ecode), codes.FailedPrecondition)
		assert.Equal(t, ErrsFromGRPCCode(codes.FailedPrecondition), FailedPrecondition.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(Aborted.Ecode), codes.Aborted)
		assert.Equal(t, ErrsFromGRPCCode(codes.Aborted), Aborted.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(OutOfRange.Ecode), codes.OutOfRange)
		assert.Equal(t, ErrsFromGRPCCode(codes.OutOfRange), OutOfRange.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(Unimplemented.Ecode), codes.Unimplemented)
		assert.Equal(t, ErrsFromGRPCCode(codes.Unimplemented), Unimplemented.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(Internal.Ecode), codes.Internal)
		assert.Equal(t, ErrsFromGRPCCode(codes.Internal), Internal.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(Unavailable.Ecode), codes.Unavailable)
		assert.Equal(t, ErrsFromGRPCCode(codes.Unavailable), Unavailable.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(DataLoss.Ecode), codes.DataLoss)
		assert.Equal(t, ErrsFromGRPCCode(codes.DataLoss), DataLoss.Ecode)
		assert.Equal(t, GRPCCodeFromeErrs(Unauthenticated.Ecode), codes.Unauthenticated)
		assert.Equal(t, ErrsFromGRPCCode(codes.Unauthenticated), Unauthenticated.Ecode)

		assert.Equal(t, codes.Unknown, GRPCCodeFromeErrs(int32(100)))
		assert.Equal(t, int32(UnknownCode), ErrsFromGRPCCode(codes.Code(100)))
	})
}
