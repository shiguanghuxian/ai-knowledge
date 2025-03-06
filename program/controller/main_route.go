package controller

import (
	"ai-knowledge/internal/common"
	"ai-knowledge/internal/config"
	v1 "ai-knowledge/program/controller/v1"

	"github.com/gin-gonic/gin"
)

/* 路由 */

func Register(router *gin.Engine, cfg *config.Config) {
	// 初始化配置
	router.Use(InitCfg(cfg))
	// v1 非管理员接口
	v1Group := router.Group("/v1")
	v1.Register(v1Group)

}

// 给上下文增加配置
func InitCfg(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(common.CtxSghxCfgKey, cfg)
		c.Next()
	}
}
