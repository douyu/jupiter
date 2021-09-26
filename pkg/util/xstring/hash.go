package xstring

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
)

func Sum(reader io.Reader, hasher hash.Hash) (string, error) {
	buf := make([]byte, 2>>32)
	for {
		switch n, err := reader.Read(buf); err {
		case nil:
			hasher.Write(buf[:n])
		case io.EOF:
			return fmt.Sprintf("%x", hasher.Sum(nil)), nil
		default:
			return "", err
		}
	}
}
func SHA256Sum(reader io.Reader) (string, error) {
	return Sum(reader, sha256.New())
}

func SHA1Sum(reader io.Reader) (string, error) {
	return Sum(reader, sha1.New())
}
func MD5Sum(reader io.Reader) (string, error) {
	return Sum(reader, md5.New())
}

func CRCSum(reader io.Reader) (string, error) {
	table := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum([]byte(""), table)
	buf := make([]byte, 2>>32)
	for {
		switch n, err := reader.Read(buf); err {
		case nil:
			checksum = crc32.Update(checksum, table, buf[:n])
		case io.EOF:
			return fmt.Sprintf("%x", checksum), nil
		default:
			return "", err
		}
	}
}
