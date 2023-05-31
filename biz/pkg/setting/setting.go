package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg *ini.File

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	Proxy      string
	Prefix     string
	StorageDir string

	Dns string

	SIGN string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}
	LoadServer()
	LoadIO()
	LoadMysql()
	LoadJWT()
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadIO() {
	sec, err := Cfg.GetSection("io")
	if err != nil {
		log.Fatalf("Fail to get section 'request': %v", err)
	}
	Proxy = sec.Key("PROXY").MustString("images")
	Prefix = sec.Key("PREFIX").MustString("http://localhost:8000/")
	StorageDir = sec.Key("STORAGE_DIR").MustString("images")
}

func LoadMysql() {
	sec, err := Cfg.GetSection("mysql")
	if err != nil {
		log.Fatalf("Fail to get section 'request': %v", err)
	}
	Dns = sec.Key("DNS").MustString("root:123456@tcp(localhost:3306)/easyio?charset=utf8mb4&parseTime=True&loc=Local")
}

func LoadJWT() {
	sec, err := Cfg.GetSection("jwt")
	if err != nil {
		log.Fatalf("Fail to get section 'request': %v", err)
	}
	Dns = sec.Key("SIGN").MustString("easyio-key")
}
