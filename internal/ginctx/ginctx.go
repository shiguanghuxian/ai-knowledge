package ginctx

import (
	"ai-knowledge/internal/common"
	"ai-knowledge/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Register(router *gin.RouterGroup) // 注册路由
}

type HandlerFunc func(c *Context)

// Handle 补全必要数据
func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			Context: c,
		}

		// 配置
		cfgV, exists := c.Get(common.CtxSghxCfgKey)
		if exists {
			if cfg, ok := cfgV.(*config.Config); ok && cfg != nil {
				ctx.Cfg = cfg
			}
		}

		h(ctx)
	}
}

type Context struct {
	*gin.Context
	Cfg *config.Config // 配置
}

// 输出json对象
func (c *Context) JSON(code uint32, data interface{}, msg ...string) {
	if len(msg) == 0 {
		msg = []string{""}
	}
	c.Context.JSON(http.StatusOK, c.JSONRoot(code, data, msg[0]))
}

// 输出根结构
func (c *Context) JSONRoot(code uint32, data interface{}, msg string) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"data": data,
		"msg":  msg,
	}
}
