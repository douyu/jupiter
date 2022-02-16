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

package xstring

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Addr2Hex converts address string to hex string, only support ipv4.
func Addr2Hex(str string) (string, error) {
	ipStr, portStr, err := net.SplitHostPort(str)
	if err != nil {
		return "", err
	}

	ip := net.ParseIP(ipStr).To4()
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return "", nil
	}

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(port))
	ip = append(ip, buf...)

	return hex.EncodeToString(ip), nil
}

// Hex2Addr converts hex string to address.
func Hex2Addr(str string) (string, error) {
	buf, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}
	if len(buf) < 4 {
		return "", fmt.Errorf("bad hex string length")
	}
	return fmt.Sprintf("%s:%d", net.IP(buf[:4]).String(), binary.BigEndian.Uint16(buf[4:])), nil
}

// Strings ...
type Strings []string

// KickEmpty kick empty elements from ss
func KickEmpty(ss []string) Strings {
	var ret = make([]string, 0)
	for _, str := range ss {
		if str != "" {
			ret = append(ret, str)
		}
	}
	return Strings(ret)
}

func Kick(ss []string, remove func(item string) bool) Strings {
	for i := 0; i < len(ss); i++ {
		if remove(ss[i]) {
			ss = append(ss[:i], ss[i+1:]...)
			i--
		}
	}
	return ss
}

// AnyBlank return true if ss has empty element
func AnyBlank(ss []string) bool {
	for _, str := range ss {
		if str == "" {
			return true
		}
	}

	return false
}

func Any(ss []string, match func(item string) bool) bool {
	for _, str := range ss {
		if match(str) {
			return true
		}
	}
	return false
}

func All(ss []string, match func(item string) bool) bool {
	if len(ss) == 0 {
		return false
	}
	for _, str := range ss {
		if !match(str) {
			return false
		}
	}
	return true
}

// HeadT ...
func (ss Strings) HeadT() (string, Strings) {
	if len(ss) > 0 {
		return ss[0], Strings(ss[1:])
	}

	return "", Strings{}
}

// Head ...
func (ss Strings) Head() string {
	if len(ss) > 0 {
		return ss[0]
	}
	return ""
}

// Head2 ...
func (ss Strings) Head2() (h0, h1 string) {
	if len(ss) > 0 {
		h0 = ss[0]
	}
	if len(ss) > 1 {
		h1 = ss[1]
	}
	return
}

// Head3 ...
func (ss Strings) Head3() (h0, h1, h2 string) {
	if len(ss) > 0 {
		h0 = ss[0]
	}
	if len(ss) > 1 {
		h1 = ss[1]
	}
	if len(ss) > 2 {
		h2 = ss[2]
	}
	return
}

// Head4 ...
func (ss Strings) Head4() (h0, h1, h2, h3 string) {
	if len(ss) > 0 {
		h0 = ss[0]
	}
	if len(ss) > 1 {
		h1 = ss[1]
	}
	if len(ss) > 2 {
		h2 = ss[2]
	}
	if len(ss) > 3 {
		h3 = ss[3]
	}
	return
}

// Split ...
func Split(raw string, sep string) Strings {
	return Strings(strings.Split(raw, sep))
}

func Pointer(str string) *string {
	return &str
}

func IndexString(ss []string, str string) int {
	for ind, s := range ss {
		if str == s {
			return ind
		}
	}
	return -1
}

func HasString(ss []string, str string) bool {
	return IndexString(ss, str) > -1
}

func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}
