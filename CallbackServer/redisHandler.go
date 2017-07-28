package main

import (
	"encoding/json"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/mediocregopher/radix.v2/sentinel"
	"github.com/mediocregopher/radix.v2/util"
	"strconv"
	"strings"
	"time"
)

var sentinelPool *sentinel.Client
var securitylPool *sentinel.Client
var redisPool *pool.Pool

func InitiateRedis() {

	var err error

	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		if err = client.Cmd("AUTH", redisPassword).Err; err != nil {
			client.Close()
			return nil, err
		}
		if err = client.Cmd("select", redisDb).Err; err != nil {
			client.Close()
			return nil, err
		}
		return client, nil
	}

	if redisMode == "sentinel" {
		sentinelIps := strings.Split(sentinelHosts, ",")

		if len(sentinelIps) > 1 {
			sentinelIp := fmt.Sprintf("%s:%s", sentinelIps[0], sentinelPort)
			sentinelPool, err = sentinel.NewClientCustom("tcp", sentinelIp, 10, df, redisClusterName)

			if err != nil {
				errHndlr(err)
			}

			securitylPool, err = sentinel.NewClientCustom("tcp", sentinelIp, 10, df, redisClusterName)

			if err != nil {
				errHndlr(err)
			}
		} else {
			fmt.Println("Not enough sentinel servers")
		}
	} else {
		redisPool, err = pool.NewCustom("tcp", redisIp, 10, df)

		if err != nil {
			errHndlr(err)
		}
	}
}

// Redis String Methods
func RedisAdd(key, value string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSet", r)
		}
	}()
	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	isExists, _ := client.Cmd("EXISTS", key).Int()

	if isExists == 1 {
		return "Key Already exists"
	} else {
		result, sErr := client.Cmd("set", key, value).Str()
		errHndlr(sErr)
		fmt.Println(result)
		return result
	}
}

func RedisSet(key, value string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSet", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, sErr := client.Cmd("set", key, value).Str()
	errHndlr(sErr)
	fmt.Println(result)
	return result
}

func RedisGet(key string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisGet", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	strObj, _ := client.Cmd("get", key).Str()
	fmt.Println(strObj)
	return strObj
}

func AppendIfMissing(windowList []string, i string) []string {
	for _, ele := range windowList {
		if ele == i {
			return windowList
		}
	}
	return append(windowList, i)
}

func RedisSearchKeys(pattern string) []string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSearchKeys", r)
		}
	}()
	matchingKeys := make([]string, 0)

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	fmt.Println("Start ScanAndGetKeys:: ", pattern)
	scanResult := util.NewScanner(client, util.ScanOpts{Command: "SCAN", Pattern: pattern, Count: 1000})

	for scanResult.HasNext() {
		//fmt.Println("next:", scanResult.Next())
		matchingKeys = AppendIfMissing(matchingKeys, scanResult.Next())
	}

	fmt.Println("Scan Result:: ", matchingKeys)
	return matchingKeys
}

func RedisIncr(key string) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSet", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, sErr := client.Cmd("incr", key).Int()
	errHndlr(sErr)
	fmt.Println(result)
	return result
}

func RedisIncrBy(key string, value int) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSet", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, sErr := client.Cmd("incrby", key, value).Int()
	errHndlr(sErr)
	fmt.Println(result)
	return result
}

func RedisRemove(key string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisRemove", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, sErr := client.Cmd("del", key).Int()
	errHndlr(sErr)
	fmt.Println(result)
	if result == 1 {
		return true
	} else {
		return false
	}
}

func RedisCheckKeyExist(key string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in CheckKeyExist", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, sErr := client.Cmd("exists", key).Int()
	errHndlr(sErr)
	fmt.Println(result)
	if result == 1 {
		return true
	} else {
		return false
	}
}

// Redis Hashes Methods

func RedisHashGetAll(hkey string) map[string]string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashGetAll", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}
	strHash, _ := client.Cmd("hgetall", hkey).Map()
	bytes, err := json.Marshal(strHash)
	if err != nil {
		fmt.Println(err)
	}
	text := string(bytes)
	fmt.Println(text)
	return strHash
}

func RedisHashGetField(hkey, field string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashGetAll", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}
	strValue, _ := client.Cmd("hget", hkey, field).Str()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(strValue)
	return strValue
}

func RedisHashSetField(hkey, field, value string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashSetField", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, _ := client.Cmd("hset", hkey, field, value).Int()
	if result > 0 {
		return true
	} else {
		return false
	}
}

func RedisHashSetMultipleField(hkey string, data map[string]string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashSetField", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}
	for key, value := range data {
		client.Cmd("hset", hkey, key, value)
	}
	fmt.Println(true)
	return true
}

func RedisRemoveHashField(hkey, field string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisRemoveHashField", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, _ := client.Cmd("hdel", hkey, field).Int()
	if result > 0 {
		return true
	} else {
		return false
	}
}

// Redis List Methods

func RedisListLpop(lname string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpop", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	lpopItem, _ := client.Cmd("lpop", lname).Str()
	fmt.Println(lpopItem)
	return lpopItem
}

func RedisListLpush(lname, value string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpush", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, _ := client.Cmd("lpush", lname, value).Int()
	if result > 0 {
		return true
	} else {
		return false
	}
}

func RedisListRpush(lname, value string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpush", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, _ := client.Cmd("rpush", lname, value).Int()
	if result > 0 {
		return true
	} else {
		return false
	}
}

func RedisListLlen(lname string) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLlen", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, _ := client.Cmd("llen", lname).Int()
	return result
}

// Redis Shorted List
func RedisZadd(lname, value string, score float64) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpop", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	zItem, err := client.Cmd("zadd", lname, score, value).Int()
	errHndlr(err)
	strZItem := strconv.Itoa(zItem)
	fmt.Println(strZItem)
	return strZItem, err
}

func RedisZScore(lname, value string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpop", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	zScoreStr, err := client.Cmd("zscore", lname, value).Str()
	errHndlr(err)
	fmt.Println(zScoreStr)
	return zScoreStr
}

func RedisZRangeByScore(lname, min, max string) []string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpop", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	zScoreStr, err := client.Cmd("zrangebyscore", lname, min, max).List()
	errHndlr(err)
	fmt.Println(zScoreStr)
	return zScoreStr
}

func RedisZRemove(lname, value string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisZRemove", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlr(err)
		defer redisPool.Put(client)
	}

	result, _ := client.Cmd("zrem", lname, value).Int()
	if result > 0 {
		return true
	} else {
		return false
	}
}

func SecurityGet(key string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisGet", r)
		}
	}()

	var client *redis.Client
	var err error

	if redisMode == "sentinel" {
		client, err = securitylPool.GetMaster(redisClusterName)
		client.Cmd("select", "0")
		errHndlr(err)
		defer securitylPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redis.DialTimeout("tcp", securityIp, time.Duration(10)*time.Second)
		errHndlr(err)
		defer client.Close()

		//authServer
		authE := client.Cmd("auth", redisPassword)
		errHndlr(authE.Err)
	}

	strObj, _ := client.Cmd("get", key).Str()
	//fmt.Println(strObj)
	return strObj
}
