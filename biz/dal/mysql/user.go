package mysql

import "EasyIO/biz/dal/model"

// QueryUserByName 通过用户名查询用户
func QueryUserByName(name string) *model.User {
	var user model.User
	result := DB.Where("username = ?", name).First(&user)
	if result.Error != nil {
		return nil
	}
	return &user
}

// CreateUser 新增用户
func CreateUser(user *model.User) error {
	result := DB.Create(user)
	if result.Error != nil {
		// 处理插入数据错误
		return result.Error
	}
	return nil
}
