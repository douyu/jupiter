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

package xgin

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Server(t *testing.T) {
	c := DefaultConfig()
	c.Port = 0
	s := c.MustBuild()
	stoped := make(chan struct{}, 1)
	go func() {
		time.AfterFunc(time.Second, func() {
			stoped <- struct{}{}
		})
		assert.True(t, s.Healthz())
		assert.NotNil(t, s.Info())
		s.Serve()
	}()
	<-stoped
}
