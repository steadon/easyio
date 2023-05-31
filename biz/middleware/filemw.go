package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"path/filepath"
	"strings"
)

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	randomString := base64.URLEncoding.EncodeToString(bytes)[:length]
	randomString = strings.ReplaceAll(randomString, "-", "")
	randomString = strings.ReplaceAll(randomString, "_", "")
	return randomString, nil
}

// GetFileNameWithoutExtension 去掉图片后缀名
func GetFileNameWithoutExtension(filePath string) string {
	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
}
