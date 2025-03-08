package milvus

import (
	"ai-knowledge/internal/config"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

var (
	MilvusHandler *MilvusOperator
)

const (
	CollectionName                   = "knowledge"
	dim                              = 1024
	idCol, questionCol, embeddingCol = "ID", "question", "embeddings"
)

// milvus 向量数据库操作
type MilvusOperator struct {
	c client.Client
}

func InitMilvus(cfg *config.MilvusConfig) {
	if cfg == nil {
		log.Panicln("milvus config is nil")
		return
	}
	c, err := client.NewClient(context.Background(), client.Config{
		Address: fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
	})
	if err != nil {
		log.Panicln("milvus connect error", err)
		return
	}

	ctx := context.Background()

	// 检查表是否存在
	has, err := c.HasCollection(ctx, CollectionName)
	if err != nil {
		log.Panicln("milvus has collection error", err)
		return
	}
	// if has {
	// 	c.DropCollection(ctx, CollectionName)
	// 	has = false
	// }
	if !has {
		log.Println("创建集合", CollectionName)
		schema := entity.NewSchema().WithName(CollectionName).WithDescription("存储问题").
			WithField(entity.NewField().WithName(idCol).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithIsAutoID(true)).
			WithField(entity.NewField().WithName(questionCol).WithDataType(entity.FieldTypeVarChar).WithMaxLength(1024)).
			WithField(entity.NewField().WithName(embeddingCol).WithDataType(entity.FieldTypeFloatVector).WithDim(dim))

		// 创建集合
		if err := c.CreateCollection(ctx, schema, entity.DefaultShardNumber); err != nil {
			log.Panicln("创建集合失败，错误:", err)
		}
		// ---
		// build index
		log.Println("开始创建 IVF_FLAT 索引")
		idx, err := entity.NewIndexIvfFlat(entity.L2, 128)
		if err != nil {
			log.Println("创建IVF_FLAT索引失败，错误: ", err)
			return
		}
		if err := c.CreateIndex(ctx, CollectionName, embeddingCol, idx, false); err != nil {
			log.Println("创建索引失败，错误: ", err)
			return
		}
	}

	MilvusHandler = &MilvusOperator{
		c: c,
	}
}

// 获取客户的
func (m *MilvusOperator) GetClient() client.Client {
	return m.c
}

// 批量插入数据
func (m *MilvusOperator) Insert(ctx context.Context, questions []string, embeddings [][]float32) (ids []int64, err error) {
	// 插入数据
	columns := []entity.Column{
		// entity.NewColumnInt64(idCol, nil),
		entity.NewColumnVarChar(questionCol, questions),
		entity.NewColumnFloatVector(embeddingCol, dim, embeddings),
	}
	result, err := m.c.Insert(ctx, CollectionName, "", columns...)
	if err != nil {
		return nil, err
	}
	err = m.c.Flush(ctx, CollectionName, false)
	if err != nil {
		return nil, err
	}
	for i := range result.Len() {
		val, err := result.Get(i)
		if err != nil {
			return nil, err
		}
		fmt.Println("插入结果", val)
		if id, ok := val.(int64); ok {
			ids = append(ids, id)
		}
	}
	return
}

// 更新索引
func (m *MilvusOperator) Update(ctx context.Context, ids []int64, questions []string, embeddings [][]float32) (newIds []int64, err error) {
	// 先删除
	err = m.c.DeleteByPks(ctx, CollectionName, "", entity.NewColumnInt64(idCol, ids))
	if err != nil {
		return nil, err
	}
	// 插入数据
	return m.Insert(ctx, questions, embeddings)
}

// 删除
func (m *MilvusOperator) Delete(ctx context.Context, ids []int64) error {
	// 先删除
	return m.c.DeleteByPks(ctx, CollectionName, "", entity.NewColumnInt64(idCol, ids))
}

// 查询结果
type SearchResult struct {
	Id    int64
	Score float32
	Text  string
}

// 查询
func (m *MilvusOperator) Search(ctx context.Context, vector32 []float32, topK int) (results []*SearchResult, err error) {
	// 加载集合
	err = m.c.LoadCollection(ctx, CollectionName, false)
	if err != nil {
		return nil, err
	}
	// 使用向量查询数据
	sp, _ := entity.NewIndexIvfFlatSearchParam(16)
	vec2search := []entity.Vector{
		entity.FloatVector(vector32),
	}
	sRet, err := m.c.Search(ctx, CollectionName, nil, "", []string{idCol, questionCol}, vec2search,
		embeddingCol, entity.L2, topK, sp)
	if err != nil {
		return nil, err
	}
	// 处理结果
	for _, col := range sRet {
		var idColumn *entity.ColumnInt64
		var questionColumn *entity.ColumnVarChar
		for _, field := range col.Fields {
			if field.Name() == idCol {
				vv, ok := field.(*entity.ColumnInt64)
				if ok {
					idColumn = vv
				}
			}
			if field.Name() == questionCol {
				vv, ok := field.(*entity.ColumnVarChar)
				if ok {
					questionColumn = vv
				}
			}
		}
		if idColumn == nil {
			return nil, errors.New("idColumn is nil")
		}

		for i := range col.ResultCount {
			id, err := idColumn.ValueByIdx(i)
			if err != nil {
				return nil, err
			}
			question, err := questionColumn.ValueByIdx(i)
			if err != nil {
				return nil, err
			}
			results = append(results, &SearchResult{
				Id:    id,
				Score: col.Scores[i],
				Text:  question,
			})
		}
	}
	return results, nil
}

// 销毁
func (m *MilvusOperator) Destroy() {
	m.c.Close()
}
