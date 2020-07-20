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

package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// Get 从redis获取string
func (r *Redis) Get(key string) string {
	var mes string
	strObj := r.Client.Get(key)
	if err := strObj.Err(); err != nil {
		mes = ""
	} else {
		mes = strObj.Val()
	}
	return mes
}

// GetRaw ...
func (r *Redis) GetRaw(key string) ([]byte, error) {
	c, err := r.Client.Get(key).Bytes()
	if err != nil && err != redis.Nil {
		return []byte{}, err
	}
	return c, nil
}

// MGet ...
func (r *Redis) MGet(keys ...string) ([]string, error) {
	sliceObj := r.Client.MGet(keys...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	tmp := sliceObj.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice, nil
}

// MGets ...
func (r *Redis) MGets(keys []string) ([]interface{}, error) {
	ret, err := r.Client.MGet(keys...).Result()
	if err != nil && err != redis.Nil {
		return []interface{}{}, err
	}
	return ret, nil
}

// Set 设置redis的string
func (r *Redis) Set(key string, value interface{}, expire time.Duration) bool {
	err := r.Client.Set(key, value, expire).Err()
	return err == nil
}

// HGetAll 从redis获取hash的所有键值对
func (r *Redis) HGetAll(key string) map[string]string {
	hashObj := r.Client.HGetAll(key)
	hash := hashObj.Val()
	return hash
}

// HGet 从redis获取hash单个值
func (r *Redis) HGet(key string, fields string) (string, error) {
	strObj := r.Client.HGet(key, fields)
	err := strObj.Err()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == redis.Nil {
		return "", nil
	}
	return strObj.Val(), nil
}

// HMGetMap 批量获取hash值，返回map
func (r *Redis) HMGetMap(key string, fields []string) map[string]string {
	if len(fields) == 0 {
		return make(map[string]string)
	}
	sliceObj := r.Client.HMGet(key, fields...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return make(map[string]string)
	}

	tmp := sliceObj.Val()
	hashRet := make(map[string]string, len(tmp))

	var tmpTagID string

	for k, v := range tmp {
		tmpTagID = fields[k]
		if v != nil {
			hashRet[tmpTagID] = v.(string)
		} else {
			hashRet[tmpTagID] = ""
		}
	}
	return hashRet
}

// HMSet 设置redis的hash
func (r *Redis) HMSet(key string, hash map[string]interface{}, expire time.Duration) bool {
	if len(hash) > 0 {
		err := r.Client.HMSet(key, hash).Err()
		if err != nil {
			return false
		}
		if expire > 0 {
			r.Client.Expire(key, expire)
		}
		return true
	}
	return false
}

// HSet hset
func (r *Redis) HSet(key string, field string, value interface{}) bool {
	err := r.Client.HSet(key, field, value).Err()
	return err == nil
}

// HDel ...
func (r *Redis) HDel(key string, field ...string) bool {
	IntObj := r.Client.HDel(key, field...)
	err := IntObj.Err()
	return err == nil
}

// SetWithErr ...
func (r *Redis) SetWithErr(key string, value interface{}, expire time.Duration) error {
	err := r.Client.Set(key, value, expire).Err()
	return err
}

// SetNx 设置redis的string 如果键已存在
func (r *Redis) SetNx(key string, value interface{}, expiration time.Duration) bool {

	result, err := r.Client.SetNX(key, value, expiration).Result()

	if err != nil {
		return false
	}

	return result
}

// SetNxWithErr 设置redis的string 如果键已存在
func (r *Redis) SetNxWithErr(key string, value interface{}, expiration time.Duration) (bool, error) {
	result, err := r.Client.SetNX(key, value, expiration).Result()
	return result, err
}

// Incr redis自增
func (r *Redis) Incr(key string) bool {
	err := r.Client.Incr(key).Err()
	return err == nil
}

// IncrWithErr ...
func (r *Redis) IncrWithErr(key string) (int64, error) {
	ret, err := r.Client.Incr(key).Result()
	return ret, err
}

// IncrBy 将 key 所储存的值加上增量 increment 。
func (r *Redis) IncrBy(key string, increment int64) (int64, error) {
	intObj := r.Client.IncrBy(key, increment)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// Decr redis自减
func (r *Redis) Decr(key string) bool {
	err := r.Client.Decr(key).Err()
	return err == nil
}

// Type ...
func (r *Redis) Type(key string) (string, error) {
	statusObj := r.Client.Type(key)
	if err := statusObj.Err(); err != nil {
		return "", err
	}

	return statusObj.Val(), nil
}

// ZRevRange 倒序获取有序集合的部分数据
func (r *Redis) ZRevRange(key string, start, stop int64) ([]string, error) {
	strSliceObj := r.Client.ZRevRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// ZRevRangeWithScores ...
func (r *Redis) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	zSliceObj := r.Client.ZRevRangeWithScores(key, start, stop)
	if err := zSliceObj.Err(); err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}
	return zSliceObj.Val(), nil
}

// ZRange ...
func (r *Redis) ZRange(key string, start, stop int64) ([]string, error) {
	strSliceObj := r.Client.ZRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// ZRevRank ...
func (r *Redis) ZRevRank(key string, member string) (int64, error) {
	intObj := r.Client.ZRevRank(key, member)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// ZRevRangeByScore ...
func (r *Redis) ZRevRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	res, err := r.Client.ZRevRangeByScore(key, opt).Result()
	if err != nil && err != redis.Nil {
		return []string{}, err
	}

	return res, nil
}

// ZRevRangeByScoreWithScores ...
func (r *Redis) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error) {
	res, err := r.Client.ZRevRangeByScoreWithScores(key, opt).Result()
	if err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}

	return res, nil
}

// HMGet 批量获取hash值
func (r *Redis) HMGet(key string, fileds []string) []string {
	sliceObj := r.Client.HMGet(key, fileds...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	tmp := sliceObj.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice
}

// ZCard 获取有序集合的基数
func (r *Redis) ZCard(key string) (int64, error) {
	IntObj := r.Client.ZCard(key)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}
	return IntObj.Val(), nil
}

// ZScore 获取有序集合成员 member 的 score 值
func (r *Redis) ZScore(key string, member string) (float64, error) {
	FloatObj := r.Client.ZScore(key, member)
	err := FloatObj.Err()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	return FloatObj.Val(), err
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中
func (r *Redis) ZAdd(key string, members ...redis.Z) (int64, error) {
	IntObj := r.Client.ZAdd(key, members...)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (r *Redis) ZCount(key string, min, max string) (int64, error) {
	IntObj := r.Client.ZCount(key, min, max)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// Del redis删除
func (r *Redis) Del(key string) int64 {
	result, err := r.Client.Del(key).Result()
	if err != nil {
		return 0
	}
	return result
}

// DelWithErr ...
func (r *Redis) DelWithErr(key string) (int64, error) {
	result, err := r.Client.Del(key).Result()
	return result, err
}

// HIncrBy 哈希field自增
func (r *Redis) HIncrBy(key string, field string, incr int) int64 {
	result, err := r.Client.HIncrBy(key, field, int64(incr)).Result()
	if err != nil {
		return 0
	}
	return result
}

// HIncrByWithErr 哈希field自增并且返回错误
func (r *Redis) HIncrByWithErr(key string, field string, incr int) (int64, error) {
	return r.Client.HIncrBy(key, field, int64(incr)).Result()
}

// Exists 键是否存在
func (r *Redis) Exists(key string) bool {
	result, err := r.Client.Exists(key).Result()
	if err != nil {
		return false
	}
	return result == 1
}

// ExistsWithErr ...
func (r *Redis) ExistsWithErr(key string) (bool, error) {
	result, err := r.Client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
func (r *Redis) LPush(key string, values ...interface{}) (int64, error) {
	IntObj := r.Client.LPush(key, values...)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
func (r *Redis) RPush(key string, values ...interface{}) (int64, error) {
	IntObj := r.Client.RPush(key, values...)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPop 移除并返回列表 key 的尾元素。
func (r *Redis) RPop(key string) (string, error) {
	strObj := r.Client.RPop(key)
	if err := strObj.Err(); err != nil {
		return "", err
	}

	return strObj.Val(), nil
}

// LRange 获取列表指定范围内的元素
func (r *Redis) LRange(key string, start, stop int64) ([]string, error) {
	result, err := r.Client.LRange(key, start, stop).Result()
	if err != nil {
		return []string{}, err
	}

	return result, nil
}

// LLen ...
func (r *Redis) LLen(key string) int64 {
	IntObj := r.Client.LLen(key)
	if err := IntObj.Err(); err != nil {
		return 0
	}

	return IntObj.Val()
}

// LLenWithErr ...
func (r *Redis) LLenWithErr(key string) (int64, error) {
	ret, err := r.Client.LLen(key).Result()
	return ret, err
}

// LRem ...
func (r *Redis) LRem(key string, count int64, value interface{}) int64 {
	IntObj := r.Client.LRem(key, count, value)
	if err := IntObj.Err(); err != nil {
		return 0
	}

	return IntObj.Val()
}

// LIndex ...
func (r *Redis) LIndex(key string, idx int64) (string, error) {
	ret, err := r.Client.LIndex(key, idx).Result()
	return ret, err
}

// LTrim ...
func (r *Redis) LTrim(key string, start, stop int64) (string, error) {
	ret, err := r.Client.LTrim(key, start, stop).Result()
	return ret, err
}

// ZRemRangeByRank 移除有序集合中给定的排名区间的所有成员
func (r *Redis) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	result, err := r.Client.ZRemRangeByRank(key, start, stop).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

// Expire 设置过期时间
func (r *Redis) Expire(key string, expiration time.Duration) (bool, error) {
	result, err := r.Client.Expire(key, expiration).Result()
	if err != nil {
		return false, err
	}

	return result, err
}

// ZRem 从zset中移除变量
func (r *Redis) ZRem(key string, members ...interface{}) (int64, error) {
	result, err := r.Client.ZRem(key, members...).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// SAdd 向set中添加成员
func (r *Redis) SAdd(key string, member ...interface{}) (int64, error) {
	intObj := r.Client.SAdd(key, member...)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// SMembers 返回set的全部成员
func (r *Redis) SMembers(key string) ([]string, error) {
	strSliceObj := r.Client.SMembers(key)
	if err := strSliceObj.Err(); err != nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// SIsMember ...
func (r *Redis) SIsMember(key string, member interface{}) (bool, error) {
	boolObj := r.Client.SIsMember(key, member)
	if err := boolObj.Err(); err != nil {
		return false, err
	}
	return boolObj.Val(), nil
}

// HKeys 获取hash的所有域
func (r *Redis) HKeys(key string) []string {
	strObj := r.Client.HKeys(key)
	if err := strObj.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	return strObj.Val()
}

// HLen 获取hash的长度
func (r *Redis) HLen(key string) int64 {
	intObj := r.Client.HLen(key)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0
	}
	return intObj.Val()
}

// GeoAdd 写入地理位置
func (r *Redis) GeoAdd(key string, location *redis.GeoLocation) (int64, error) {
	res, err := r.Client.GeoAdd(key, location).Result()
	if err != nil {
		return 0, err
	}

	return res, nil
}

// GeoRadius 根据经纬度查询列表
func (r *Redis) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	res, err := r.Client.GeoRadius(key, longitude, latitude, query).Result()
	if err != nil {
		return []redis.GeoLocation{}, err
	}

	return res, nil
}

// TTL 查询过期时间
func (r *Redis) TTL(key string) (int64, error) {
	if result, err := r.Client.TTL(key).Result(); err != nil {
		return 0, err
	} else {
		return int64(result.Seconds()), nil
	}
}

// Close closes the cluster client, releasing any open resources.
//
// It is rare to Close a ClusterClient, as the ClusterClient is meant
// to be long-lived and shared between many goroutines.
func (r *Redis) Close() (err error) {
	err = nil
	if r.Client != nil {
		if r.Cluster() != nil {
			err = r.Cluster().Close()
		}

		if r.Stub() != nil {
			err = r.Stub().Close()
		}
	}
	return err
}
