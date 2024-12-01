package system

import (
	"project/logic/model"
	"time"
)

type User struct {
	model.Global
	Name          string    `gorm:"size:255;not null;unique" json:"name"` // 用户名
	Password      string    `gorm:"size:255;not null" json:"password"`    // 用户密码
	Nickname      string    `gorm:"size:255;not null" json:"nickname"`    // 昵称
	Picture       string    `gorm:"type:text" json:"avatar"`              // 头像
	Status        bool      `gorm:"default:true" json:"status"`           // 状态
	LastLoginTime time.Time `json:"lastLoginTime"`                        // 最后登录时间
	LastLoginIP   string    `gorm:"type:varchar(45)" json:"lastLoginIp"`  // 最后登录IP (调整为适用于多数据库的类型)
	Role          []Role    `gorm:"many2many:user_bind_role" json:"-"`    // 用户角色关系
}
