package generate

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	StringAlpha        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	StringAlphaNumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	StringNumeric      = "0123456789"
	StringHex          = "abcdef0123456789"
)

// UUID generator
func UUID() string {
	uuidv4 := uuid.NewV4()
	sessionUUID := uuidv4.String()
	sessionTime := time.Now().Unix()

	sessionString := fmt.Sprintf("%d%s", sessionTime, sessionUUID)

	hasher := md5.New()
	hasher.Write([]byte(sessionString))

	return hex.EncodeToString(hasher.Sum(nil))
}

// MD5 generator
func MD5(value string) string {
	hasher := md5.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}

// SHA1 generator
func SHA1(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// RandomString generator
func RandomString(l int, letters string) string {

	var letterRunes = []rune{}

	// default to ALPHA if no runes
	if len(letters) == 0 {
		letterRunes = []rune(StringAlpha)
	} else {
		letterRunes = []rune(letters)
	}

	b := make([]rune, l)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
