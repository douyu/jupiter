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

package tests

import (
	"net/http"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

type HTTPTestCase struct {
	Host              string
	Method            string
	Path              string
	Body              string
	Header            map[string]string
	Query             string
	ExpectStatusRange httpexpect.StatusRange
	ExpectBody        string
}

// RunHTTPTestCase runs a test case against the given handler.
func RunHTTPTestCase(htc HTTPTestCase) {
	ginkgoT := ginkgo.GinkgoT()
	expect := httpexpect.New(ginkgoT, htc.Host)

	req := &httpexpect.Request{}
	req.WithTimeout(time.Second)

	switch htc.Method {
	case http.MethodGet:
		req = expect.GET(htc.Path)
	case http.MethodPost:
		req = expect.POST(htc.Path)
	case http.MethodPut:
		req = expect.PUT(htc.Path)
	case http.MethodDelete:
		req = expect.DELETE(htc.Path)
	case http.MethodOptions:
		req = expect.OPTIONS(htc.Path)
	}

	assert.NotNil(ginkgoT, req)

	if len(htc.Query) > 0 {
		req.WithQueryString(htc.Query)
	}

	if len(htc.Body) > 0 {
		req.WithText(htc.Body)
	}

	resp := req.Expect()

	if htc.ExpectStatusRange > 0 {
		resp.StatusRange(htc.ExpectStatusRange)
	}

	if len(htc.ExpectBody) == 0 {
		resp.Body().Empty()
	} else {
		resp.Body().Contains(htc.ExpectBody)
	}
}
