基于 [gomodule/redigo](https://github.com/gomodule/redigo) 的二次封装，提供 `stand-alone` `sentinel` `cluster` 3种部署模式下的统一接口，使得更换 `redis` 部署模式对业务透明  


### 开发进度
------------

| Mode部署模式                 | 代码完成度 | 测试完成度 | 依赖包                                                  |
| :--------------------------- | :--------: | :--------: | :------------------------------------------------------ |
| alone 单例，Twemproxy，Codis | 100%       | 100%       |                                                         |
| sentinel 哨兵模式            | 100%       | 100%       | [FZambia/sentinel](https://github.com/FZambia/sentinel) |
| cluster 集群模式             | 100%       | 100%       | [letsfire/redisc](https://github.com/letsfire/redisc)   |


### 方法说明
------------

* 订阅命令方法定义 `type SubFunc func(c redis.PubSubConn) (err error)`
* 普通命令方法定义 `type ExecFunc func(c redis.Conn) (res interface{}, err error)`
* 执行订阅回调方法 Sub(fn SubFunc) (err error)
* 执行命令回调方法 Exec(fn ExecFunc) (interface{}, error)
* 返回当前模式名称 Mode() string
* 下列方法都是封装了Exec方法，格式化返回值的类型
* Int(fn ExecFunc) (int, error)
* Ints(fn ExecFunc) ([]int, error)
* IntMap(fn ExecFunc) (map[string]int, error)
* Int64(fn ExecFunc) (int64, error)
* Int64s(fn ExecFunc) ([]int64, error)
* Int64Map(fn ExecFunc) (map[string]int64, error)
* Uint64(fn ExecFunc) (uint64, error)
* Bool(fn ExecFunc) (bool, error)
* String(fn ExecFunc) (string, error)
* StringMap(fn ExecFunc) (map[string]string, error)
* Strings(fn ExecFunc) ([]string, error)
* Bytes(fn ExecFunc) ([]byte, error)
* ByteSlices(fn ExecFunc) ([][]byte, error)
* Positions(fn ExecFunc) ([]\*[2]float64, error)
* Float64(fn ExecFunc) (float64, error)
* Float64s(fn ExecFunc) ([]float64, error)
* Values(fn ExecFunc) ([]interface{}, error)


### Alone示例
------------

```go
var echoStr = "hello world"
	
var aloneMode = alone.New(
    alone.Addr("192.168.0.110:6379"),
    alone.PoolOpts(
        mode.MaxActive(0),       // 最大连接数，默认0无限制
        mode.MaxIdle(0),         // 最多保持空闲连接数，默认2*runtime.GOMAXPROCS(0)
        mode.Wait(false),        // 连接耗尽时是否等待，默认false
        mode.IdleTimeout(0),     // 空闲连接超时时间，默认0不超时
        mode.MaxConnLifetime(0), // 连接的生命周期，默认0不失效
        mode.TestOnBorrow(nil),  // 空间连接取出后检测是否健康，默认nil
    ),
    alone.DialOpts(
        redis.DialReadTimeout(time.Second),    // 读取超时，默认time.Second
        redis.DialWriteTimeout(time.Second),   // 写入超时，默认time.Second
        redis.DialConnectTimeout(time.Second), // 连接超时，默认500*time.Millisecond
        redis.DialPassword(""),                // 鉴权密码，默认空
        redis.DialDatabase(0),                 // 数据库号，默认0
        redis.DialKeepAlive(time.Minute*5),    // 默认5*time.Minute
        redis.DialNetDial(nil),                // 自定义dial，默认nil
        redis.DialUseTLS(false),               // 是否用TLS，默认false
        redis.DialTLSSkipVerify(false),        // 服务器证书校验，默认false
        redis.DialTLSConfig(nil),              // 默认nil，详见tls.Config
    ),
)

var instance = redigo.New(aloneMode)

res, err := instance.String(func(c redis.Conn) (res interface{}, err error) {
    return c.Do("ECHO", echoStr)
})

if err != nil {
    log.Fatal(err)
} else if res != echoStr {
    log.Fatalf("unexpected result, expect = %s, but = %s", echoStr, res)
}
```


### Sentinel示例
----------------

```go
var echoStr = "hello world"
    
var sentinelMode = sentinel.New(
    sentinel.Addrs([]string{"192.168.0.110:26379"}),
    // 这两项配置和Alone模式完全相同
    // sentinel.PoolOpts(...),
    // sentinel.DialOpts(...),
    // 连接哨兵配置，用法于sentinel.DialOpts()一致
    // 默认未配置的情况则直接使用sentinel.DialOpts()的配置
    // sentinel.SentinelDialOpts()
)

var instance = redigo.New(sentinelMode)

res, err := instance.String(func(c redis.Conn) (res interface{}, err error) {
    return c.Do("ECHO", echoStr)
})

if err != nil {
    log.Fatal(err)
} else if res != echoStr {
    log.Fatalf("unexpected result, expect = %s, but = %s", echoStr, res)
}
```


### Cluster示例
----------------

```go
var echoStr = "hello world"
    
var clusterMode = cluster.New(
    cluster.Nodes([]string{
        "192.168.0.110:30001", "192.168.0.110:30002", "192.168.0.110:30003",
        "192.168.0.110:30004", "192.168.0.110:30005", "192.168.0.110:30006",
    }),
    // 这两项配置和Alone模式完全相同
    // cluster.PoolOpts(...),
    // cluster.DialOpts(...),
)

var instance = redigo.New(clusterMode)

res, err := instance.String(func(c redis.Conn) (res interface{}, err error) {
    return c.Do("ECHO", echoStr)
})

if err != nil {
    log.Fatal(err)
} else if res != echoStr {
    log.Fatalf("unexpected result, expect = %s, but = %s", echoStr, res)
}
```