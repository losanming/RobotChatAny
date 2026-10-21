package pkg

import "math/rand"

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// 获取10位随机大写字符
func GetRandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
