package llm

import (
	"ai-knowledge/internal/config"
	"context"
	"log"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

var (
	LLMHandler *LLMOperator
)

// 大模型相关
type LLMOperator struct {
	llm *openai.LLM
}

func InitLLM(cfg *config.LLMConfig) {
	if cfg == nil {
		log.Panicln("llm config is nil")
		return
	}
	llm, err := openai.New(openai.WithBaseURL(cfg.BaseUrl),
		openai.WithToken(cfg.ApiKey), openai.WithModel(cfg.Model))
	if err != nil {
		log.Panicln("init llm failed, err:", err)
		return
	}
	LLMHandler = &LLMOperator{
		llm: llm,
	}
}

func (l *LLMOperator) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return l.llm.Call(ctx, prompt, options...)
}

// LoadStuffQA
func (l *LLMOperator) LoadStuffQA(ctx context.Context, question string, documents []schema.Document) (any, error) {
	stuffQAChain := chains.LoadStuffQA(l.llm)
	// docs := make([]schema.Document, 0)
	// for _, document := range documents {
	// 	docs = append(docs, schema.Document{
	// 		PageContent: document,
	// 		Score:       0,
	// 	})
	// }
	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": documents,
		"question":        question,
	})
	if err != nil {
		return "", err
	}

	return answer["text"], nil
}
