package embedding

import (
	"ai-knowledge/internal/config"
	"context"
	"log"

	"github.com/tmc/langchaingo/llms/openai"
)

var (
	TextEmbeddingHandler *TextEmbeddingOperator
)

// 处理文本相关向量操作
type TextEmbeddingOperator struct {
	cfg *config.EmbeddingConfig
	llm *openai.LLM
}

func InitTextEmbeddingOperator(cfg *config.EmbeddingConfig) {
	if cfg == nil {
		log.Fatalln("embedding config is nil")
		return
	}
	// openai格式客户的
	llm, err := openai.New(openai.WithBaseURL(cfg.BaseUrl),
		openai.WithToken(cfg.ApiKey), openai.WithEmbeddingModel(cfg.Model))
	if err != nil {
		log.Panicln("init llm failed, err:", err)
		return
	}

	TextEmbeddingHandler = &TextEmbeddingOperator{
		cfg: cfg,
		llm: llm,
	}
}

// 单个计算文本向量
func (t *TextEmbeddingOperator) CalculateEmbedding(ctx context.Context, text string) ([]float32, error) {
	vectors, err := t.llm.CreateEmbedding(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(vectors) == 0 {
		return nil, ErrVectorTransform
	}
	return vectors[0], nil
}

// 批量计算文本向量
func (t *TextEmbeddingOperator) CalculateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	return t.llm.CreateEmbedding(ctx, texts)
}
