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

package xcrypt

import (
	"encoding/base64"
	"testing"
)

// 私钥生成
//openssl genrsa -out rsa_private_key.pem 1024
var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCxZNKFvMK5XtATIHfJyJEYB2JjxGsun2MkEX6/i+UINN5lmBwe
EjEeW8WtplXFHdHhOjRxPtk21tRRWm/6lvwUyli3JAPJQkgxj5US6rEie5FC8fxq
XBV08eoKW138vz2bt1Ls/F027s0kDv9tat61lAQ7sxqaW/Ftz1MxOSkucwIDAQAB
AoGABAWzOFEVYTqjISvlS2/+yjqwom57t6zphJHY++LiKJN6T3dpe80RzAxsqQlS
fIu2jJLTSZYROssYOVgBnf76bDyRV7NQDfVt0FE1urF5XyHo9qSHGpb7XqcgIidD
UBHKcOdrcyNjvxnZoZHQHw43fdkyb1WWrxuyK6T/KqW5mIECQQDlu9B+tLBHvlMO
Sp/skCzfi95kPZJ5H5gZXQadz/C0gIWh1P0Z5U+cduBEdq7xO4XpfGF18X8QsNFO
jZmJSeHhAkEAxa0MbumUwsLDdDCf69tBiLW069new0TTzWPK/mwNOuM2JqiQpPx4
ckEAsyh3GGFdt5MhjniT0hhppzi7tghC0wJAEIna0qRTZHbRJ+A7bx5Z/KXnFrRQ
DSQ3IOxPg6DqpTPzatkYd3rIpmzwbD1XDsrIMyzfH0yJZzwzdUJAYV/OQQJANmIL
X6AnawWGHDscZBjoCKJk6dYAsRwIYSMpP6GeaisERNJvKNTEljpH5QIm8bAnxk9W
FgoaMNzChFzZV5UiPQJAA2YWGrFR/sma24PBILK7LQQO3IkDFuQR6/23KQH8fU9j
rOv1pjAeZCw7NMETaAqtIxRcLDqdhIGe/C/JS8jPkA==
-----END RSA PRIVATE KEY-----
`)

// 公钥: 根据私钥生成
//openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCxZNKFvMK5XtATIHfJyJEYB2Jj
xGsun2MkEX6/i+UINN5lmBweEjEeW8WtplXFHdHhOjRxPtk21tRRWm/6lvwUyli3
JAPJQkgxj5US6rEie5FC8fxqXBV08eoKW138vz2bt1Ls/F027s0kDv9tat61lAQ7
sxqaW/Ftz1MxOSkucwIDAQAB
-----END PUBLIC KEY-----
`)

func TestRsaEncrypt(t *testing.T) {
	str := "jupiter"
	b, err := RsaEncrypt([]byte(str), publicKey)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(base64.StdEncoding.EncodeToString(b))
}

func TestRsaDecrypt(t *testing.T) {
	str := "jupiter"
	b, err := RsaEncrypt([]byte(str), publicKey)
	if err != nil {
		t.Error(err.Error())
	}
	ss, err := RsaDecrypt(b, privateKey)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(string(ss))
}
