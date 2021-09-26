// Copyright 2021 rex lv
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

package autoproc

import (
	"runtime"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	// 初始化注册中心
	if _, err := maxprocs.Set(); err != nil {
		xlog.Panic("auto max procs", xlog.FieldMod(ecode.ModProc), xlog.FieldErrKind(ecode.ErrKindAny), xlog.FieldErr(err))
	}
	conf.OnLoaded(func(c *conf.Configuration) {
		if maxProcs := conf.GetInt("maxProc"); maxProcs != 0 {
			runtime.GOMAXPROCS(maxProcs)
		}
		xlog.Info("auto max procs", xlog.FieldMod(ecode.ModProc), xlog.Int64("procs", int64(runtime.GOMAXPROCS(-1))))
	})
}
