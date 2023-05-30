package main

import (
	"EasyIO/biz/dal/mysql"
	"EasyIO/biz/pkg/setting"
	"EasyIO/biz/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 创建一个新的 Gin 引擎实例
	r := router.InitRouter()

	// 启动 mysql 驱动
	_ = mysql.InitDb()

	// 逆向创建表格
	_ = mysql.InitTable()

	// 启动 HTTP 服务器
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        r,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务器并处理错误
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
