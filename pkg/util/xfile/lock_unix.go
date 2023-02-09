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

//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd
// +build darwin dragonfly freebsd linux netbsd openbsd

package xfile

import (
	"io"
	"os"
	"syscall"
)

// lockCloser hides all of an os.File's methods, except for Close.
type lockCloser struct {
	f *os.File
}

// Close ...
func (l lockCloser) Close() error {
	return l.f.Close()
}

// Lock ...
func Lock(name string) (io.Closer, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	/*
		Some people tell me FcntlFlock does not exist, so use flock here
	*/
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return nil, err
	}

	// spec := syscall.Flock_t{
	// 	Type:   syscall.F_WRLCK,
	// 	Whence: int16(os.SEEK_SET),
	// 	Start:  0,
	// 	Len:    0, // 0 means to lock the entire file.
	// 	Pid:    int32(os.Getpid()),
	// }
	// if err := syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, &spec); err != nil {
	// 	f.Close()
	// 	return nil, err
	// }

	return lockCloser{f}, nil
}
