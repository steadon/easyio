package mysql

import (
	"EasyIO/biz/dal/model"
	"EasyIO/biz/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// InitDb 初始化函数
func InitDb() error {
	var err error
	dsn := setting.Dns
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

// InitTable 逆向创建表格
func InitTable() error {
	var err error
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Print("逆向创建表失败")
	}
	return nil
}
