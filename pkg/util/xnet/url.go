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
	"time"

	"github.com/douyu/jupiter/pkg/util/xcast"
)

// URL wrap url.URL.
type URL struct {
	Scheme     string
	Opaque     string        // encoded opaque data
	User       *url.Userinfo // username and password information
	Host       string        // host or host:port
	Path       string        // path (relative paths may omit leading slash)
	RawPath    string        // encoded path hint (see EscapedPath method)
	ForceQuery bool          // append a query ('?') even if RawQuery is empty
	RawQuery   string        // encoded query values, without '?'
	Fragment   string        // fragment for references, without '#'
	HostName   string
	Port       string
	params     url.Values
}

// ParseURL parses raw into URL.
func ParseURL(raw string) (*URL, error) {
	u, e := url.Parse(raw)
	if e != nil {
		return nil, e
	}

	return &URL{
		Scheme:     u.Scheme,
		Opaque:     u.Opaque,
		User:       u.User,
		Host:       u.Host,
		Path:       u.Path,
		RawPath:    u.RawPath,
		ForceQuery: u.ForceQuery,
		RawQuery:   u.RawQuery,
		Fragment:   u.Fragment,
		HostName:   u.Hostname(),
		Port:       u.Port(),
		params:     u.Query(),
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
// if value is empty, expect returns
func (u *URL) QueryInt(field string, expect int) (ret int) {
	ret, err := xcast.ToIntE(u.Query().Get(field))
	if err != nil {
		return expect
	}

	return ret
}

// QueryInt64 returns provided field's value in int64 type.
// if value is empty, expect returns
func (u *URL) QueryInt64(field string, expect int64) (ret int64) {
	ret, err := xcast.ToInt64E(u.Query().Get(field))
	if err != nil {
		return expect
	}

	return ret
}

// QueryString returns provided field's value in string type.
// if value is empty, expect returns
func (u *URL) QueryString(field string, expect string) (ret string) {
	ret = expect
	if mi := u.Query().Get(field); mi != "" {
		ret = mi
	}

	return
}

// QueryDuration returns provided field's value in duration type.
// if value is empty, expect returns
func (u *URL) QueryDuration(field string, expect time.Duration) (ret time.Duration) {
	ret, err := xcast.ToDurationE(u.Query().Get(field))
	if err != nil {
		return expect
	}

	return ret
}

// QueryBool returns provided field's value in bool
// if value is empty, expect returns
func (u *URL) QueryBool(field string, expect bool) (ret bool) {
	ret, err := xcast.ToBoolE(u.Query().Get(field))
	if err != nil {
		return expect
	}
	return ret
}

// Query parses RawQuery and returns the corresponding values.
// It silently discards malformed value pairs.
// To check errors use ParseQuery.
func (u *URL) Query() url.Values {
	v, _ := url.ParseQuery(u.RawQuery)
	return v
}
