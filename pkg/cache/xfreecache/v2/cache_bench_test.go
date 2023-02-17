package xfreecache

import (
	"strconv"
	"testing"
)

func BenchmarkLocalCache_GetCacheData(b *testing.B) {
	localCache := New[string, Student](DefaultConfig())

	b.Run("read", func(b *testing.B) {
		student := Student{10, "student1"}
		for i := 0; i < b.N; i++ {
			_, _ = localCache.GetAndSetCacheData("mytest", student.Name, func() (Student, error) {
				res := student
				return res, nil
			})
		}
	})

	b.Run("read & write", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			student := Student{10, "student" + strconv.Itoa(i)}
			_, _ = localCache.GetAndSetCacheData("mytest", student.Name, func() (Student, error) {
				res := student
				return res, nil
			})
		}
	})
}

func BenchmarkLocalCache_GetCacheMap(b *testing.B) {

	localCache := New[int64, int64](DefaultConfig())

	b.Run("read", func(b *testing.B) {
		uidList := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		for i := 0; i < b.N; i++ {
			_, _ = localCache.GetAndSetCacheMap("mytest2", uidList, func(in []int64) (map[int64]int64, error) {
				res := make(map[int64]int64)
				for _, uid := range in {
					res[uid] = uid
				}
				return res, nil
			})
		}
	})

	b.Run("read & write", func(b *testing.B) {
		uidList := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		for i := 0; i < b.N; i++ {
			uidListT := make([]int64, 0, 10)
			for _, uid := range uidList {
				uidListT = append(uidListT, uid+int64(i))
			}
			_, _ = localCache.GetAndSetCacheMap("mytest2", uidListT, func(in []int64) (map[int64]int64, error) {
				res := make(map[int64]int64)
				for _, uid := range in {
					res[uid] = uid
				}
				return res, nil
			})
		}
	})
}
