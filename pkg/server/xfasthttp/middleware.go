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

package xfasthttp

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Middleware func(h fasthttp.RequestHandler) fasthttp.RequestHandler

// recoverMiddleware ...
func recoverMiddleware(c *Config) Middleware {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			var err error
			var beg = time.Now()
			var fields = make([]xlog.Field, 0, 8)

			defer func() {
				fields = append(fields, zap.Float64("cost", time.Since(beg).Seconds()))
				if rec := recover(); rec != nil {
					switch rec := rec.(type) {
					case error:
						err = rec
					default:
						err = fmt.Errorf("%v", rec)
					}

					stack := make([]byte, 4096)
					length := runtime.Stack(stack, true)
					fields = append(fields, zap.ByteString("stack", stack[:length]))

					if !c.DisablePrintStack {
						fmt.Fprintln(os.Stderr, "[PANIC RECOVER]", err)
						fmt.Fprintln(os.Stderr, string(stack))
					}
				}

				fields = append(fields,
					zap.String("method", string(ctx.Request.Header.Method())),
					zap.Int("code", ctx.Response.StatusCode()),
					zap.String("host", string(ctx.Host())),
					zap.String("path", string(ctx.Path())),
				)
				if err != nil {
					fields = append(fields, zap.String("err", err.Error()))
					c.logger.Error("access", fields...)
					return
				}

				c.logger.Info("access", fields...)
			}()

			next(ctx)
		}
	}
}
