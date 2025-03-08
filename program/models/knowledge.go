package models

import (
	"ai-knowledge/internal/db"
	"time"

	"github.com/google/uuid"
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
	GroupKey  string `gorm:"column:group_key" json:"group_key"`
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (Knowledge) TableName() string {
	return "knowledge"
}

// 生成一个分组标识
func (m *Knowledge) GenGroupKey() string {
	return uuid.NewString()
}

// 批量创建
func (m *Knowledge) BatchCreate(knowledges []*Knowledge) error {
	return db.GormHandler.Table(m.TableName()).Create(knowledges).Error
}

// 根据向量id列表查询
func (m *Knowledge) BatchGetByIds(ids []int64) ([]*Knowledge, error) {
	knowledges := make([]*Knowledge, 0)
	err := db.GormHandler.Table(m.TableName()).Where("vector_id in (?)", ids).Find(&knowledges).Error
	if err == gorm.ErrRecordNotFound {
		return knowledges, nil
	}
	return knowledges, err
}

// 根据id查询
func (m *Knowledge) GetById(id int64) (*Knowledge, error) {
	knowledge := new(Knowledge)
	err := db.GormHandler.Table(m.TableName()).Where("id = ?", id).First(knowledge).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return knowledge, err
}

// 更新数据
func (m *Knowledge) UpdateById(id int64, data map[string]any) error {
	data["updated_at"] = time.Now().Unix()
	return db.GormHandler.Table(m.TableName()).Where("id =?", id).Updates(data).Error
}

// 分组查询分页
func (m *Knowledge) GetList(page, pageSize int, typ int32) (list []*Knowledge, total int64, err error) {
	mydb := db.GormHandler.Table(m.TableName())
	if typ > 0 {
		mydb = mydb.Where("type =?", typ)
	}
	err = mydb.Count(&total).Error
	if err != nil {
		return
	}
	err = mydb.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id desc").
		Scan(&list).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

// 根据分组标识查询
func (m *Knowledge) GetByGroupKey(groupKey string) (list []*Knowledge, err error) {
	err = db.GormHandler.Table(m.TableName()).Where("group_key =?", groupKey).Order("id asc").Scan(&list).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

// 根据id删除
func (m *Knowledge) DelByIds(ids []int64) error {
	return db.GormHandler.Table(m.TableName()).Where("id in (?)", ids).Delete(m).Error
}
