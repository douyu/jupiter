package xfreecache

import (
	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	"golang.org/x/sync/errgroup"
	"strconv"
	"testing"
)

// BenchmarkLocalCache_GetCacheData race检测
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

	b.Run("read & write & race", func(b *testing.B) {
		eg := errgroup.Group{}
		for i := 0; i < b.N; i++ {
			eg.Go(func() error {
				student := Student{10, "student" + strconv.Itoa(i)}
				_, _ = localCache.GetAndSetCacheData("mytest", student.Name, func() (Student, error) {
					res := student
					return res, nil
				})
				return nil
			})
		}
		_ = eg.Wait()
	})
}

// BenchmarkLocalCache_GetCacheData_Proto race检测
func BenchmarkLocalCache_GetCacheData_Proto(b *testing.B) {
	localCache := New[int, *helloworldv1.SayHiResponse](DefaultConfig())
	b.Run("read", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			student := &helloworldv1.SayHiResponse{Error: uint32(i)}
			_, _ = localCache.GetAndSetCacheData("mytest", i, func() (*helloworldv1.SayHiResponse, error) {
				res := student
				return res, nil
			})
		}
	})

	b.Run("read & write & race", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			student := &helloworldv1.SayHiResponse{Error: uint32(1)}
			data, _ := localCache.GetAndSetCacheData("mytest", 1, func() (*helloworldv1.SayHiResponse, error) {
				res := student
				return res, nil
			})
			_ = data.Data
		}
	})

	b.Run("read & write & race", func(b *testing.B) {
		eg := errgroup.Group{}
		for i := 0; i < b.N; i++ {
			eg.Go(func() error {
				student := &helloworldv1.SayHiResponse{Error: uint32(1)}
				data, _ := localCache.GetAndSetCacheData("mytest", 1, func() (*helloworldv1.SayHiResponse, error) {
					res := student
					return res, nil
				})
				_ = data.Data

				return nil
			})
		}
		_ = eg.Wait()
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
