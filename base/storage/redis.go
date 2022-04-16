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
		logs.Error("redisPool get connect fail： %s", err.Error())
	}
	defer con.Close()
	return logCommand(con, commandName, args...)
}

func logCommand(con redis.Conn, commandName string, args ...interface{}) (reply interface{}, err error) {
	res, err := con.Do(commandName, args...)
	if err != nil {
		logs.Error("redis command ==>", commandName, args, "result:", res, "error:", err)
	} else if beego.BConfig.RunMode != beego.PROD {
		logs.Info("redis command ==>", commandName, args, "result:", res, "error:", err)
	}
	return res, err
}

/**
redis  SET
*/
func (pool *RedisPool) Set(key, v interface{}, ex ...time.Duration) (err error) {
	if len(ex) > 0 {
		_, err := pool.do("SET", key, v, "EX", ex[0].Seconds())
		return err
	}
	_, err = pool.do("SET", key, v)
	return
}

/**
redis  GET
*/
func (pool *RedisPool) Get(key string) (string, error) {
	return redis.String(pool.do("GET", key))
}

/**
redis  GET
*/
func (pool *RedisPool) GetInt64(key string) (int64, error) {
	res, err := pool.do("GET", key)
	if err != nil || res == nil {
		return 0, err
	}
	return redis.Int64(res, err)
}

/**
redis EXPIRE
*/
func (pool *RedisPool) Expire(key string, ex time.Duration) error {
	_, err := pool.do("EXPIRE", key, ex.Seconds())
	return err
}

func (pool *RedisPool) SetExp(key, v string, ex time.Duration) error {
	_, err := pool.do("SET", key, v, "EX", ex.Seconds())
	return err
}

// GetTtl 剩余过期时间
func (pool *RedisPool) GetTtl(key string) (int64, error) {
	return redis.Int64(pool.do("TTL", key))
}

/**
redis EXISTS
*/
func (pool *RedisPool) Exist(key string) (bool, error) {
	return redis.Bool(pool.do("EXISTS", key))
}

/**
redis DEL
*/
func (pool *RedisPool) Del(key string) error {
	_, err := pool.do("DEL", key)
	return err
}

/**
redis SETNX
*/
func (pool *RedisPool) SetNX(key string, value interface{}) error {
	_, err := pool.do("SETNX", key, value)
	return err
}

/**
redis GET
return map
*/
func (pool *RedisPool) SetJson(key string, value interface{}, ex ...time.Duration) error {
	valueStr, errJson := json.Marshal(value)
	if errJson != nil {
		logs.Error("json nil", errJson.Error())
		return errJson
	}
	return pool.Set(key, valueStr, ex...)
}

/**
redis GET
*/
func (pool *RedisPool) GetJson(key string, res interface{}) error {
	bv, err := redis.Bytes(pool.do("GET", key))
	if err != nil {
		return err
	}
	errJson := json.Unmarshal(bv, res)
	if errJson != nil {
		logs.Error("json nil", errJson.Error())
		return err
	}
	return nil
}

/**
redis hSet 注意 设置什么类型 取的时候需要获取对应类型
*/
func (pool *RedisPool) HSet(key string, field string, data interface{}) error {
	_, err := pool.do("HSET", key, field, data)
	return err
}

/**
redis hGet  设置什么类型 取的时候需要获取对应类型
*/
func (pool *RedisPool) HGet(key, field string) (interface{}, error) {
	return pool.do("HGET", key, field)
}

/**
redis hGetAll
return map
*/
func (pool *RedisPool) HGetAll(key string) (map[string]string, error) {
	return redis.StringMap(pool.do("HGETALL", key))
}

/**
redis INCR 将 key 中储存的数字值增一
*/
func (pool *RedisPool) Incr(key string) error {
	_, err := pool.do("INCR", key)
	return err
}

/**
redis INCRBY 将 key 所储存的值加上增量 n
*/
func (pool *RedisPool) IncrBy(key string, n int) error {
	_, err := pool.do("INCRBY", key, n)
	return err
}

/**
redis DECR 将 key 中储存的数字值减一。
*/
func (pool *RedisPool) Decr(key string) error {
	_, err := pool.do("DECR", key)
	return err
}

/**
redis DECRBY 将 key 所储存的值减去减量 n
*/
func (pool *RedisPool) DecrBy(key string, n int) error {
	_, err := pool.do("DECRBY", key, n)
	return err
}

/**
Hyperloglog 基数统计
*/
func (pool *RedisPool) PFADD(key string, value interface{}) (err error) {
	_, err = pool.do("PFADD", key, value)
	return
}

/**
Hyperloglog 基数统计
*/
func (pool *RedisPool) PFCOUNT(key string) (int64, error) {
	return redis.Int64(pool.do("PFCOUNT", key))
}

/**
SetBit 将bitmap中index偏移量值置为1（value == nil or value[0] == true）,0(value[0]==false)
*/
func (pool *RedisPool) SetBit(key string, index int64, value ...bool) (err error) {
	val := 1
	if len(value) > 0 && !value[0] {
		val = 0
	}
	_, err = pool.do("SETBIT", key, index, val)
	return
}

/**
GetBit 获取bitmap偏移量为index的值
*/
func (pool *RedisPool) GetBit(key string, index int) (bool, error) {
	return redis.Bool(pool.do("GETBIT", key, index))
}

/**
BitCount 统计bitmap
*/
func (pool *RedisPool) BitCount(key string) (int64, error) {
	return redis.Int64(pool.do("BITCOUNT", key))
}
