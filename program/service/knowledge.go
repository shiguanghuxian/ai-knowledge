package service

import (
	"ai-knowledge/internal/embedding"
	"ai-knowledge/internal/llm"
	"ai-knowledge/internal/logger"
	"ai-knowledge/internal/milvus"
	"ai-knowledge/program/models"
	"context"
	"errors"
	"fmt"

	"github.com/tmc/langchaingo/schema"
)

var (
	Knowledge = new(KnowledgeService)
)

// 知识问答
type KnowledgeService struct {
}

// 保存问答知识
func (s *KnowledgeService) SaveQAndA(ctx context.Context, questions []string, answer string) (err error) {
	if len(questions) == 0 || answer == "" {
		err = ErrInvalidParams
		return
	}
	// 处理问题为向量
	vectors, err := embedding.TextEmbeddingHandler.CalculateEmbeddings(ctx, questions)
	if err != nil {
		logger.Logger.Errorw("处理问题为向量错误", "err", err, "questions", questions, "answer", answer)
		return err
	}

	if len(vectors) != len(questions) {
		return errors.New("向量数与问题数不一致")
	}

	// 写入向量数据库
	ids, err := milvus.MilvusHandler.Insert(ctx, questions, vectors)
	if err != nil {
		logger.Logger.Errorw("写入向量数据库错误", "err", err, "questions", questions, "answer", answer, "vectors", vectors)
		return err
	}
	if len(ids) != len(questions) {
		return errors.New("向量插入数与问题数不一致")
	}

	groupKey := new(models.Knowledge).GenGroupKey()

	// 写入db
	knowledges := make([]*models.Knowledge, 0)
	for i, question := range questions {
		knowledges = append(knowledges, &models.Knowledge{
			Question: question,
			Answer:   answer,
			VectorId: ids[i],
			Type:     models.KnowledgeTypeQAndA,
			GroupKey: groupKey,
		})
	}
	err = new(models.Knowledge).BatchCreate(knowledges)
	if err != nil {
		logger.Logger.Errorw("写入db错误", "err", err, "knowledges", knowledges)
		return err
	}

	return
}

// 更新保存问答知识
func (s *KnowledgeService) UpQAndA(ctx context.Context, ids []int64, questions []string, answer string) (err error) {
	if len(ids) == 0 || len(questions) == 0 || answer == "" {
		err = ErrInvalidParams
		return
	}
	// 处理问题为向量
	vectors, err := embedding.TextEmbeddingHandler.CalculateEmbeddings(ctx, questions)
	if err != nil {
		logger.Logger.Errorw("处理问题为向量错误", "err", err, "questions", questions, "answer", answer)
		return err
	}

	if len(vectors) != len(questions) {
		return errors.New("向量数与问题数不一致")
	}

	// 根据id查询数据
	oldList, err := new(models.Knowledge).BatchGetByIds(ids)
	if err != nil {
		logger.Logger.Errorw("查询db错误", "err", err, "ids", ids)
		return err
	}
	if len(oldList) == 0 {
		err = ErrDataNotFound
		return
	}
	if len(oldList) != len(ids) {
		return errors.New("查询数据与id数不一致")
	}

	// 写入向量数据库
	vectorIds, err := milvus.MilvusHandler.Update(ctx, ids, questions, vectors)
	if err != nil {
		logger.Logger.Errorw("写入向量数据库错误", "err", err, "ids", ids, "questions", questions, "answer", answer, "vectors", vectors)
		return err
	}
	if len(vectorIds) != len(questions) {
		return errors.New("向量插入数与问题数不一致")
	}

	// 写入db
	knowledgeHandler := new(models.Knowledge)
	for i, id := range ids {
		err = knowledgeHandler.UpdateById(id, map[string]any{
			"question":  questions[i],
			"answer":    answer,
			"vector_id": vectorIds[i],
		})
		if err != nil {
			logger.Logger.Errorw("写入db错误", "err", err, "id", id, "question", questions[i], "answer", answer, "vector_id", vectorIds[i])
			return err
		}
	}

	return
}

type SearchKnowledge struct {
	*models.Knowledge
	Score float32
}

// 搜索知识
func (s *KnowledgeService) Search(ctx context.Context, question string, topK int) (answer any, knowledges []*SearchKnowledge, err error) {
	vector, err := embedding.TextEmbeddingHandler.CalculateEmbedding(ctx, question)
	if err != nil {
		logger.Logger.Errorw("处理问题为向量错误", "err", err, "question", question)
		return
	}
	if len(vector) == 0 {
		err = ErrVectorTransform
		return
	}
	results, err := milvus.MilvusHandler.Search(ctx, vector, topK)
	if err != nil {
		logger.Logger.Errorw("搜索向量数据库错误", "err", err, "question", question, "vector", vector, "top_k", topK)
		return
	}
	if len(results) == 0 {
		return
	}
	ids := make([]int64, 0)
	for _, result := range results {
		ids = append(ids, result.Id)
	}
	list, err := new(models.Knowledge).BatchGetByIds(ids)
	if err != nil {
		logger.Logger.Errorw("查询db错误", "err", err, "ids", ids)
		return
	}
	// 整理数据
	for _, v := range list {
		score := float32(0)
		for _, vv := range results {
			if vv.Id == v.VectorId {
				score = vv.Score
				break
			}
		}
		knowledges = append(knowledges, &SearchKnowledge{
			Knowledge: v,
			Score:     score,
		})
	}
	// 调用大模型
	documents := make([]schema.Document, 0)
	for _, v := range knowledges {
		switch v.Type {
		case models.KnowledgeTypePure:
			documents = append(documents, schema.Document{
				PageContent: v.Text,
				Score:       v.Score,
			})
		case models.KnowledgeTypeQAndA:
			documents = append(documents, schema.Document{
				PageContent: fmt.Sprintf("问：%s\n答：%s", v.Question, v.Answer),
				Score:       v.Score,
			})
		}
	}
	answer, err = llm.LLMHandler.LoadStuffQA(ctx, question, documents)
	if err != nil {
		logger.Logger.Errorw("调用大模型错误", "err", err, "question", question, "documents", documents)
		return
	}

	return
}

// 保存知识
func (s *KnowledgeService) SaveKnowledge(ctx context.Context, texts []string) (err error) {
	if len(texts) == 0 {
		err = ErrInvalidParams
		return
	}
	// 处理问题为向量
	vectors, err := embedding.TextEmbeddingHandler.CalculateEmbeddings(ctx, texts)
	if err != nil {
		logger.Logger.Errorw("处理问题为向量错误", "err", err, "texts", texts)
		return err
	}

	if len(vectors) != len(texts) {
		return errors.New("向量数与问题数不一致")
	}

	// 写入向量数据库
	ids, err := milvus.MilvusHandler.Insert(ctx, texts, vectors)
	if err != nil {
		logger.Logger.Errorw("写入向量数据库错误", "err", err, "texts", texts, "vectors", vectors)
		return err
	}
	if len(ids) != len(texts) {
		return errors.New("向量插入数与问题数不一致")
	}

	groupKey := new(models.Knowledge).GenGroupKey()

	// 写入db
	knowledges := make([]*models.Knowledge, 0)
	for i, text := range texts {
		knowledges = append(knowledges, &models.Knowledge{
			Text:     text,
			VectorId: ids[i],
			Type:     models.KnowledgeTypePure,
			GroupKey: groupKey,
		})
	}
	err = new(models.Knowledge).BatchCreate(knowledges)
	if err != nil {
		logger.Logger.Errorw("写入db错误", "err", err, "knowledges", knowledges)
		return err
	}

	return
}

// 保存知识
func (s *KnowledgeService) UpKnowledge(ctx context.Context, ids []int64, texts []string) (err error) {
	if len(ids) == 0 || len(texts) == 0 {
		err = ErrInvalidParams
		return
	}
	// 处理问题为向量
	vectors, err := embedding.TextEmbeddingHandler.CalculateEmbeddings(ctx, texts)
	if err != nil {
		logger.Logger.Errorw("处理问题为向量错误", "err", err, "texts", texts)
		return err
	}

	if len(vectors) != len(texts) {
		return errors.New("向量数与问题数不一致")
	}

	// 根据id查询数据
	oldList, err := new(models.Knowledge).BatchGetByIds(ids)
	if err != nil {
		logger.Logger.Errorw("查询db错误", "err", err, "ids", ids)
		return err
	}
	if len(oldList) == 0 {
		err = ErrDataNotFound
		return
	}
	if len(oldList) != len(ids) {
		return errors.New("查询数据与id数不一致")
	}

	// 写入向量数据库
	vectorIds, err := milvus.MilvusHandler.Update(ctx, ids, texts, vectors)
	if err != nil {
		logger.Logger.Errorw("写入向量数据库错误", "err", err, "texts", texts, "vectors", vectors)
		return err
	}
	if len(ids) != len(texts) {
		return errors.New("向量插入数与问题数不一致")
	}

	// 写入db
	knowledgeHandler := new(models.Knowledge)
	for i, id := range ids {
		err = knowledgeHandler.UpdateById(id, map[string]any{
			"text":      texts[i],
			"vector_id": vectorIds[i],
		})
		if err != nil {
			logger.Logger.Errorw("写入db错误", "err", err, "id", id, "text", texts[i], "vector_id", vectorIds[i])
			return err
		}
	}

	return
}

// 分页查询
func (s *KnowledgeService) GetList(ctx context.Context, page, pageSize int, typ int32) (list []*models.Knowledge, total int64, err error) {
	return new(models.Knowledge).GetList(page, pageSize, typ)
}

// 获取组内全部数据
func (s *KnowledgeService) GetByGroupKey(ctx context.Context, groupKey string) (list []*models.Knowledge, err error) {
	return new(models.Knowledge).GetByGroupKey(groupKey)
}

// 根据id列表删除
func (s *KnowledgeService) DelByIds(ctx context.Context, ids []int64) (err error) {
	if len(ids) == 0 {
		err = ErrInvalidParams
		return
	}
	// 查询出数据
	list, err := new(models.Knowledge).BatchGetByIds(ids)
	if err != nil {
		logger.Logger.Errorw("查询db错误", "err", err, "ids", ids)
		return err
	}
	vectorIds := make([]int64, 0)
	for _, v := range list {
		vectorIds = append(vectorIds, v.VectorId)
	}

	// 删除向量
	err = milvus.MilvusHandler.Delete(ctx, vectorIds)
	if err != nil {
		logger.Logger.Errorw("删除向量数据库错误", "err", err, "vector_ids", vectorIds)
		return err
	}
	// 删除db
	err = new(models.Knowledge).DelByIds(ids)
	if err != nil {
		logger.Logger.Errorw("删除db错误", "err", err, "ids", ids)
		return err
	}

	return
}
