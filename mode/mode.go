package mode

import (
	"runtime"
	"time"
	
	"github.com/gomodule/redigo/redis"
)

type IMode interface {
	GetConn() redis.Conn
	NewConn() (redis.Conn, error)
}

type PoolOption func(pool *redis.Pool)

func TestOnBorrow(value func(c redis.Conn, t time.Time) (err error)) PoolOption {
	return func(pool *redis.Pool) {
		pool.TestOnBorrow = value
	}
}

func Wait(value bool) PoolOption {
	return func(pool *redis.Pool) {
		pool.Wait = value
	}
}

func MaxIdle(value int) PoolOption {
	return func(pool *redis.Pool) {
		pool.MaxIdle = value
	}
}

func MaxActive(value int) PoolOption {
	return func(pool *redis.Pool) {
		pool.MaxActive = value
	}
}

func IdleTimeout(value time.Duration) PoolOption {
	return func(pool *redis.Pool) {
		pool.IdleTimeout = value
	}
}

func MaxConnLifetime(value time.Duration) PoolOption {
	return func(pool *redis.Pool) {
		pool.MaxConnLifetime = value
	}
}

func DefaultDialOpts() []redis.DialOption {
	return []redis.DialOption{
		redis.DialConnectTimeout(time.Millisecond * 500),
		redis.DialReadTimeout(time.Second),
		redis.DialWriteTimeout(time.Second),
	}
}

func DefaultPoolOpts() []PoolOption {
	return []PoolOption{
		Wait(false),
		MaxIdle(2 * runtime.GOMAXPROCS(0)),
	}
}
