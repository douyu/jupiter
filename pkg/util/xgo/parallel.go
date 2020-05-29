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

package xgo

import (
	"sync"

	"golang.org/x/sync/errgroup"
)

// ParallelWithError ...
func ParallelWithError(fns ...func() error) func() error {
	return func() error {
		eg := errgroup.Group{}
		for _, fn := range fns {
			eg.Go(fn)
		}

		return eg.Wait()
	}
}

// ParallelWithErrorChan calls the passed functions in a goroutine, returns a chan of errors.
// fns会并发执行，chan error
func ParallelWithErrorChan(fns ...func() error) chan error {
	total := len(fns)
	errs := make(chan error, total)

	var wg sync.WaitGroup
	wg.Add(total)

	go func(errs chan error) {
		wg.Wait()
		close(errs)
	}(errs)

	for _, fn := range fns {
		go func(fn func() error, errs chan error) {
			defer wg.Done()
			errs <- try(fn, nil)
		}(fn, errs)
	}

	return errs
}

// RestrictParallelWithErrorChan calls the passed functions in a goroutine, limiting the number of goroutines running at the same time,
// returns a chan of errors.
func RestrictParallelWithErrorChan(concurrency int, fns ...func() error) chan error {
	total := len(fns)

	if concurrency <= 0 {
		concurrency = 1
	}

	if concurrency > total {
		concurrency = total
	}

	var wg sync.WaitGroup
	wg.Add(total)

	errs := make(chan error, total)
	sem := make(chan struct{}, concurrency)
	go func(sem chan struct{}, errs chan error) {
		wg.Wait()
		close(errs)
		close(sem)
	}(sem, errs)

	for _, fn := range fns {
		go func(fn func() error, sem chan struct{}, errs chan error) {
			defer wg.Done()
			<-sem
			errs <- try(fn, nil)
			sem <- struct{}{}
		}(fn, sem, errs)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}

	return errs
}
