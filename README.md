# 介绍

scan服务扫链，数据存储到mysql

web服务提供http接口

# 开发 编译

开发

go run cmd/scan/scan.go -c etc/scan/dev.yaml

编译

go build cmd/scan/scan.go

# 目录

```
├── cmd  # 启动目录
│ ├── scan
│ └── web
├── etc  # 配置文件
├── internal  # 核心业务逻辑
│ ├── common
│ │ └── const
│ └── scan
│     ├── config
│     ├── dao
│     └── services
├── pkg  # 基础包
│ ├── log
│ ├── mysql
│ └── token
├── sql  # sql文件
└── vendor  # 包依赖
```