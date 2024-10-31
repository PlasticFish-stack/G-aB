package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const (
	LETTER = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// DIGITS 数字常量
	DIGITS = "0123456789"
	// SPECIALS 特殊字符常量
	SPECIALS = "~=+%^*/()[]{}/!@#$?|"
	// ALL 全部字符常量
	ALL = LETTER + DIGITS + SPECIALS
)

func RandString(length int, source string) string {
	var runes = []rune(source)
	b := make([]rune, length)
	for i := range b {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(runes))))
		if err != nil {
			return ""
		}
		b[i] = runes[index.Int64()]
	}
	return string(b)
}

func HashExec(str string) (string, error) {
	hashe, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码无法被hash加密: %w", err)
	}
	return string(hashe), nil
}

func CheckExec(hashe, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashe), []byte(str))
	return err == nil
}
