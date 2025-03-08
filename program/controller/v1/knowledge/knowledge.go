package knowledge

import (
	"ai-knowledge/internal/common"
	"ai-knowledge/internal/ginctx"
	"ai-knowledge/internal/logger"
	"ai-knowledge/program/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 知识点管理

type KnowledgeController struct {
}

func (kc *KnowledgeController) Register(router *gin.RouterGroup) {
	router.POST("/knowledge/queryQAndA", ginctx.Handle(kc.QueryQAndA))
	router.POST("/knowledge/saveQAndA", ginctx.Handle(kc.SaveQAndA))
	router.POST("/knowledge/upQAndA", ginctx.Handle(kc.UpQAndA))
	router.POST("/knowledge/saveKnowledge", ginctx.Handle(kc.SaveKnowledge))
	router.POST("/knowledge/upKnowledge", ginctx.Handle(kc.UpKnowledge))
	router.GET("/knowledge/getList", ginctx.Handle(kc.GetList))
	router.GET("/knowledge/getByGroupKey", ginctx.Handle(kc.GetByGroupKey))
	router.POST("/knowledge/delByIds", ginctx.Handle(kc.DelByIds))
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

type UpKnowledgeReq struct {
	Texts []struct {
		Id   int64  `json:"id"`
		Text string `json:"text"`
	} `json:"texts"`
}

// 保存知识点
func (kc *KnowledgeController) UpKnowledge(c *ginctx.Context) {
	req := new(UpKnowledgeReq)
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

	ids := make([]int64, 0)
	texts := make([]string, 0)
	for _, v := range req.Texts {
		if v.Id <= 0 || v.Text == "" {
			logger.Logger.Warnw("参数不合法", "req", req, "v", v)
			c.JSON(1, nil, "参数错误")
			return
		}
		ids = append(ids, v.Id)
		texts = append(texts, v.Text)
	}

	err := service.Knowledge.UpKnowledge(c, ids, texts)
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

type UpQAndAReq struct {
	Questions []struct {
		Id       int64  `json:"id"`
		Question string `json:"question"`
	} `json:"questions"`
	Answer string `json:"answer"`
}

// 修改问答
func (kc *KnowledgeController) UpQAndA(c *ginctx.Context) {
	req := new(UpQAndAReq)
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
	ids := make([]int64, 0)
	questions := make([]string, 0)
	for _, v := range req.Questions {
		if v.Id <= 0 || v.Question == "" {
			logger.Logger.Warnw("参数不合法", "req", req, "v", v)
			c.JSON(1, nil, "参数错误")
			return
		}
		ids = append(ids, v.Id)
		questions = append(questions, v.Question)
	}

	err := service.Knowledge.UpQAndA(c, ids, questions, req.Answer)
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

// 获取列表
func (kc *KnowledgeController) GetList(c *ginctx.Context) {
	typ, _ := strconv.Atoi(c.Query("type"))
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = common.DefaultPageSize
	}
	list, total, err := service.Knowledge.GetList(c, page, pageSize, int32(typ))
	if err != nil {
		logger.Logger.Errorw("处理数据错误", "err", err)
		c.JSON(1, nil, "处理数据错误")
		return
	}

	c.JSON(0, map[string]any{
		"list":  list,
		"total": total,
	}, "成功")
}

// 获取组内全部数据
func (kc *KnowledgeController) GetByGroupKey(c *ginctx.Context) {
	groupKey := c.Query("group_key")
	if groupKey == "" {
		logger.Logger.Warnw("参数不合法", "group_key", groupKey)
		c.JSON(1, nil, "参数错误")
		return
	}

	list, err := service.Knowledge.GetByGroupKey(c, groupKey)
	if err != nil {
		logger.Logger.Errorw("处理数据错误", "err", err)
		c.JSON(1, nil, "处理数据错误")
		return
	}

	c.JSON(0, map[string]any{
		"list": list,
	}, "成功")
}

type DelByIdsReq struct {
	Ids []int64 `json:"ids"`
}

// 删除
func (kc *KnowledgeController) DelByIds(c *ginctx.Context) {
	req := new(DelByIdsReq)
	if err := c.Bind(req); err != nil {
		logger.Logger.Warnw("参数不合法", "err", err)
		c.JSON(1, nil, "参数错误")
		return
	}
	if len(req.Ids) == 0 {
		logger.Logger.Warnw("参数不合法", "req", req)
		c.JSON(1, nil, "参数错误")
		return
	}

	err := service.Knowledge.DelByIds(c, req.Ids)
	if err != nil {
		logger.Logger.Errorw("处理数据错误", "err", err, "req", req)
		c.JSON(1, nil, "处理数据错误")
		return
	}

	c.JSON(0, nil, "成功")
}
