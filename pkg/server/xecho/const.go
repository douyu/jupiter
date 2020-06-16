// Copyright 2020 Douyu
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

package xecho

//code
const (
	codeMS                   = 1000
	codeMSInvalidParam       = 1001
	codeMSInvoke             = 1002
	codeMSInvokeLen          = 1003
	codeMSSecondItemNotError = 1004
	codeMSResErr             = 1005
)

// Headers
const (
	// HeaderAcceptEncoding ...
	HeaderAcceptEncoding = "Accept-Encoding"
	// HeaderContentType ...
	HeaderContentType = "Content-Type"
	// HRPC Errord
	HeaderHRPCErr = "HRPC-Errord"
)

// MIME types
const (
	// MIMEApplicationJSON ...
	MIMEApplicationJSON = "application/json"
	// MIMEApplicationJSONCharsetUTF8 ...
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
	// MIMEApplicationProtobuf ...
	MIMEApplicationProtobuf = "application/protobuf"
)
const (
	charsetUTF8 = "charset=utf-8"
)
