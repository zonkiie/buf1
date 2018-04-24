package funclib

import (
	"math/rand"
	"strings"
	"time"

	crypt "github.com/amoghe/go-crypt"
	"github.com/pkg/errors"
)

var src = rand.NewSource(time.Now().UnixNano())

/**
 * @see https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
 */
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// Functio does only check sha256+ passwords, other password hashes let crypt return "invalid argument"
//cryptfunction := func(hash string, text string) error {
func CryptCheck(hash string, text string) error {
	encr, cerr := crypt.Crypt(text, hash)
	if cerr != nil {
		return cerr
	} else {
		if encr == hash {
			return nil
		} else {
			return errors.New("Password does not match!")
		}
	}
}

/**
 * algo can be:
 * - MD5
 * - SHA-256
 * - SHA-512
 */
func CreateSalt(algo string) (string, error) {
	ualgo := strings.ToUpper(algo)
	var salt string
	switch ualgo {
	case "MD5":
		salt = "$1$" + RandStringBytesMaskImprSrc(8)
	case "SHA-256":
		salt = "$5$rounds=97531$" + RandStringBytesMaskImprSrc(16) + "$"
	case "SHA-512":
		salt = "$6$rounds=97531$" + RandStringBytesMaskImprSrc(16) + "$"
	default:
		return "", errors.New("Algo not known!")
	}
	return salt, nil
}

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
