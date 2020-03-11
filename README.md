**提供 `alone` `sentinel` `cluster` 通用API，更换部署模式对业务透明**

| Mode部署模式                 | 代码完成度 | 测试完成度 | 依赖包                                                  |
| :--------------------------- | :--------: | :--------: | :------------------------------------------------------ |
| alone 单机模式               | 100%       | 100%       | [gomodule/redigo](https://github.com/gomodule/redigo)   |
| sentinel 哨兵模式            | 100%       | 100%       | [FZambia/sentinel](https://github.com/FZambia/sentinel) |
| cluster 集群模式             | 100%       | 100%       | [mna/redisc](https://github.com/mna/redisc)             |