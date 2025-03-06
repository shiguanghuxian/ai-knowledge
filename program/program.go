package program

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"ai-knowledge/internal/config"
	"ai-knowledge/internal/db"
	"ai-knowledge/internal/embedding"
	"ai-knowledge/internal/llm"
	"ai-knowledge/internal/logger"
	"ai-knowledge/internal/milvus"
	"ai-knowledge/program/controller"

	"github.com/gin-gonic/gin"
)

// Program 程序实体
type Program struct {
	cfg *config.Config
	srv *http.Server
}

// New 创建程序实例
func New() (*Program, error) {
	// 初始化配置文件
	cfgChan, err := config.NewConfig("")
	if err != nil {
		return nil, err
	}
	cfg := <-cfgChan
	// 初始化日志
	logger.InitLogger(cfg.Debug)

	return &Program{
		cfg: cfg,
	}, nil
}

// Run 启动程序
func (p *Program) Run() {
	if p.cfg.Debug {
		js, _ := json.Marshal(p.cfg)
		log.Println(string(js))
	}
	// 连接数据库
	db.InitDB(p.cfg.Debug, p.cfg.DB)
	// 初始化milvus
	milvus.InitMilvus(p.cfg.Milvus)
	// 初始化向量处理
	embedding.InitTextEmbeddingOperator(p.cfg.Embedding)
	// 初始化llm
	llm.InitLLM(p.cfg.LLM)

	// 启动http监听
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	// 注册路由
	controller.Register(router, p.cfg)
	// 运行服务
	p.srv = &http.Server{
		Addr:    p.cfg.Address,
		Handler: router,
	}
	go func() {
		err := p.srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
		if err != nil {
			log.Println("启动http服务错误", err)
		}
	}()
}

// Stop 程序结束要做的事
func (p *Program) Stop() {
	if p.srv == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// 销毁链接
	milvus.MilvusHandler.Destroy()

	log.Println("Server exiting")
}
