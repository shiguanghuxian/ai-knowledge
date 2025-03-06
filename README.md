# 一个简单的基于AI的问答系统

## 项目介绍
本项目是一个基于AI的问答系统，使用了OpenAI接口规范模型访问(测试使用ollama)+langchaingo(大模型编排框架)+milvus(向量数据库)+mysql(元数据存储)等技术。

#### 支持两种数据存入方式：
1. 基于问题和答案形式的数据，支持多问题+答案的形式，多问题用于向量索引提高问题命中率。
2. 基于纯知识的形式的数据，支持纯文本形式的数据。

#### 使用方式
1. [x] 使用接口调用
2. [ ] 接入企微自定义应用
...

#### 接口调用
1.基于问题和答案形式的数据

```bash
curl --request POST \
  --url http://127.0.0.1:19090/v1/knowledge/saveQAndA \
  --header 'Content-Type: application/json' \
  --data '{
	"questions": ["猕猴桃是甜的吗", "西瓜是甜的吗"],
	"answer": "是的"
}'
```

2.基于纯知识的形式的数据

```bash
curl --request POST \
  --url http://127.0.0.1:19090/v1/knowledge/saveKnowledge \
  --header 'Content-Type: application/json' \
  --data '{
	"texts": ["猕猴桃是甜的", "西瓜是甜的"]
}'
```

3.基于问题和答案形式的数据查询

```bash
curl --request POST \
  --url http://127.0.0.1:19090/v1/knowledge/queryQAndA \
  --header 'Content-Type: application/json' \
  --data '{
	"question": "西瓜甜吗",
	"top_k": 3
}'
```

## 项目结构
```
.
├── LICENSE
├── Makefile # 编译脚本
├── README.md
├── bin # 可执行文件，发布目录
│   ├── ai-knowledge
│   └── config
│       └── cfg.toml
├── docker # 中间件
│   ├── docker-compose.yml
│   └── volumes
├── go.mod
├── go.sum
├── internal
│   ├── common # 公共模块
│   │   ├── common.go
│   │   └── const.go
│   ├── config # 配置模块
│   │   ├── config.go
│   │   └── watch.go
│   ├── db # mysql数据库模块
│   │   └── db.go
│   ├── embedding # 向量处理模块
│   │   ├── embedding.go
│   │   └── text.go
│   ├── ginctx # 接口上下文模块
│   │   └── ginctx.go
│   ├── llm # llm模块
│   │   └── llm.go
│   ├── logger # 日志模块
│   │   └── logger.go
│   └── milvus # milvus向量数据库模块
│       └── milvus.go
├── main.go
├── program # 业务逻辑
│   ├── controller # 接口模块
│   │   ├── main_route.go
│   │   └── v1
│   │       ├── knowledge
│   │       │   └── knowledge.go
│   │       └── v1.go
│   ├── models # 模型模块
│   │   └── knowledge.go
│   ├── program.go # 主程序
│   └── service # 业务逻辑模块
│       ├── knowledge.go
│       └── service.go
└── sql # 数据库脚本
    └── mysql.sql
```
