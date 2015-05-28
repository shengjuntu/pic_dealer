// insert_db.go
package main

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

var (
	RedisClients *redis.Pool
	REDIS_HOST = "127.0.0.1:6379"
	REDIS_DB = 0
	MAX_IDLE = 1
	MAX_ACTIVE = 10
	IDLE_TIMEOUT = 180*time.Second
	CUR_IMG_ID = "cur_img_id"
)


func init_pool() {
	RedisClients = &redis.Pool {
		MaxIdle : MAX_IDLE,
		MaxActive : MAX_ACTIVE,
		IdleTimeout : IDLE_TIMEOUT,
		Dial : func()(redis.Conn, error) {
			c,err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}
}

func db_get_image_id()(int64, error) {
	conn := RedisClients.Get()
	defer conn.Close()
	
	reply, err := conn.Do("GET", CUR_IMG_ID)
	str_id, err := redis.String(reply, err)
	if err == nil {
		return strconv.ParseInt(str_id, 10, 64)
	}
    return -1, err
}

func db_get_image(id string)(_image [] byte, _e error) {
	conn := RedisClients.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", id)
	image,err := redis.Bytes(reply, err)
	return image,err
}

func db_list_image()([] string, error) {
	conn := RedisClients.Get()
	defer conn.Close()
	reply, err := conn.Do("KEYS", "raw_img_*")
	return redis.Strings(reply, err)
}

func db_insert_image(image [] byte, id string)(int64, error) {
	conn := RedisClients.Get()
	defer conn.Close()
	
	conn.Send("MULTI")
	conn.Send("SET", id, image)
	conn.Send("INCR", CUR_IMG_ID)
	reply, err := conn.Do("EXEC")
	
	return redis.Int64(reply, err)
}
