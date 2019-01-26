package alone

import (
	"github.com/gomodule/redigo/redis"
	"github.com/letsfire/redigo/mode"
)

type standAloneMode struct {
	pool *redis.Pool
}

func (sam *standAloneMode) GetConn() redis.Conn {
	return sam.pool.Get()
}

func (sam *standAloneMode) NewConn() (redis.Conn, error) {
	return sam.pool.Dial()
}

var _ mode.IMode = &standAloneMode{}

func New(optFuncs ...OptFunc) *standAloneMode {
	opts := options{
		addr:     "127.0.0.1:6379",
		dialOpts: mode.DefaultDialOpts(),
		poolOpts: mode.DefaultPoolOpts(),
	}
	for _, optFunc := range optFuncs {
		optFunc(&opts)
	}
	pool := &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", opts.addr, opts.dialOpts...)
		},
	}
	for _, poolOptFunc := range opts.poolOpts {
		poolOptFunc(pool)
	}
	return &standAloneMode{pool: pool}
}
