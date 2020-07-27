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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)
// encrypt aes
func AesEncrypt(ori []byte, key string) (string, error) {

	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()

	ori = PKCS7Padding(ori, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])

	crypt := make([]byte, len(ori))

	blockMode.CryptBlocks(crypt, ori)

	return base64.RawURLEncoding.EncodeToString(crypt), nil
}

// decrypt aes
func AesDecrypt(crypt string, key string) (string, error) {

	cryptByte, err := base64.RawURLEncoding.DecodeString(crypt)
	if err != nil {
		return "", err
	}
	k := []byte(key)

	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()

	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])

	orig := make([]byte, len(cryptByte))

	blockMode.CryptBlocks(orig, cryptByte)

	orig = PKCS7UnPadding(orig)
	return string(orig), nil
}

func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}