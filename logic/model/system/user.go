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

// type UserUpdate struct {
// 	Global
// 	Nickname *string `json:"nickname"` // 昵称
// 	Picture  *string `json:"avatar"`   // 头像
// 	Status   *bool   `json:"status"`   // 状态
// }

// func SearchUser(db *gorm.DB) ([]User, error) {
// 	var users []User
// 	if err := db.Omit("password").Preload("Role").Find(&users).Error; err != nil {
// 		return nil, fmt.Errorf("查询用户失败: %v", err)
// 	}
// 	return users, nil
// }

// func (user *User) Search(db *gorm.DB) (*User, error) {
// 	var resultUser User
// 	if err := db.Omit("password").Find(&resultUser, user.Id).Error; err != nil {
// 		return nil, fmt.Errorf("查询用户失败: %v", err)
// 	}
// 	return &resultUser, nil
// }

// func (user *User) Register(db *gorm.DB) error {
// 	if err := db.Create(&user).Error; err != nil {
// 		if errors.Is(err, gorm.ErrDuplicatedKey) {
// 			return fmt.Errorf("用户名称已存在: %v", err)
// 		}
// 		return fmt.Errorf("新建用户失败: %v", err)
// 	}
// 	return nil
// }

// func (user *User) PasswordReset(db *gorm.DB, newPassword string) error {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			panic(r)
// 		} else {
// 			tx.Commit()
// 		}
// 	}()
// 	var resultUser *User
// 	if err := tx.First(&resultUser, user.Id).Error; err != nil {
// 		return fmt.Errorf("未找到该用户: %v", err)
// 	}
// 	if resultUser.Password != user.Password {
// 		return fmt.Errorf("原密码不正确")
// 	}
// 	if err := tx.Model(resultUser).Update("password", newPassword).Error; err != nil {
// 		return fmt.Errorf("未找到该用户: %v", err)
// 	}
// 	return nil
// }

// func (user *User) Login(db *gorm.DB) (*User, error) {
// 	var resUser User // 改为非指针类型
// 	if err := db.Where("name = ? AND status = ?", user.Name, true).First(&resUser).Error; err != nil {
// 		return nil, fmt.Errorf("查询该用户失败: %v", err)
// 	}
// 	if isValid := utils.CheckExec(resUser.Password, user.Password); !isValid {
// 		return nil, fmt.Errorf("密码错误")
// 	}
// 	resUser.Password = "" // 清除密码
// 	return &resUser, nil  // 返回指针
// }

// func (user *User) Delete(db *gorm.DB) error {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			panic(r)
// 		} else {
// 			tx.Commit()
// 		}
// 	}()
// 	currentTime := time.Now().Format("2006-01-02 15:04:05")
// 	user.Name = user.Name + "_is_deleted" + currentTime
// 	if err := tx.Updates(&user).Error; err != nil {
// 		return fmt.Errorf("删除用户失败,请检查: %v", err)
// 	}
// 	if err := tx.Delete(&user).Error; err != nil {
// 		return fmt.Errorf("删除用户失败,请检查: %v", err)
// 	}
// 	return nil
// }

// func (user *UserUpdate) Update(db *gorm.DB) error {
// 	var resultUser *User
// 	if err := db.First(&resultUser, user.Id).Error; err != nil {
// 		return fmt.Errorf("未查询到该用户")
// 	}
// 	if err := db.Model(resultUser).Updates(user).Error; err != nil {
// 		return fmt.Errorf("更新用户失败,请检查: %v", err)
// 	}
// 	return nil
// }

// func (user *User) Bind(db *gorm.DB, RoleIds []uint) error {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			panic(r) // re-throw panic after Rollback
// 		} else if tx.Error != nil {
// 			tx.Rollback()
// 		} else {
// 			tx.Commit()
// 		}
// 	}()
// 	// 查询当前用户绑定的角色
// 	if err := tx.First(&user, user.Id).Error; err != nil {
// 		return fmt.Errorf("该用户不存在: %v", err)
// 	}
// 	var roles []Role
// 	for _, id := range RoleIds {
// 		roles = append(roles, Role{Global: Global{Id: id}})
// 	}
// 	// if err := db.Model(&user).Association("UserBindRole").Clear().Error; err != nil {
// 	// 	return fmt.Errorf("绑定角色失败-清空权限步骤: %v", err)
// 	// }
// 	if err := db.Model(user).Association("Role").Replace(roles); err != nil {
// 		return fmt.Errorf("绑定角色失败-增加权限步骤: %v", err)
// 	}
// 	return nil
// }

// func (user *User) GetBind(db *gorm.DB) ([]Role, error) {
// 	var roles []Role
// 	users, err := user.Search(db)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// 执行查询，连接 UserAndRole 和 Role 表
// 	if err := db.Model(users).Association("Role").Find(&roles); err != nil {
// 		return nil, fmt.Errorf("获取用户绑定角色数组失败: %v", err)
// 	}

// 	return roles, nil
// }

// func UpdateLoginMsg(db *gorm.DB, userId uint, ip string) error {
// 	// 更新最后登录时间和IP
// 	if err := db.Model(&User{}).Where("id = ?", userId).Updates(User{
// 		LastLoginTime: time.Now(),
// 		LastLoginIP:   ip,
// 	}).Error; err != nil {
// 		return fmt.Errorf("更新登录信息失败: %v", err)
// 	}
// 	return nil
// }
