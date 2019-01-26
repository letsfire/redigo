package redigo

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/letsfire/redigo/mode/alone"
	"github.com/letsfire/redigo/mode/cluster"
	"github.com/letsfire/redigo/mode/sentinel"
)

func BenchmarkAloneMode_Exec(b *testing.B) {
	echoStr := "hello world"
	aRegido := New(alone.New(
		alone.Addr("192.168.0.110:6379"),
	))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res, err := aRegido.String(func(c redis.Conn) (res interface{}, err error) {
				return c.Do("ECHO", echoStr)
			})
			if err != nil {
				b.Errorf("exec failed, err = %s", err)
			} else if res != echoStr {
				b.Errorf("unexpected result, expect = %s, but = %s", echoStr, res)
			}
		}
	})
}

func BenchmarkAloneMode_Sub(b *testing.B) {
	message := "hello world"
	channel := "test-channel"
	aRegido := New(alone.New(
		alone.Addr("192.168.0.110:6379"),
	))
	var counter int32
	notifyChan := make(chan struct{})
	go aRegido.Sub(func(c redis.PubSubConn) (err error) {
		if err = c.Subscribe(channel); err != nil {
			b.Error(err)
		}
		for {
			switch msg := c.ReceiveWithTimeout(0).(type) {
			case redis.Subscription:
				notifyChan <- struct{}{}
			case redis.Message:
				atomic.AddInt32(&counter, -1)
				if string(msg.Data) != message {
					b.Errorf("unexpected result, expect = %s, but = %s", message, msg.Data)
				}
			case error:
				b.Errorf("receive failed, err = %s", msg)
			}
		}
	})
	<-notifyChan
	go func() {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := aRegido.Exec(func(c redis.Conn) (res interface{}, err error) {
					return c.Do("PUBLISH", channel, message)
				})
				if err != nil {
					b.Errorf("exec failed, err = %s", err)
				} else {
					atomic.AddInt32(&counter, 1)
				}
			}
		})
		notifyChan <- struct{}{}
	}()
	<-notifyChan
	time.Sleep(time.Millisecond * 100)
	if counter != 0 {
		b.Errorf("unexpected result, expect = 0, but = %d", counter)
	}
}

func BenchmarkSentinelMode_Exec(b *testing.B) {
	echoStr := "hello world"
	sRedigo := New(sentinel.New(
		sentinel.Addrs([]string{"192.168.0.110:26379"}),
	))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res, err := sRedigo.String(func(c redis.Conn) (res interface{}, err error) {
				return c.Do("ECHO", echoStr)
			})
			if err != nil {
				b.Errorf("exec failed, err = %s", err)
			} else if res != echoStr {
				b.Errorf("unexpected result, expect = %s, but = %s", echoStr, res)
			}
		}
	})
}

func BenchmarkSentinelMode_Sub(b *testing.B) {
	message := "hello world"
	channel := "test-channel"
	sRegido := New(sentinel.New(
		sentinel.Addrs([]string{"192.168.0.110:26379"}),
	))
	var counter int32
	notifyChan := make(chan struct{})
	go sRegido.Sub(func(c redis.PubSubConn) (err error) {
		if err = c.Subscribe(channel); err != nil {
			b.Error(err)
		}
		for {
			switch msg := c.ReceiveWithTimeout(0).(type) {
			case redis.Subscription:
				notifyChan <- struct{}{}
			case redis.Message:
				atomic.AddInt32(&counter, -1)
				if string(msg.Data) != message {
					b.Errorf("unexpected result, expect = %s, but = %s", message, msg.Data)
				}
			case error:
				b.Errorf("receive failed, err = %s", msg)
			}
		}
	})
	<-notifyChan
	go func() {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := sRegido.Exec(func(c redis.Conn) (res interface{}, err error) {
					return c.Do("PUBLISH", channel, message)
				})
				if err != nil {
					b.Errorf("exec failed, err = %s", err)
				} else {
					atomic.AddInt32(&counter, 1)
				}
			}
		})
		notifyChan <- struct{}{}
	}()
	<-notifyChan
	time.Sleep(time.Millisecond * 100)
	if counter != 0 {
		b.Errorf("unexpected result, expect = 0, but = %d", counter)
	}
}

func BenchmarkClusterMode_Exec(b *testing.B) {
	echoStr := "hello world"
	cRedigo := New(cluster.New(
		cluster.Nodes([]string{
			"192.168.0.110:30001", "192.168.0.110:30002", "192.168.0.110:30003",
			"192.168.0.110:30004", "192.168.0.110:30005", "192.168.0.110:30006",
		}),
	))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res, err := cRedigo.String(func(c redis.Conn) (res interface{}, err error) {
				return c.Do("ECHO", echoStr)
			})
			if err != nil {
				b.Errorf("exec failed, err = %s", err)
			} else if res != echoStr {
				b.Errorf("unexpected result, expect = %s, but = %s", echoStr, res)
			}
		}
	})
}

func BenchmarkClusterMode_Sub(b *testing.B) {
	message := "hello world"
	channel := "test-channel"
	cRegido := New(cluster.New(
		cluster.Nodes([]string{
			"192.168.0.110:30001", "192.168.0.110:30002", "192.168.0.110:30003",
			"192.168.0.110:30004", "192.168.0.110:30005", "192.168.0.110:30006",
		}),
	))
	var counter int32
	notifyChan := make(chan struct{})
	go cRegido.Sub(func(c redis.PubSubConn) (err error) {
		if err = c.Subscribe(channel); err != nil {
			b.Error(err)
		}
		for {
			switch msg := c.ReceiveWithTimeout(0).(type) {
			case redis.Subscription:
				notifyChan <- struct{}{}
			case redis.Message:
				atomic.AddInt32(&counter, -1)
				if string(msg.Data) != message {
					b.Errorf("unexpected result, expect = %s, but = %s", message, msg.Data)
				}
			case error:
				b.Errorf("receive failed, err = %s", msg)
			}
		}
	})
	<-notifyChan
	go func() {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := cRegido.Exec(func(c redis.Conn) (res interface{}, err error) {
					return c.Do("PUBLISH", channel, message)
				})
				if err != nil {
					b.Errorf("exec failed, err = %s", err)
				} else {
					atomic.AddInt32(&counter, 1)
				}
			}
		})
		notifyChan <- struct{}{}
	}()
	<-notifyChan
	time.Sleep(time.Millisecond * 100)
	if counter != 0 {
		b.Errorf("unexpected result, expect = 0, but = %d", counter)
	}
}
