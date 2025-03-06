package v1

import (
	"ai-knowledge/internal/ginctx"
	"ai-knowledge/program/controller/v1/knowledge"

	"github.com/gin-gonic/gin"
)

// 初始化v1路由

var (
	allController = make([]ginctx.Controller, 0)
)

func init() {
	allController = append(allController, new(knowledge.KnowledgeController))
}

func Register(router *gin.RouterGroup) {
	for _, v := range allController {
		v.Register(router)
	}
}
