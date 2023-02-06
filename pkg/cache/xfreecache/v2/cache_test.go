package xfreecache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Student struct {
	Age  int
	Name string
}

func Test_cache_GetAndSetCacheData(t *testing.T) {
	oneCache := NewLocalCache[string, Student](DefaultConfig().
		SetExpire(60 * time.Second).SetName("local"))
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
	for _, tt := range tests {
		c := NewLocalCache[int64, int64](DefaultConfig())
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
