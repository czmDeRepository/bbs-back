package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
)

type RedisPool struct {
	*redis.Pool
}

var redisPool *RedisPool

func initRedis() {
	//连接地址
	redisConn := beego.AppConfig.DefaultString("redisConn", "127.0.0.1:6379")
	//db分区
	redisDbNum := beego.AppConfig.DefaultInt("redisDbNum", 0)
	//密码
	redisPassword := beego.AppConfig.DefaultString("redisPassword", "")

	redisPool = new(RedisPool)
	//建立连接池
	redisPool.Pool = &redis.Pool{
		//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxIdle: beego.AppConfig.DefaultInt("redisMaxIdle", 2),
		//最大的激活连接数，表示同时最多有N个连接
		MaxActive: beego.AppConfig.DefaultInt("redisMaxActive", 5),
		//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		IdleTimeout: time.Duration(beego.AppConfig.DefaultInt64("RedisIdleTimeout", int64(60*time.Second))),
		//建立连接
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConn)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			if redisPassword != "" {
				if _, authErr := c.Do("AUTH", redisPassword); authErr != nil {
					return nil, fmt.Errorf("redis auth password error: %s", authErr)
				}
			}
			//选择分区
			c.Do("SELECT", redisDbNum)
			return c, nil
		},
		//ping
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
		Wait: true, // 线程用完阻塞
	}
}

func GetRedisPool() *RedisPool {
	if redisPool == nil {
		Init()
	}
	return redisPool
}

func (pool *RedisPool) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	con := GetRedisPool().Pool.Get()
	err, ok := con.(error)
	if ok {
		panic(err)
	}
	defer con.Close()
	return logCommand(con, commandName, args...)
}

func logCommand(con redis.Conn, commandName string, args ...interface{}) (reply interface{}, err error) {
	res, err := con.Do(commandName, args...)
	if err != nil {
		logs.Error("redis command ==>", commandName, args, "result:", res, "error:", err)
	} else if beego.BConfig.RunMode != beego.PROD {
		logs.Critical("redis command ==>", commandName, args, "result:", res, "error:", err)
	}
	return res, err
}

/**
redis  SET
*/
func (pool *RedisPool) Set(key, v interface{}) (bool, error) {
	return redis.Bool(redisPool.do("SET", key, v))
}

/**
redis  GET
*/
func (pool *RedisPool) Get(key string) (string, error) {
	return redis.String(redisPool.do("GET", key))
}

/**
redis EXPIRE
*/
func (pool *RedisPool) Expire(key string, ex time.Duration) error {
	_, err := redisPool.do("EXPIRE", key, ex.Seconds())
	return err
}

func (pool *RedisPool) SetExp(key, v string, ex time.Duration) error {
	_, err := redisPool.do("SET", key, v, "EX", ex.Seconds())
	return err
}

/**
redis EXISTS
*/
func (pool *RedisPool) Exist(key string) (bool, error) {
	return redis.Bool(redisPool.do("EXISTS", key))
}

/**
redis DEL
*/
func (pool *RedisPool) Del(key string) error {
	_, err := redisPool.do("DEL", key)
	return err
}

/**
redis SETNX
*/
func (pool *RedisPool) SetNX(key string, value interface{}) error {
	_, err := redisPool.do("SETNX", key, value)
	return err
}

/**
redis GET
return map
*/
func (pool *RedisPool) GetJson(key string) (map[string]string, error) {
	var jsonData map[string]string
	bv, err := redis.Bytes(redisPool.do("GET", key))
	if err != nil {
		return nil, err
	}
	errJson := json.Unmarshal(bv, &jsonData)
	if errJson != nil {
		logs.Error("json nil", err.Error())
		return nil, err
	}
	return jsonData, nil
}

/**
redis hSet 注意 设置什么类型 取的时候需要获取对应类型
*/
func (pool *RedisPool) HSet(key string, field string, data interface{}) error {
	_, err := redisPool.do("HSET", key, field, data)
	return err
}

/**
redis hGet  设置什么类型 取的时候需要获取对应类型
*/
func (pool *RedisPool) HGet(key, field string) (interface{}, error) {
	return redisPool.do("HGET", key, field)
}

/**
redis hGetAll
return map
*/
func (pool *RedisPool) HGetAll(key string) (map[string]string, error) {
	return redis.StringMap(redisPool.do("HGETALL", key))
}

/**
redis INCR 将 key 中储存的数字值增一
*/
func (pool *RedisPool) Incr(key string) error {
	_, err := redisPool.do("INCR", key)
	return err
}

/**
redis INCRBY 将 key 所储存的值加上增量 n
*/
func (pool *RedisPool) IncrBy(key string, n int) error {
	_, err := redisPool.do("INCRBY", key, n)
	return err
}

/**
redis DECR 将 key 中储存的数字值减一。
*/
func (pool *RedisPool) Decr(key string) error {
	_, err := redisPool.do("DECR", key)
	return err
}

/**
redis DECRBY 将 key 所储存的值减去减量 n
*/
func (pool *RedisPool) DecrBy(key string, n int) error {
	_, err := redisPool.do("DECRBY", key, n)
	return err
}
