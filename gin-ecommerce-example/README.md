gin-ecommerce-example/
├── cmd/
│   └── main.go                # 入口文件，启动 Gin 服务
├── config/
│   └── config.yaml            # 配置文件（数据库连接等）
├── internal/
│   ├── handlers/              # Gin 路由处理器
│   │   └── product_handler.go
│   ├── models/                # 数据模型
│   │   └── product.go
│   ├── repositories/          # 仓库层（数据库操作）
│   │   ├── mysql_repo.go      # MySQL 操作
│   │   ├── redis_repo.go      # Redis 操作
│   │   └── es_repo.go         # Elasticsearch 操作
│   └── services/              # 业务逻辑层
│       └── product_service.go
├── pkg/
│   ├── config/
│   │   └── config.go          # Viper 配置加载
│   └── database/              # 数据库连接
│       ├── mysql.go
│       ├── redis.go
│       └── es.go
├── frontend/
│   └── index.html             # 简单前端网页
├── Dockerfile                 # Go 应用的 Docker 文件
├── docker-compose.yml         # Docker Compose 配置
└── go.mod / go.sum