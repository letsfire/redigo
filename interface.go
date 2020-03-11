package redigo

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gomodule/redigo/redis"
)

type ModeInterface interface {
	fmt.Stringer
	GetConn() redis.Conn
	NewConn() (redis.Conn, error)
}

// DefaultDialOpts 默认连接配置
func DefaultDialOpts() []redis.DialOption {
	return []redis.DialOption{
		redis.DialConnectTimeout(time.Second),
		redis.DialReadTimeout(time.Second * 3),
		redis.DialWriteTimeout(time.Second * 3),
	}
}

// DefaultPoolOpts 默认连接池配置
func DefaultPoolOpts() []PoolOption {
	return []PoolOption{
		Wait(false),
		MaxIdle(2 * runtime.GOMAXPROCS(0)),
		IdleTimeout(time.Second * 15),
	}
}
