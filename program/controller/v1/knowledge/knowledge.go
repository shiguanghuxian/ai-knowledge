package knowledge

import (
	"ai-knowledge/internal/common"
	"ai-knowledge/internal/ginctx"
	"ai-knowledge/internal/logger"
	"ai-knowledge/program/service"

	"github.com/gin-gonic/gin"
)

// 知识点管理

type KnowledgeController struct {
}

func (kc *KnowledgeController) Register(router *gin.RouterGroup) {
	router.POST("/knowledge/saveQAndA", ginctx.Handle(kc.SaveQAndA))
	router.POST("/knowledge/queryQAndA", ginctx.Handle(kc.QueryQAndA))
	router.POST("/knowledge/saveKnowledge", ginctx.Handle(kc.SaveKnowledge))
}

type SaveKnowledgeReq struct {
	Texts []string `json:"texts"`
}

// 保存知识点
func (kc *KnowledgeController) SaveKnowledge(c *ginctx.Context) {
	req := new(SaveKnowledgeReq)
	if err := c.Bind(req); err != nil {
		logger.Logger.Warnw("参数不合法", "err", err)
		c.JSON(1, nil, "参数错误")
		return
	}
	if len(req.Texts) == 0 {
		logger.Logger.Warnw("参数不合法", "req", req)
		c.JSON(1, nil, "参数错误")
		return
	}
	for _, v := range req.Texts {
		if v == "" {
			logger.Logger.Warnw("参数不合法", "req", req, "v", v)
			c.JSON(1, nil, "参数错误")
			return
		}
	}

	err := service.Knowledge.SaveKnowledge(c, req.Texts)
	if err != nil {
		logger.Logger.Errorw("处理数据错误", "err", err, "req", req)
		c.JSON(1, nil, "处理数据错误")
		return
	}
	c.JSON(0, nil, "成功")
}

type SaveQAndAReq struct {
	Questions []string `json:"questions"`
	Answer    string   `json:"answer"`
}

// 保存问答
func (kc *KnowledgeController) SaveQAndA(c *ginctx.Context) {
	req := new(SaveQAndAReq)
	if err := c.Bind(req); err != nil {
		logger.Logger.Warnw("参数不合法", "err", err)
		c.JSON(1, nil, "参数错误")
		return
	}
	if len(req.Questions) == 0 || req.Answer == "" {
		logger.Logger.Warnw("参数不合法", "req", req)
		c.JSON(1, nil, "参数错误")
		return
	}
	for _, v := range req.Questions {
		if v == "" {
			logger.Logger.Warnw("参数不合法", "req", req, "v", v)
			c.JSON(1, nil, "参数错误")
			return
		}
	}

	err := service.Knowledge.SaveQAndA(c, req.Questions, req.Answer)
	if err != nil {
		logger.Logger.Errorw("处理数据错误", "err", err, "req", req)
		c.JSON(1, nil, "处理数据错误")
		return
	}

	c.JSON(0, nil, "成功")
}

type QueryQAndAReq struct {
	Question string `json:"question"`
	TopK     int    `json:"top_k"`
}

// 查询一个问题答案
func (kc *KnowledgeController) QueryQAndA(c *ginctx.Context) {
	req := new(QueryQAndAReq)
	if err := c.Bind(req); err != nil {
		logger.Logger.Warnw("参数不合法", "err", err)
		c.JSON(1, nil, "参数错误")
		return
	}
	if req.Question == "" {
		logger.Logger.Warnw("参数不合法", "req", req)
		c.JSON(1, nil, "参数错误")
		return
	}
	if req.TopK <= 0 {
		req.TopK = common.DefaultTopK
	}

	answer, knowledges, err := service.Knowledge.Search(c, req.Question, req.TopK)
	if err != nil {
		logger.Logger.Errorw("处理数据错误", "err", err, "req", req)
		c.JSON(1, nil, "处理数据错误")
		return
	}
	c.JSON(0, map[string]any{
		"answer":     answer,
		"knowledges": knowledges,
	}, "成功")
}
