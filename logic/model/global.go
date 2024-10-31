package model

import (
	"time"

	"gorm.io/gorm"
)

type Global struct {
	Id        uint           `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt time.Time      `json:"createTime"`           // 创建时间
	UpdatedAt time.Time      `json:"updateTime"`           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`       // 删除时间
}

func SearchRepart() {

}
