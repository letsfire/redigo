package redigo

import (
	"github.com/gomodule/redigo/redis"
	"github.com/letsfire/redigo/mode"
)

type SubFunc func(c redis.PubSubConn) (err error)
type ExecFunc func(c redis.Conn) (res interface{}, err error)

type Redigo struct{ mode mode.IMode }

func New(mode mode.IMode) *Redigo {
	return &Redigo{mode: mode}
}

func (r *Redigo) Mode() string {
	return r.mode.String()
}

func (r *Redigo) Sub(fn SubFunc) (err error) {
	conn, err := r.mode.NewConn()
	if err != nil {
		return
	}
	psConn := redis.PubSubConn{Conn: conn}
	err = fn(psConn)
	psConn.Close()
	return
}

func (r *Redigo) Exec(fn ExecFunc) (res interface{}, err error) {
	conn := r.mode.GetConn()
	res, err = fn(conn)
	conn.Close()
	if err != nil {
		if _, ok := err.(redis.Error); ok {
			return
		} else if nconn, nerr := r.mode.NewConn(); nerr != nil {
			return
		} else {
			res, err = fn(nconn)
			nconn.Close()
		}
	}
	return
}

func (r *Redigo) Int(fn ExecFunc) (int, error) {
	return redis.Int(r.Exec(fn))
}

func (r *Redigo) Ints(fn ExecFunc) ([]int, error) {
	return redis.Ints(r.Exec(fn))
}

func (r *Redigo) IntMap(fn ExecFunc) (map[string]int, error) {
	return redis.IntMap(r.Exec(fn))
}

func (r *Redigo) Int64(fn ExecFunc) (int64, error) {
	return redis.Int64(r.Exec(fn))
}

func (r *Redigo) Int64s(fn ExecFunc) ([]int64, error) {
	return redis.Int64s(r.Exec(fn))
}
func (r *Redigo) Int64Map(fn ExecFunc) (map[string]int64, error) {
	return redis.Int64Map(r.Exec(fn))
}

func (r *Redigo) Uint64(fn ExecFunc) (uint64, error) {
	return redis.Uint64(r.Exec(fn))
}

func (r *Redigo) Bool(fn ExecFunc) (bool, error) {
	return redis.Bool(r.Exec(fn))
}

func (r *Redigo) String(fn ExecFunc) (string, error) {
	return redis.String(r.Exec(fn))
}

func (r *Redigo) StringMap(fn ExecFunc) (map[string]string, error) {
	return redis.StringMap(r.Exec(fn))
}

func (r *Redigo) Strings(fn ExecFunc) ([]string, error) {
	return redis.Strings(r.Exec(fn))
}

func (r *Redigo) Bytes(fn ExecFunc) ([]byte, error) {
	return redis.Bytes(r.Exec(fn))
}

func (r *Redigo) ByteSlices(fn ExecFunc) ([][]byte, error) {
	return redis.ByteSlices(r.Exec(fn))
}

func (r *Redigo) Positions(fn ExecFunc) ([]*[2]float64, error) {
	return redis.Positions(r.Exec(fn))
}

func (r *Redigo) Float64(fn ExecFunc) (float64, error) {
	return redis.Float64(r.Exec(fn))
}

func (r *Redigo) Float64s(fn ExecFunc) ([]float64, error) {
	return redis.Float64s(r.Exec(fn))
}

func (r *Redigo) Values(fn ExecFunc) ([]interface{}, error) {
	return redis.Values(r.Exec(fn))
}
