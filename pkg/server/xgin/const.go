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

package xgin

//code
const (
	codeMS                   = 1000
	codeMSInvalidParam       = 1001
	codeMSInvoke             = 1002
	codeMSInvokeLen          = 1003
	codeMSSecondItemNotError = 1004
	codeMSResErr             = 1005
)
const (
	// StatusContinue ...
	StatusContinue = 100
	// StatusSwitchingProtocols ...
	StatusSwitchingProtocols = 101

	// StatusOK ...
	StatusOK = 200
	// StatusCreated ...
	StatusCreated = 201
	// StatusAccepted ...
	StatusAccepted = 202
	// StatusNonAuthoritativeInfo ...
	StatusNonAuthoritativeInfo = 203
	// StatusNoContent ...
	StatusNoContent = 204
	// StatusResetContent ...
	StatusResetContent = 205
	// StatusPartialContent ...
	StatusPartialContent = 206

	// StatusMultipleChoices ...
	StatusMultipleChoices = 300
	// StatusMovedPermanently ...
	StatusMovedPermanently = 301
	// StatusFound ...
	StatusFound = 302
	// StatusSeeOther ...
	StatusSeeOther = 303
	// StatusNotModified ...
	StatusNotModified = 304
	// StatusUseProxy ...
	StatusUseProxy = 305
	// StatusTemporaryRedirect ...
	StatusTemporaryRedirect = 307

	// StatusBadRequest ...
	StatusBadRequest = 400
	// StatusUnauthorized ...
	StatusUnauthorized = 401
	// StatusPaymentRequired ...
	StatusPaymentRequired = 402
	// StatusForbidden ...
	StatusForbidden = 403
	// StatusNotFound ...
	StatusNotFound = 404
	// StatusMethodNotAllowed ...
	StatusMethodNotAllowed = 405
	// StatusNotAcceptable ...
	StatusNotAcceptable = 406
	// StatusProxyAuthRequired ...
	StatusProxyAuthRequired = 407
	// StatusRequestTimeout ...
	StatusRequestTimeout = 408
	// StatusConflict ...
	StatusConflict = 409
	// StatusGone ...
	StatusGone = 410
	// StatusLengthRequired ...
	StatusLengthRequired = 411
	// StatusPreconditionFailed ...
	StatusPreconditionFailed = 412
	// StatusRequestEntityTooLarge ...
	StatusRequestEntityTooLarge = 413
	// StatusRequestURITooLong ...
	StatusRequestURITooLong = 414
	// StatusUnsupportedMediaType ...
	StatusUnsupportedMediaType = 415
	// StatusRequestedRangeNotSatisfiable ...
	StatusRequestedRangeNotSatisfiable = 416
	// StatusExpectationFailed ...
	StatusExpectationFailed = 417
	// StatusTeapot ...
	StatusTeapot = 418
	// StatusPreconditionRequired ...
	StatusPreconditionRequired = 428
	// StatusTooManyRequests ...
	StatusTooManyRequests = 429
	// StatusRequestHeaderFieldsTooLarge ...
	StatusRequestHeaderFieldsTooLarge = 431
	// StatusUnavailableForLegalReasons ...
	StatusUnavailableForLegalReasons = 451
	// StatusInternalServerError ...
	StatusInternalServerError = 500
	// StatusNotImplemented ...
	StatusNotImplemented = 501
	// StatusBadGateway ...
	StatusBadGateway = 502
	// StatusServiceUnavailable ...
	StatusServiceUnavailable = 503
	// StatusGatewayTimeout ...
	StatusGatewayTimeout = 504
	// StatusHTTPVersionNotSupported ...
	StatusHTTPVersionNotSupported = 505
	// StatusNetworkAuthenticationRequired ...
	StatusNetworkAuthenticationRequired = 511

	// StatusErrorCodeReturned 针对微服务定制的错误返回status
	StatusErrorCodeReturned = 800
)

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}

var statusText = map[int]string{
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "Switching Protocols",

	StatusOK:                   "OK",
	StatusCreated:              "Created",
	StatusAccepted:             "Accepted",
	StatusNonAuthoritativeInfo: "Non-Authoritative Information",
	StatusNoContent:            "No Content",
	StatusResetContent:         "Reset Content",
	StatusPartialContent:       "Partial Content",

	StatusMultipleChoices:   "Multiple Choices",
	StatusMovedPermanently:  "Moved Permanently",
	StatusFound:             "Found",
	StatusSeeOther:          "See Other",
	StatusNotModified:       "Not Modified",
	StatusUseProxy:          "Use Proxy",
	StatusTemporaryRedirect: "Temporary Redirect",

	StatusBadRequest:                   "Bad Request",
	StatusUnauthorized:                 "Unauthorized",
	StatusPaymentRequired:              "Payment Required",
	StatusForbidden:                    "Forbidden",
	StatusNotFound:                     "Not Found",
	StatusMethodNotAllowed:             "Method Not Allowed",
	StatusNotAcceptable:                "Not Acceptable",
	StatusProxyAuthRequired:            "Proxy Authentication Required",
	StatusRequestTimeout:               "Request Timeout",
	StatusConflict:                     "Conflict",
	StatusGone:                         "Gone",
	StatusLengthRequired:               "Length Required",
	StatusPreconditionFailed:           "Precondition Failed",
	StatusRequestEntityTooLarge:        "Request Entity Too Large",
	StatusRequestURITooLong:            "Request URI Too Long",
	StatusUnsupportedMediaType:         "Unsupported Media Type",
	StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	StatusExpectationFailed:            "Expectation Failed",
	StatusTeapot:                       "I'm a teapot",
	StatusPreconditionRequired:         "Precondition Required",
	StatusTooManyRequests:              "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	StatusInternalServerError:           "Internal Server Errord",
	StatusNotImplemented:                "Not Implemented",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "Service Unavailable",
	StatusGatewayTimeout:                "Gateway Timeout",
	StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

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
