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
	"context"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

type HTTPTestCase struct {
	Host         string
	Method       string
	Path         string
	Body         string
	Timeout      time.Duration
	Header       map[string]string
	Query        string
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

	req := resty.New().R()
	req.SetQueryString(htc.Query)
	req.SetBody(htc.Body)
	req.SetContext(ctx)

	res, err := req.Execute(htc.Method, htc.Host+htc.Path)

	assert.Nil(ginkgoT, err)

	if htc.ExpectStatus > 0 {
		assert.Equal(ginkgoT, htc.ExpectStatus, res.StatusCode())
	}

	assert.Equal(ginkgoT, htc.ExpectBody, res.String())
}
