package util

import (
	"math/rand"
	"time"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func RandomString(n int) string {
	b := make([]rune, n)
	k := len(letters)

	for i := range b {
		b[i] = letters[rand.Intn(k)]
	}
	return string(b)
}

func RandomTime() time.Time {
	// 生成随机的年份
	year := time.Now().Year() - rand.Intn(10) + 1

	// 生成随机的月份
	month := time.Now().Month() - time.Month(rand.Intn(12)+1)

	// 生成随机的日期
	day := rand.Intn(30) + 1

	// 生成随机的小时
	hour := rand.Intn(24)

	// 生成随机的分钟
	minute := rand.Intn(60)

	return time.Date(year, month, day, hour, minute, 0, 0, time.Now().Location())
}
