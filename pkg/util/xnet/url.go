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

package xnet

import (
	"net/url"
	"strconv"
	"time"
)

// URL wrap url.URL.
type URL struct {
	url.URL
}

// ParseURL parses raw into URL.
func ParseURL(raw string) (*URL, error) {
	u, e := url.Parse(raw)
	if e != nil {
		return nil, e
	}

	return &URL{
		URL: *u,
	}, nil
}

// Password gets password from URL.
func (u *URL) Password() (string, bool) {
	if u.User != nil {
		return u.User.Password()
	}
	return "", false
}

// Username gets username from URL.
func (u *URL) Username() string {
	return u.User.Username()
}

// QueryInt returns provided field's value in int type.
func (u *URL) QueryInt(field string, expect int) (ret int) {
	ret = expect
	if mi := u.Query().Get(field); mi != "" {
		if m, e := strconv.Atoi(mi); e == nil {
			if m > 0 {
				ret = m
			}
		}
	}

	return
}

// QueryInt64 returns provided field's value in int64 type.
func (u *URL) QueryInt64(field string, expect int64) (ret int64) {
	ret = expect
	if mi := u.Query().Get(field); mi != "" {
		if m, e := strconv.ParseInt(mi, 10, 64); e == nil {
			if m > 0 {
				ret = m
			}
		}
	}

	return
}

// QueryString returns provided field's value in string type.
func (u *URL) QueryString(field string, expect string) (ret string) {
	ret = expect
	if mi := u.Query().Get(field); mi != "" {
		ret = mi
	}

	return
}

// QuerySecond returns provided field's value in duration type.
// Deprecated: use QueryDuration instead.
func (u *URL) QuerySecond(field string, expect int64) (ret time.Duration) {
	return u.QueryDuration(field, expect)
}

// QueryDuration returns provided field's value in duration type.
func (u *URL) QueryDuration(field string, expect int64) (ret time.Duration) {
	ret = time.Duration(expect)
	if mi := u.Query().Get(field); mi != "" {
		if m, e := strconv.ParseInt(mi, 10, 64); e == nil {
			if m > 0 {
				ret = time.Duration(m)
			}
		}
	}

	return
}
