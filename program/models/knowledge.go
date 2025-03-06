package models

import (
	"ai-knowledge/internal/db"

	"gorm.io/gorm"
)

var (
	// 类型 0问答 1纯知识
	KnowledgeTypeQAndA = int32(0)
	KnowledgeTypePure  = int32(1)
)

// 知识库
type Knowledge struct {
	Id        int64  `gorm:"column:id;primary_key" json:"id"`
	Question  string `gorm:"column:question" json:"question"`
	Answer    string `gorm:"column:answer" json:"answer"`
	Text      string `gorm:"column:text" json:"text"`
	VectorId  int64  `gorm:"column:vector_id" json:"vector_id"`
	Type      int32  `gorm:"column:type" json:"type"`
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (Knowledge) TableName() string {
	return "knowledge"
}

// 批量创建
func (m *Knowledge) BatchCreate(knowledges []*Knowledge) error {
	return db.GormHandler.Table(m.TableName()).Create(knowledges).Error
}

// 根据向量id列表查询
func (m *Knowledge) BatchGetByIds(ids []int64) ([]*Knowledge, error) {
	knowledges := make([]*Knowledge, 0)
	err := db.GormHandler.Table(m.TableName()).Where("vector_id in ?", ids).Find(&knowledges).Error
	if err == gorm.ErrRecordNotFound {
		return knowledges, nil
	}
	return knowledges, err
}
