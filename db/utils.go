package db

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"
)

func deriveUserKey(login string, password string) string {
	hash := sha512.New()
	hash.Write([]byte(login + password + time.Now().String()))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	fmt.Println(mdStr)
	return mdStr
}
