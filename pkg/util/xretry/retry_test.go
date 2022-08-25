package xretry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	tests := []struct {
		name       string
		retray     int
		retrayTime time.Duration
		expect     int
	}{
		{
			"retry-3-100ms",
			3,
			time.Millisecond * 200,
			2,
		},
		{
			"retry-5-200ms",
			5,
			time.Millisecond * 100,
			0,
		},
		{
			"retry-5-200ms",
			0,
			time.Millisecond * 100,
			1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start := time.Now()
			cnt := 0
			Do(test.retray, test.retrayTime, func() error {
				cnt++
				if test.expect > 0 && cnt >= test.expect {
					return nil
				}
				return errors.New("")
			})
			cost := time.Since(start).Milliseconds()
			if test.expect > 0 {
				assert.EqualValues(t, cnt, test.expect)
			} else {
				assert.EqualValues(t, cnt, test.retray+1)

			}

			if cost < int64(test.expect-1)*test.retrayTime.Milliseconds() {
				t.Fail()
			}
		})

	}

}
