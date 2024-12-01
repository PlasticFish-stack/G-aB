package controll

import (
	"fmt"
	"project/logic"
	"project/logic/model"
	"time"
)

type JwtToken struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpTime      time.Time `json:"expires"`
}
type ResultUser struct {
	Avatar   *string  `json:"avatar"`
	Username string   `json:"username"`
	Nickname *string  `json:"nickname"`
	Role     []string `json:"roles"`
	JwtToken
}

func Login(username string, password string, ip string) (*ResultUser, error) {
	var userResult ResultUser
	var result *model.User
	var err error
	var roles []model.Role
	user := &model.User{
		Name:     username,
		Password: password,
	}

	if result, err = user.Login(logic.Gorm); err != nil {
		return nil, err
	}
	if err = model.UpdateLoginMsg(logic.Gorm, result.Id, ip); err != nil {
		return nil, err
	}
	if roles, err = GetBindUserRole(result.Id); err != nil {
		return nil, err
	}
	var rolesId []uint
	var role []string
	for _, v := range roles {
		rolesId = append(rolesId, v.Id)
		role = append(role, v.Identifier)
	}
	userResult.AccessToken,
		userResult.RefreshToken,
		userResult.ExpTime,
		err = GenerateJwt(result.Name, result.Id, rolesId)
	if err != nil {
		return nil, fmt.Errorf("签发token失败: %v", err)
	}
	// roles, err := result.GetBind(logic.Gorm)
	// fmt.Println(roles)
	// if err != nil {
	// 	return nil, fmt.Errorf("获取用户角色失败: %v", err)
	// }
	// var role []string
	// for _, v := range roles {
	// 	role = append(role, v.Identifier)
	// }

	userResult.Avatar = &result.Picture
	userResult.Username = result.Name
	userResult.Nickname = &result.Nickname

	userResult.Role = role

	return &userResult, err
}

func Register(username string, password string, picture string, nickname string) error {
	user := &model.User{
		Name:     username,
		Password: password,
		Picture:  picture,
		Nickname: nickname,
	}
	if err := user.Register(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func UserDelete(id uint, username string) error {
	var User = model.User{
		Global: model.Global{Id: id},
	}
	resultUser, err := User.Search(logic.Gorm)
	if err != nil {
		return err
	}
	if username != resultUser.Name {
		return fmt.Errorf("用户名与id不匹配")
	}
	if err := resultUser.Delete(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func UserUpdate(id uint, nickname string, picture string, status bool) error {
	updateInfo := model.UserUpdate{
		Global:   model.Global{Id: id},
		Nickname: &nickname,
		Picture:  &picture,
		Status:   &status,
	}
	if err := updateInfo.Update(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func UserGetAll() ([]model.User, error) {
	var users []model.User
	var err error
	if users, err = model.SearchUser(logic.Gorm); err != nil {
		return nil, err
	}
	// var ResultUsers []ResultUser
	// for _, v := range users {
	// 	role, err := v.GetBind(logic.Gorm)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	UserResult := GetAllResult{
	// 		User: v,
	// 		Role: role,
	// 	}
	// 	ResultUsers = append(ResultUsers, &UserResult)
	// }
	return users, nil
}

func BindUserRole(userId uint, roleId []uint) error {
	user := &model.User{
		Global: model.Global{Id: userId},
	}
	// roleIds := make([]uint, len(roleId))
	// for i, idStr := range roleId {
	// 	id, err := strconv.Atoi(idStr)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	roleIds[i] = uint(id)
	// }
	err := user.Bind(logic.Gorm, roleId)
	if err != nil {
		return err
	}
	return nil
}

func GetBindUserRole(userId uint) ([]model.Role, error) {
	user := &model.User{
		Global: model.Global{Id: userId},
	}
	roles, err := user.GetBind(logic.Gorm)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
