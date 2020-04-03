package redis

import (
	"errors"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Pool struct {
	pool *redis.Pool
}

type RedisConfig struct {
	Addr         string `json:"addr"`
	DB           int    `json:"dbnum"`
	Password     string `json:"password"`
	IdleTimeout  int    `json:"idletimeout"` //second
	ConnTimeout  int    `json:"conntimeout"` //microsecond
	ReadTimeout  int    `json:"readtimeout"` //microsecond
	WriteTimeout int    `json:"writetimeout"` //microsecond
	MaxIdle      int    `json:"maxidle"`
	MaxActive    int    `json:"maxactive"`
}

var defaultConfig = RedisConfig{
	Addr:         ":6379",
	DB:           0,
	Password:     "",
	IdleTimeout:  100,
	ConnTimeout:  200,
	ReadTimeout:  500,
	WriteTimeout: 200,
	MaxIdle:      10,
	MaxActive:    10,
}

func NewPool(config RedisConfig) (*Pool, error) {
	if config.Addr == "" {
		return nil, errors.New("addr is empty")
	}

	// default config params
	if config.ConnTimeout == 0 {
		config.ConnTimeout = defaultConfig.ConnTimeout
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = defaultConfig.WriteTimeout
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = defaultConfig.ReadTimeout
	}
	if config.MaxIdle == 0 {
		config.MaxIdle = defaultConfig.MaxIdle
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = defaultConfig.IdleTimeout
	}
	if config.MaxActive == 0 {
		config.MaxActive = config.MaxIdle
	}

	var redisPool = &redis.Pool{
		MaxIdle:     config.MaxIdle,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
		Wait:        true,
		MaxActive:   config.MaxActive,
		Dial: func() (redis.Conn, error) {
			var options []redis.DialOption
			if config.Password != "" {
				options = append(options, redis.DialPassword(config.Password))
			}
			if config.DB != 0 {
				options = append(options, redis.DialDatabase(config.DB))
			}
			// dial timeout
			options = append(options, redis.DialConnectTimeout(time.Duration(config.ConnTimeout)*time.Microsecond))
			options = append(options, redis.DialReadTimeout(time.Duration(config.ReadTimeout)*time.Microsecond))
			options = append(options, redis.DialWriteTimeout(time.Duration(config.WriteTimeout)*time.Microsecond))

			c, err := redis.Dial("tcp", config.Addr, options...)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	return &Pool{
		pool: redisPool,
	}, nil
}

func (p *Pool) getByKey(key string) (interface{}, error){
	conn := p.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print(err.Error())
		}
	}()
	reply, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (p *Pool) Set(key string, value interface{}) (bool, error){
	conn := p.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print(err.Error())
		}
	}()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return false, err
	}
	return true, nil
}