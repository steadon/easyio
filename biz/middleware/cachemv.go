package middleware

import (
	"github.com/patrickmn/go-cache"
	"log"
	"time"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

// DeleteCache 删除缓存
func DeleteCache(k string) {
	_, found := c.Get(k)
	if found == true {
		c.Delete(k)
		log.Println("Cache deleted:", k)
	} else {
		log.Println("Cache not found:", k)
	}
}

// AddCache 添加缓存
func AddCache(k string, v interface{}) {
	err := c.Add(k, v, cache.DefaultExpiration)
	if err != nil {
		log.Println("Cache Added Failed:", err)
		return
	}
	log.Println("Cache Added:", k)
}

// GetCache 获取缓存
func GetCache(k string) interface{} {
	v, _ := c.Get(k)
	return v
}
