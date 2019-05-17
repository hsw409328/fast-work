/**
 * Author: haoshuaiwei 
 * Date: 2019-05-14 18:14 
 */

package fast_drive

import (
	"fast-work/fast-sys"
	"gopkg.in/redis.v5"
	"strconv"
	"time"
)

var (
	RedisDriver *redis.Client
)

func initRedis() {
	host, _ := fast_sys.GoConfig.GetString("redis", "host")
	port, _ := fast_sys.GoConfig.GetString("redis", "port")
	authPassword, _ := fast_sys.GoConfig.GetString("redis", "auth_password")
	dbStr, _ := fast_sys.GoConfig.GetString("redis", "db")
	db, _ := strconv.Atoi(dbStr)
	RedisDriver = redis.NewClient(&redis.Options{
		Addr:        host + ":" + port,
		Password:    authPassword,
		DB:          db,
		DialTimeout: time.Second * 2,
		//IdleTimeout: time.Second * 1000000,
	})
}
