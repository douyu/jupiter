package xgolanglru

import (
	"bytes"
	"errors"
	"fmt"
	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/stretchr/testify/assert"
)

type Student struct {
	Age  int
	Name string
}

func Test_cache_GetAndSetCacheData(t *testing.T) {
	var configStr = `
		[jupiter.xgolanglru]
			[jupiter.xgolanglru.test1]
				expire = "60s"
			[jupiter.xgolanglru.test2]
				expire = "10s"
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))

	for i := 1; i <= 2; i++ {
		oneCache := StdNew[string, Student](fmt.Sprintf("test%d", i))
		missCount := 0

		tests := []struct {
			stu Student
		}{
			{
				stu: Student{
					Age:  1,
					Name: "Student 1",
				},
			},
			{
				stu: Student{
					Age:  2,
					Name: "Student 2",
				},
			},
			{
				stu: Student{
					Age:  1,
					Name: "Student 1",
				},
			},
			{
				stu: Student{
					Age:  1,
					Name: "Student 3",
				},
			},
			{
				stu: Student{
					Age:  2,
					Name: "Student 2",
				},
			},
		}

		for _, tt := range tests {
			key := fmt.Sprintf("%d-%s", tt.stu.Age, tt.stu.Name)
			result, _ := oneCache.GetAndSetCacheData(key, tt.stu.Name, func() (Student, error) {
				missCount++
				fmt.Println("local cache miss hit")
				return tt.stu, nil
			})
			fmt.Println(result)
			assert.Equalf(t, tt.stu, result, "GetAndSetCacheData(%v) cache value error", key)
		}
		assert.Equalf(t, missCount, 3, "GetAndSetCacheData miss count error")

		missCount = 0
		for _, tt := range tests {
			key := fmt.Sprintf("%d-%s-new", tt.stu.Age, tt.stu.Name)

			result := oneCache.GetCacheValue(key, tt.stu.Name)
			fmt.Println(result)
			if result == (Student{}) {
				missCount++
				fmt.Println("local cache miss hit")
			} else {
				assert.Equalf(t, tt.stu, result, "split GetAndSetCacheData(%v) cache value error", key)
			}

			err := oneCache.SetCacheValue(key, tt.stu.Name, func() (Student, error) {
				return tt.stu, nil
			})
			assert.Emptyf(t, err, "split GetAndSetCacheData err not nil")
		}
		assert.Equalf(t, missCount, 3, "split GetAndSetCacheData miss count error")
	}
}

func Test_cache_GetAndSetCacheData_proto(t *testing.T) {
	var configStr = `
		[jupiter.xgolanglru]
			[jupiter.xgolanglru.test1]
				expire = "60s"
			[jupiter.xgolanglru.test2]
				expire = "10s"
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))

	for i := 1; i <= 2; i++ {
		oneCache := StdNew[string, *helloworldv1.SayHiRequest](fmt.Sprintf("test%d", i))
		missCount := 0

		tests := []struct {
			stu *helloworldv1.SayHiRequest
		}{
			{
				stu: &helloworldv1.SayHiRequest{
					Name: "Student 1",
				},
			},
			{
				stu: &helloworldv1.SayHiRequest{
					Name: "Student 2",
				},
			},
			{
				stu: &helloworldv1.SayHiRequest{
					Name: "Student 1",
				},
			},
			{
				stu: &helloworldv1.SayHiRequest{
					Name: "Student 3",
				},
			},
			{
				stu: &helloworldv1.SayHiRequest{
					Name: "Student 2",
				},
			},
			{
				stu: nil,
			},
			{
				stu: nil,
			},
		}

		for _, tt := range tests {
			key := tt.stu.GetName()
			result, _ := oneCache.GetAndSetCacheData(key, tt.stu.GetName(), func() (*helloworldv1.SayHiRequest, error) {
				missCount++
				fmt.Println("local cache miss hit")
				return tt.stu, nil
			})
			fmt.Println(result)
			assert.Equalf(t, tt.stu.GetName(), result.GetName(), "GetAndSetCacheData(%v) cache value error", key)
		}
		assert.Equalf(t, missCount, 4, "GetAndSetCacheData miss count error")

		missCount = 0
		for _, tt := range tests {
			key := tt.stu.GetName() + "-new"

			result := oneCache.GetCacheValue(key, tt.stu.GetName())
			fmt.Println(result)
			if result == nil {
				missCount++
				fmt.Println("local cache miss hit")
			} else {
				assert.Equalf(t, tt.stu.GetName(), result.GetName(), "split GetAndSetCacheData(%v) cache value error", key)
			}

			err := oneCache.SetCacheValue(key, tt.stu.GetName(), func() (*helloworldv1.SayHiRequest, error) {
				return tt.stu, nil
			})
			assert.Emptyf(t, err, "split GetAndSetCacheData err not nil")
		}
		assert.Equalf(t, missCount, 5, "split GetAndSetCacheData miss count error")
	}
}

func Test_cache_GetAndSetCacheMap(t *testing.T) {
	type args struct {
		ids []int64
	}
	tests := []struct {
		args  args
		wantV map[int64]int64
	}{
		{
			args: args{
				ids: []int64{1, 2, 1, 3},
			},
			wantV: map[int64]int64{1: 1, 2: 2, 3: 3},
		},
		{
			args: args{
				ids: []int64{2, 3, 4},
			},
			wantV: map[int64]int64{2: 2, 3: 3, 4: 4},
		},
		{
			args: args{
				ids: []int64{9, 6},
			},
			wantV: map[int64]int64{9: 9, 6: 6},
		},
		{
			args: args{
				ids: []int64{1, 2, 3},
			},
			wantV: map[int64]int64{1: 1, 2: 2, 3: 3},
		},
	}

	missCount := 0
	c := New[int64, int64](DefaultConfig())
	for _, tt := range tests {
		gotV, err := c.GetAndSetCacheMap("mytest2", tt.args.ids, func(in []int64) (map[int64]int64, error) {
			missCount++
			res := make(map[int64]int64)
			for _, uid := range in {
				res[uid] = uid
			}
			fmt.Println("======== in =========")
			fmt.Println(res)
			return res, nil
		})
		fmt.Println("======== out =========")
		fmt.Println(gotV)
		assert.Nil(t, err, fmt.Sprintf("GetAndSetCacheMap(%v)", tt.args.ids))
		assert.Equalf(t, tt.wantV, gotV, "GetAndSetCacheMap(%v)", tt.args.ids)
	}
	assert.Equalf(t, missCount, 3, "GetAndSetCacheMap miss count error")
}

func TestStdConfig(t *testing.T) {
	var configStr = `
		[jupiter.xgolanglru]
			[jupiter.xgolanglru.test]
				expire = "1m"
				disableMetric = true
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	t.Run("std config on addr nil", func(t *testing.T) {
		var config *Config
		name := "test"
		config = StdConfig(name)
		assert.NotNil(t, config)
		assert.Equalf(t, config.Name, name, "StdConfig Name")
		assert.Equalf(t, config.Expire, 1*time.Minute, "StdConfig Expire")
		assert.Equalf(t, config.DisableMetric, true, "StdConfig DisableMetric")
	})

}

func TestGetAndSetDataWithError(t *testing.T) {
	var configStr = `
		[jupiter.cache]
			size = "128m"
			[jupiter.cache.test1]
				expire = "60s"
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))

	for i := 1; i <= 1; i++ {
		oneCache := StdNew[string, *helloworldv1.SayHiRequest](fmt.Sprintf("test%d", i))
		missCount := 0

		tests := []struct {
			stu *helloworldv1.SayHiRequest
		}{
			{
				stu: &helloworldv1.SayHiRequest{
					Name: "Student 1",
				},
			},
			{
				stu: &helloworldv1.SayHiRequest{
					Name: "Student 2",
				},
			},
			{
				stu: &helloworldv1.SayHiRequest{
					Name: "Student 1",
				},
			},
			{
				stu: &helloworldv1.SayHiRequest{
					Name: "Student 3",
				},
			},
			{
				stu: &helloworldv1.SayHiRequest{
					Name: "Student 2",
				},
			},
			{
				stu: nil,
			},
			{
				stu: nil,
			},
		}

		for _, tt := range tests {
			key := tt.stu.GetName()
			result, _ := oneCache.GetAndSetCacheData(key, tt.stu.GetName(), func() (*helloworldv1.SayHiRequest, error) {
				missCount++
				fmt.Println("local cache miss hit")
				return tt.stu, errors.New("this is a test error")
			})
			fmt.Println(result)
			assert.Equalf(t, tt.stu.GetName(), result.GetName(), "GetAndSetCacheData(%v) cache value error", key)
		}
		// 因为接口报错了，所以全部没有命中缓存
		assert.Equalf(t, missCount, 7, "GetAndSetCacheData miss count error")
	}
}
