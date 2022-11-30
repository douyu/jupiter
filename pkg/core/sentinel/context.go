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

package sentinel

import (
	"context"

	"github.com/alibaba/sentinel-golang/core/base"
)

var (
	sentinelEntryKey = struct{}{}
)

func WithContext(ctx context.Context, val *base.SentinelEntry) context.Context {
	return context.WithValue(ctx, sentinelEntryKey, val)
}

func FromContext(ctx context.Context) *base.SentinelEntry {
	if ctx == nil {
		return nil
	}
	if val, ok := ctx.Value(sentinelEntryKey).(*base.SentinelEntry); ok {
		return val
	}
	return nil
}
