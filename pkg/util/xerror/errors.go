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
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	UnknownCode = 2
)

// Err struct
type Err struct {
	Ecode int32       `json:"error"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// Error implements the error interface.
func (e *Err) Error() string {
	if e == nil {
		return "nil"
	}

	return fmt.Sprintf("xerror: ecode = %v msg = %s", e.GetEcode(), e.GetMsg())
}

// New returns an error object for the code, message.
func New(code codes.Code, message string) *Err {
	if _, ok := msgs[code]; ok {
		panic("code had been registered")
	}
	msgs[code] = message

	return &Err{
		Ecode: int32(code),
		Msg:   message,
		Data:  struct{}{},
	}
}

// GetEcode return the packaged Ecode
func (e *Err) GetEcode() int32 {

	if e != nil && e.Ecode != 0 {
		return e.Ecode
	}
	return 0
}

// GetMsg return the Msg
func (e *Err) GetMsg() string {
	if e != nil {
		return e.Msg
	}
	return ""
}

// GetData return the Data
func (e *Err) GetData() interface{} {
	if e != nil {
		return e.Data
	}
	return nil
}

// WithMsg allows the programmer to override Msg and ruturn the new Err
func (e *Err) WithMsg(msg string) *Err {
	return &Err{
		Ecode: e.Ecode,
		Msg:   msg,
		Data:  e.Data,
	}
}

// WithData allows the programmer to override Data and ruturn the new Err
func (e *Err) WithData(data interface{}) *Err {
	return &Err{
		Ecode: e.Ecode,
		Msg:   e.Msg,
		Data:  data,
	}
}

// Convert try to convert an error to *Error.
// It supports wrapped errors.
func Convert(err error) *Err {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Err); ok {
		return e
	}
	// 将status.error转化为Err
	gs, ok := status.FromError(err)
	if ok {
		return &Err{
			Ecode: ErrsFromGRPCCode(gs.Code()),
			Msg:   gs.Message(),
			Data:  struct{}{},
		}
	}
	return &Err{
		Ecode: int32(UnknownCode),
		Msg:   err.Error(),
		Data:  struct{}{},
	}
}

// GRPCStatus returns the Status represented by se.
func (e *Err) GRPCStatus() *status.Status {
	s := status.New(GRPCCodeFromeErrs(e.Ecode), e.Msg)
	return s
}
