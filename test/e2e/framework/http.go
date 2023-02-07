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

package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	jresty "github.com/douyu/jupiter/pkg/client/resty"
	"github.com/imdario/mergo"
	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

type HTTPTestCase struct {
	Conf         *jresty.Config
	Method       string
	Path         string
	Body         string
	Timeout      time.Duration
	Header       map[string]string
	Query        string
	ExpectHeader http.Header
	ExpectStatus int
	ExpectBody   string
}

// RunHTTPTestCase runs a test case against the given handler.
func RunHTTPTestCase(htc HTTPTestCase) {
	ginkgoT := ginkgo.GinkgoT()

	if htc.Timeout == 0 {
		htc.Timeout = time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), htc.Timeout)
	defer cancel()

	err := mergo.Merge(htc.Conf, jresty.DefaultConfig())
	assert.Nil(ginkgoT, err)

	req := htc.Conf.MustBuild().R()
	req.SetQueryString(htc.Query)
	req.SetBody(htc.Body)
	req.SetHeaders(htc.Header)
	req.SetContext(ctx)

	res, err := req.Execute(htc.Method, htc.Path)

	assert.Nil(ginkgoT, err, "error: %s", err)

	if htc.ExpectStatus > 0 {
		assert.Equal(ginkgoT, htc.ExpectStatus, res.StatusCode(),
			"expected: %d\nactually: %d", htc.ExpectStatus, res.StatusCode())
	}

	if len(htc.ExpectHeader) > 0 {
		assert.EqualValues(ginkgoT, htc.ExpectHeader, res.Header(),
			"expected: %s\nactually: %s", htc.ExpectHeader, res.Header())
	}

	if len(htc.ExpectBody) > 0 {
		var body bytes.Buffer
		err = json.Compact(&body, []byte(res.String()))
		// 如果Compact失败，则说明不是json格式
		if err != nil {
			body.WriteString(res.String())
		}

		assert.Equal(ginkgoT, htc.ExpectBody, body.String(),
			"expected: %s\nactually: %s", htc.ExpectBody, body.String())
	}
}
