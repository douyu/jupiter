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
func (rc *RedisClusterStub) Get(key string) string {
	var mes string
	strObj := rc.Client.Get(key)
	if err := strObj.Err(); err != nil {
		mes = ""
	} else {
		mes = strObj.Val()
	}
	return mes
}

// GetRaw ...
func (rc *RedisClusterStub) GetRaw(key string) ([]byte, error) {
	c, err := rc.Client.Get(key).Bytes()
	if err != nil && err != redis.Nil {
		return []byte{}, err
	}
	return c, nil
}

// MGet ...
func (rc *RedisClusterStub) MGet(keys ...string) ([]string, error) {
	sliceObj := rc.Client.MGet(keys...)
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
func (rc *RedisClusterStub) MGets(keys []string) ([]interface{}, error) {
	ret, err := rc.Client.MGet(keys...).Result()
	if err != nil && err != redis.Nil {
		return []interface{}{}, err
	}
	return ret, nil
}

// Set 设置redis的string
func (rc *RedisClusterStub) Set(key string, value interface{}, expire time.Duration) bool {
	err := rc.Client.Set(key, value, expire).Err()
	if err != nil {
		return false
	}
	return true
}

// HGetAll 从redis获取hash的所有键值对
func (rc *RedisClusterStub) HGetAll(key string) map[string]string {
	hashObj := rc.Client.HGetAll(key)
	hash := hashObj.Val()
	return hash
}

// HGet 从redis获取hash单个值
func (rc *RedisClusterStub) HGet(key string, fields string) (string, error) {
	strObj := rc.Client.HGet(key, fields)
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
func (rc *RedisClusterStub) HMGetMap(key string, fields []string) map[string]string {
	if len(fields) == 0 {
		return make(map[string]string)
	}
	sliceObj := rc.Client.HMGet(key, fields...)
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
func (rc *RedisClusterStub) HMSet(key string, hash map[string]interface{}, expire time.Duration) bool {
	if len(hash) > 0 {
		err := rc.Client.HMSet(key, hash).Err()
		if err != nil {
			return false
		}
		if expire > 0 {
			rc.Client.Expire(key, expire)
		}
		return true
	}
	return false
}

// HSet hset
func (rc *RedisClusterStub) HSet(key string, field string, value interface{}) bool {
	err := rc.Client.HSet(key, field, value).Err()
	if err != nil {
		return false
	}
	return true
}

// HDel ...
func (rc *RedisClusterStub) HDel(key string, field ...string) bool {
	IntObj := rc.Client.HDel(key, field...)
	if err := IntObj.Err(); err != nil {
		return false
	}

	return true
}

// SetWithErr ...
func (rc *RedisClusterStub) SetWithErr(key string, value interface{}, expire time.Duration) error {
	err := rc.Client.Set(key, value, expire).Err()
	return err
}

// SetNx 设置redis的string 如果键已存在
func (rc *RedisClusterStub) SetNx(key string, value interface{}, expiration time.Duration) bool {

	result, err := rc.Client.SetNX(key, value, expiration).Result()

	if err != nil {
		return false
	}

	return result
}

// SetNxWithErr 设置redis的string 如果键已存在
func (rc *RedisClusterStub) SetNxWithErr(key string, value interface{}, expiration time.Duration) (bool, error) {
	result, err := rc.Client.SetNX(key, value, expiration).Result()
	return result, err
}

// Incr redis自增
func (rc *RedisClusterStub) Incr(key string) bool {
	err := rc.Client.Incr(key).Err()
	if err != nil {
		return false
	}
	return true
}

// IncrWithErr ...
func (rc *RedisClusterStub) IncrWithErr(key string) (int64, error) {
	ret, err := rc.Client.Incr(key).Result()
	return ret, err
}

// IncrBy 将 key 所储存的值加上增量 increment 。
func (rc *RedisClusterStub) IncrBy(key string, increment int64) (int64, error) {
	intObj := rc.Client.IncrBy(key, increment)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// Decr redis自减
func (rc *RedisClusterStub) Decr(key string) bool {
	err := rc.Client.Decr(key).Err()
	if err != nil {
		return false
	}
	return true
}

// Type ...
func (rc *RedisClusterStub) Type(key string) (string, error) {
	statusObj := rc.Client.Type(key)
	if err := statusObj.Err(); err != nil {
		return "", err
	}

	return statusObj.Val(), nil
}

// ZRevRange 倒序获取有序集合的部分数据
func (rc *RedisClusterStub) ZRevRange(key string, start, stop int64) ([]string, error) {
	strSliceObj := rc.Client.ZRevRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// ZRevRangeWithScores ...
func (rc *RedisClusterStub) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	zSliceObj := rc.Client.ZRevRangeWithScores(key, start, stop)
	if err := zSliceObj.Err(); err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}
	return zSliceObj.Val(), nil
}

// ZRange ...
func (rc *RedisClusterStub) ZRange(key string, start, stop int64) ([]string, error) {
	strSliceObj := rc.Client.ZRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// ZRevRank ...
func (rc *RedisClusterStub) ZRevRank(key string, member string) (int64, error) {
	intObj := rc.Client.ZRevRank(key, member)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// ZRevRangeByScore ...
func (rc *RedisClusterStub) ZRevRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	res, err := rc.Client.ZRevRangeByScore(key, opt).Result()
	if err != nil && err != redis.Nil {
		return []string{}, err
	}

	return res, nil
}

// ZRevRangeByScoreWithScores ...
func (rc *RedisClusterStub) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error) {
	res, err := rc.Client.ZRevRangeByScoreWithScores(key, opt).Result()
	if err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}

	return res, nil
}

// HMGet 批量获取hash值
func (rc *RedisClusterStub) HMGet(key string, fileds []string) []string {
	sliceObj := rc.Client.HMGet(key, fileds...)
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
func (rc *RedisClusterStub) ZCard(key string) (int64, error) {
	IntObj := rc.Client.ZCard(key)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}
	return IntObj.Val(), nil
}

// ZScore 获取有序集合成员 member 的 score 值
func (rc *RedisClusterStub) ZScore(key string, member string) (float64, error) {
	FloatObj := rc.Client.ZScore(key, member)
	err := FloatObj.Err()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	return FloatObj.Val(), err
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中
func (rc *RedisClusterStub) ZAdd(key string, members ...redis.Z) (int64, error) {
	IntObj := rc.Client.ZAdd(key, members...)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (rc *RedisClusterStub) ZCount(key string, min, max string) (int64, error) {
	IntObj := rc.Client.ZCount(key, min, max)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// Del redis删除
func (rc *RedisClusterStub) Del(key string) int64 {
	result, err := rc.Client.Del(key).Result()
	if err != nil {
		return 0
	}
	return result
}

// DelWithErr ...
func (rc *RedisClusterStub) DelWithErr(key string) (int64, error) {
	result, err := rc.Client.Del(key).Result()
	return result, err
}

// HIncrBy 哈希field自增
func (rc *RedisClusterStub) HIncrBy(key string, field string, incr int) int64 {
	result, err := rc.Client.HIncrBy(key, field, int64(incr)).Result()
	if err != nil {
		return 0
	}
	return result
}

// HIncrBy 哈希field自增并且返回错误
func (rc *RedisClusterStub) HIncrByWithErr(key string, field string, incr int) (int64, error) {
	return rc.Client.HIncrBy(key, field, int64(incr)).Result()
}

// Exists 键是否存在
func (rc *RedisClusterStub) Exists(key string) bool {
	result, err := rc.Client.Exists(key).Result()
	if err != nil {
		return false
	}
	return result == 1
}

// ExistsWithErr ...
func (rc *RedisClusterStub) ExistsWithErr(key string) (bool, error) {
	result, err := rc.Client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
func (rc *RedisClusterStub) LPush(key string, values ...interface{}) (int64, error) {
	IntObj := rc.Client.LPush(key, values...)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
func (rc *RedisClusterStub) RPush(key string, values ...interface{}) (int64, error) {
	IntObj := rc.Client.RPush(key, values...)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPop 移除并返回列表 key 的尾元素。
func (rc *RedisClusterStub) RPop(key string) (string, error) {
	strObj := rc.Client.RPop(key)
	if err := strObj.Err(); err != nil {
		return "", err
	}

	return strObj.Val(), nil
}

// LRange 获取列表指定范围内的元素
func (rc *RedisClusterStub) LRange(key string, start, stop int64) ([]string, error) {
	result, err := rc.Client.LRange(key, start, stop).Result()
	if err != nil {
		return []string{}, err
	}

	return result, nil
}

// LLen ...
func (rc *RedisClusterStub) LLen(key string) int64 {
	IntObj := rc.Client.LLen(key)
	if err := IntObj.Err(); err != nil {
		return 0
	}

	return IntObj.Val()
}

// LLenWithErr ...
func (rc *RedisClusterStub) LLenWithErr(key string) (int64, error) {
	ret, err := rc.Client.LLen(key).Result()
	return ret, err
}

// LRem ...
func (rc *RedisClusterStub) LRem(key string, count int64, value interface{}) int64 {
	IntObj := rc.Client.LRem(key, count, value)
	if err := IntObj.Err(); err != nil {
		return 0
	}

	return IntObj.Val()
}

// LIndex ...
func (rc *RedisClusterStub) LIndex(key string, idx int64) (string, error) {
	ret, err := rc.Client.LIndex(key, idx).Result()
	return ret, err
}

// LTrim ...
func (rc *RedisClusterStub) LTrim(key string, start, stop int64) (string, error) {
	ret, err := rc.Client.LTrim(key, start, stop).Result()
	return ret, err
}

// ZRemRangeByRank 移除有序集合中给定的排名区间的所有成员
func (rc *RedisClusterStub) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	result, err := rc.Client.ZRemRangeByRank(key, start, stop).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

// Expire 设置过期时间
func (rc *RedisClusterStub) Expire(key string, expiration time.Duration) (bool, error) {
	result, err := rc.Client.Expire(key, expiration).Result()
	if err != nil {
		return false, err
	}

	return result, err
}

// ZRem 从zset中移除变量
func (rc *RedisClusterStub) ZRem(key string, members ...interface{}) (int64, error) {
	result, err := rc.Client.ZRem(key, members...).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// SAdd 向set中添加成员
func (rc *RedisClusterStub) SAdd(key string, member ...interface{}) (int64, error) {
	intObj := rc.Client.SAdd(key, member...)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// SMembers 返回set的全部成员
func (rc *RedisClusterStub) SMembers(key string) ([]string, error) {
	strSliceObj := rc.Client.SMembers(key)
	if err := strSliceObj.Err(); err != nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// SIsMember ...
func (rc *RedisClusterStub) SIsMember(key string, member interface{}) (bool, error) {
	boolObj := rc.Client.SIsMember(key, member)
	if err := boolObj.Err(); err != nil {
		return false, err
	}
	return boolObj.Val(), nil
}

// HKeys 获取hash的所有域
func (rc *RedisClusterStub) HKeys(key string) []string {
	strObj := rc.Client.HKeys(key)
	if err := strObj.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	return strObj.Val()
}

// HLen 获取hash的长度
func (rc *RedisClusterStub) HLen(key string) int64 {
	intObj := rc.Client.HLen(key)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0
	}
	return intObj.Val()
}

// GeoAdd 写入地理位置
func (rc *RedisClusterStub) GeoAdd(key string, location *redis.GeoLocation) (int64, error) {
	res, err := rc.Client.GeoAdd(key, location).Result()
	if err != nil {
		return 0, err
	}

	return res, nil
}

// GeoRadius 根据经纬度查询列表
func (rc *RedisClusterStub) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	res, err := rc.Client.GeoRadius(key, longitude, latitude, query).Result()
	if err != nil {
		return []redis.GeoLocation{}, err
	}

	return res, nil
}

// Ttl 查询过期时间
func (rc *RedisClusterStub) TTL(key string) (int64, error) {
	result, err := rc.Client.TTL(key).Result()
	return int64(result.Seconds()), err
}

// Close closes the cluster client, releasing any open resources.
//
// It is rare to Close a ClusterClient, as the ClusterClient is meant
// to be long-lived and shared between many goroutines.
func (rc *RedisClusterStub) Close() (err error) {
	if rc.Client != nil {
		err = rc.Client.Close()
	}
	return
}
