package xgolanglru

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试相关值为空的测试用例

func Test_cache_proto_Null_GetAndSetCacheMap(t *testing.T) {
	var configStr = `
		[jupiter.xgolanglru]
			[jupiter.xgolanglru.test1]
				expire = "60s"
			[jupiter.xgolanglru.test2]
				expire = "10s"
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))

	type args struct {
		ids []int64
	}
	// 初始化测试数据和测试用例
	data := map[int64]*helloworldv1.SayHiRequest{2: {Name: "2"}, 3: {Name: "3"}, 4: {Name: "4"}, 5: {Name: "5"}, 6: {Name: "6"}}
	tests := []struct {
		args  args
		wantV map[int64]*helloworldv1.SayHiRequest
	}{
		{
			args: args{
				ids: []int64{1, 2, 1, 3},
			},
			wantV: map[int64]*helloworldv1.SayHiRequest{2: data[2], 3: data[3]},
		},
		{
			args: args{
				ids: []int64{1, 2, 3, 4},
			},
			wantV: map[int64]*helloworldv1.SayHiRequest{2: data[2], 3: data[3], 4: data[4]},
		},
		{
			args: args{
				ids: []int64{1, 5, 6},
			},
			wantV: map[int64]*helloworldv1.SayHiRequest{5: data[5], 6: data[6]},
		},
		{
			args: args{
				ids: []int64{1, 2, 3},
			},
			wantV: map[int64]*helloworldv1.SayHiRequest{2: data[2], 3: data[3]},
		},
	}

	for i := 1; i <= 2; i++ {
		fmt.Printf("\n======== %d =========\n", i)
		missCount := 0
		c := StdNew[int64, *helloworldv1.SayHiRequest](fmt.Sprintf("test%d", i))
		for _, tt := range tests {
			gotV, err := c.GetAndSetCacheMap("Test_cache_proto_Null_GetAndSetCacheMap", tt.args.ids, func(in []int64) (map[int64]*helloworldv1.SayHiRequest, error) {
				missCount++
				res := make(map[int64]*helloworldv1.SayHiRequest)
				for _, uid := range in {
					if val, ok := data[uid]; ok {
						res[uid] = val
					}
				}
				fmt.Println("======== in =========")
				fmt.Println(res)
				return res, nil
			})
			fmt.Println("======== out =========")
			fmt.Println(gotV)
			assert.Nil(t, err, fmt.Sprintf("GetAndSetCacheMap(%v)", tt.args.ids))
			assert.Equalf(t, len(gotV), len(tt.wantV), "GetAndSetCacheMap(%v) len", tt.args.ids)
			for k, v := range gotV {
				val, ok := tt.wantV[k]
				assert.Equalf(t, ok, true, "GetAndSetCacheMap(%v) ok", tt.args.ids)
				assert.Equalf(t, v.GetName(), val.GetName(), "GetAndSetCacheMap(%v) val", tt.args.ids)
			}
		}
		assert.Equalf(t, missCount, 3, "GetAndSetCacheMap miss count error")
	}
}

func Test_cache_json_Null_GetAndSetCacheMap(t *testing.T) {
	var configStr = `
		[jupiter.xgolanglru]
			size = "64m"
			sizeLru = 2000
			[jupiter.xgolanglru.test1]
				expire = "60s"
			[jupiter.xgolanglru.test2]
				expire = "10s"
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))

	type args struct {
		ids []int64
	}
	// 初始化测试数据和测试用例
	data := map[int64]*Student{2: {Name: "2"}, 3: {Name: "3"}, 4: {Name: "4"}, 5: {Name: "5"}, 6: {Name: "6"}}
	tests := []struct {
		args  args
		wantV map[int64]*Student
	}{
		{
			args: args{
				ids: []int64{1, 2, 1, 3},
			},
			wantV: map[int64]*Student{2: data[2], 3: data[3]},
		},
		{
			args: args{
				ids: []int64{1, 2, 3, 4},
			},
			wantV: map[int64]*Student{2: data[2], 3: data[3], 4: data[4]},
		},
		{
			args: args{
				ids: []int64{1, 5, 6},
			},
			wantV: map[int64]*Student{5: data[5], 6: data[6]},
		},
		{
			args: args{
				ids: []int64{1, 2, 3},
			},
			wantV: map[int64]*Student{2: data[2], 3: data[3]},
		},
	}

	for i := 1; i <= 2; i++ {
		fmt.Printf("\n======== %d =========\n", i)
		missCount := 0
		c := StdNew[int64, *Student](fmt.Sprintf("test%d", i))
		for _, tt := range tests {
			gotV, err := c.GetAndSetCacheMap("Test_cache_json_Null_GetAndSetCacheMap", tt.args.ids, func(in []int64) (map[int64]*Student, error) {
				missCount++
				res := make(map[int64]*Student)
				for _, uid := range in {
					if val, ok := data[uid]; ok {
						res[uid] = val
					}
				}
				fmt.Println("======== in =========")
				fmt.Println(res)
				return res, nil
			})
			fmt.Println("======== out =========")
			fmt.Println(gotV)
			assert.Nil(t, err, fmt.Sprintf("GetAndSetCacheMap(%v)", tt.args.ids))
			assert.Equalf(t, len(gotV), len(tt.wantV), "GetAndSetCacheMap(%v) len", tt.args.ids)
			for k, v := range gotV {
				val, ok := tt.wantV[k]
				assert.Equalf(t, ok, true, "GetAndSetCacheMap(%v) ok", tt.args.ids)
				assert.Equalf(t, v.Name, val.Name, "GetAndSetCacheMap(%v) val", tt.args.ids)
			}
		}
		assert.Equalf(t, missCount, 3, "GetAndSetCacheMap miss count error")
	}
}

func Test_cache_struct_Null_GetAndSetCacheMap(t *testing.T) {
	var configStr = `
		[jupiter.xgolanglru]
			size = "64m"
			sizeLru = 2000
			[jupiter.xgolanglru.test1]
				expire = "60s"
			[jupiter.xgolanglru.test2]
				expire = "10s"
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))

	type args struct {
		ids []int64
	}
	// 初始化测试数据和测试用例
	data := map[int64]Student{2: {Name: "2"}, 3: {Name: "3"}, 4: {Name: "4"}, 5: {Name: "5"}, 6: {Name: "6"}}
	tests := []struct {
		args  args
		wantV map[int64]Student
	}{
		{
			args: args{
				ids: []int64{1, 2, 1, 3},
			},
			wantV: map[int64]Student{2: data[2], 3: data[3]},
		},
		{
			args: args{
				ids: []int64{1, 2, 3, 4},
			},
			wantV: map[int64]Student{2: data[2], 3: data[3], 4: data[4]},
		},
		{
			args: args{
				ids: []int64{1, 5, 6},
			},
			wantV: map[int64]Student{5: data[5], 6: data[6]},
		},
		{
			args: args{
				ids: []int64{1, 2, 3},
			},
			wantV: map[int64]Student{2: data[2], 3: data[3]},
		},
	}

	for i := 1; i <= 2; i++ {
		fmt.Printf("\n======== %d =========\n", i)
		missCount := 0
		c := StdNew[int64, Student](fmt.Sprintf("test%d", i))
		for _, tt := range tests {
			gotV, err := c.GetAndSetCacheMap("Test_cache_struct_Null_GetAndSetCacheMap", tt.args.ids, func(in []int64) (map[int64]Student, error) {
				missCount++
				res := make(map[int64]Student)
				for _, uid := range in {
					if val, ok := data[uid]; ok {
						res[uid] = val
					}
				}
				fmt.Println("======== in =========")
				fmt.Println(res)
				return res, nil
			})
			fmt.Println("======== out =========")
			fmt.Println(gotV)
			assert.Nil(t, err, fmt.Sprintf("GetAndSetCacheMap(%v)", tt.args.ids))
			assert.Equalf(t, len(gotV), len(tt.wantV), "GetAndSetCacheMap(%v) len", tt.args.ids)
			for k, v := range gotV {
				val, ok := tt.wantV[k]
				assert.Equalf(t, ok, true, "GetAndSetCacheMap(%v) ok", tt.args.ids)
				assert.Equalf(t, v.Name, val.Name, "GetAndSetCacheMap(%v) val", tt.args.ids)
			}
		}
		assert.Equalf(t, missCount, 3, "GetAndSetCacheMap miss count error")
	}
}
